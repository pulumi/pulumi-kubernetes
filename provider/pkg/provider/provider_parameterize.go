// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/jsonreference"
	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/controller/openapi/builder"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kube-openapi/pkg/validation/spec"

	"github.com/pulumi/pulumi/pkg/v3/codegen/cgstrings"
	pulumischema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gen"
)

// defaultExtensionName is the final fallback when the user does not specify -n
// and the manifest path's base name cannot be used.
const defaultExtensionName = "crds"

// defaultExtensionVersion is used when version=<v> is omitted. A CRD bundle has no
// inherent version, so the value only needs to be stable.
const defaultExtensionVersion = "1.0.0"

// definitionPrefix is the JSON Schema reference prefix used in OpenAPI specs to
// point at shared type definitions in the same document (e.g.
// "#/definitions/Foo" resolves to definitions.Foo). Used when flattenOpenAPI
// hoists a nested object to a top-level definition and leaves a $ref behind.
const definitionPrefix = "#/definitions/"

// manifestDownloadTimeout bounds how long fetching a remote CRD manifest may take.
const manifestDownloadTimeout = 60 * time.Second

// maxManifestBytes caps the size of a downloaded CRD manifest.
const maxManifestBytes = 50 << 20 // 50 MiB

// manifestHTTPClient is used to download remote CRD manifests. Tests may replace it.
var manifestHTTPClient = &http.Client{Timeout: manifestDownloadTimeout}

// parameterizedPackageMap stores the CRD schemas generated from CRD manifests.
type parameterizedPackageMap struct {
	sync.Mutex
	crdSchemas map[string]map[string]*pulumischema.PackageSpec // keyed by [name][version]
}

// get retrieves the CRD schema for the given package name.
func (pm *parameterizedPackageMap) get(name, version string) *pulumischema.PackageSpec {
	pm.Lock()
	defer pm.Unlock()
	return pm.crdSchemas[name][version]
}

// has reports whether any version of the named package is present.
func (pm *parameterizedPackageMap) has(name string) bool {
	pm.Lock()
	defer pm.Unlock()
	return len(pm.crdSchemas[name]) > 0
}

// add adds the PackageSpec for a given parameterized package.
func (pm *parameterizedPackageMap) add(name, version string, schema *pulumischema.PackageSpec) {
	pm.Lock()
	defer pm.Unlock()
	if pm.crdSchemas == nil {
		pm.crdSchemas = make(map[string]map[string]*pulumischema.PackageSpec)
	}
	if pm.crdSchemas[name] == nil {
		pm.crdSchemas[name] = make(map[string]*pulumischema.PackageSpec)
	}

	pm.crdSchemas[name][version] = schema
}

// ParameterizedArgs holds the arguments for the Kubernetes Provider
// parameterization. Currently this is only used for generating types from CRD
// manifests.
type ParameterizedArgs struct {
	ExtensionName    string
	ExtensionVersion string
	CRDManifestPaths []string
}

// String returns a string representation of the ParameterizedArgs. This is used for logging purposes.
func (p ParameterizedArgs) String() string {
	return fmt.Sprintf("name: %s, version: %s, crd-manifest: %v",
		p.ExtensionName, p.ExtensionVersion, strings.Join(p.CRDManifestPaths, ", "))
}

// parseCrdArgs parses the user-provided parameterization arguments. Each argument
// is a key=value pair, e.g. "name=foo version=1.0 crd-manifest=crds.yaml". The
// crd-manifest key may be repeated to supply multiple manifests, and each value
// may be a local path or an http(s) URL (GitHub file links are supported).
// version is optional and defaults to defaultExtensionVersion.
func parseCrdArgs(args []string) (*ParameterizedArgs, error) {
	var extensionName string
	var extensionVersion string
	var yamlPaths []string

	for _, arg := range args {
		key, value, ok := strings.Cut(arg, "=")
		if !ok {
			return nil, fmt.Errorf("invalid parameterization argument %q: expected key=value", arg)
		}
		switch key {
		case "name":
			extensionName = value
		case "version":
			extensionVersion = value
		case "crd-manifest":
			yamlPaths = append(yamlPaths, value)
		default:
			return nil, fmt.Errorf("unknown parameterization argument %q: expected name, version, or crd-manifest", key)
		}
	}

	if extensionVersion == "" {
		extensionVersion = defaultExtensionVersion
	}

	if len(yamlPaths) == 0 {
		return nil, errors.New("no CRD manifests given (crd-manifest=<path-or-url>)")
	}

	return &ParameterizedArgs{
		ExtensionName:    extensionName,
		ExtensionVersion: extensionVersion,
		CRDManifestPaths: yamlPaths,
	}, nil
}

// deriveExtensionName produces an extension name from the first CRD manifest
// path or URL. The file's base name (without extension) is used directly, e.g.
// "gateway-api-crds.yaml" → "gateway-api-crds". For URLs, only the URL path is
// considered, so query strings and fragments are ignored. Falls back to
// defaultExtensionName if no usable path is provided.
func deriveExtensionName(paths []string) string {
	for _, p := range paths {
		var base string
		if u, ok := parseRemoteManifest(p); ok {
			base = path.Base(u.Path)
		} else {
			base = filepath.Base(p)
		}
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		if name != "" && name != "." && name != "/" {
			return name
		}
	}
	return defaultExtensionName
}

// parseRemoteManifest reports whether the given crd-manifest value is an http(s)
// URL, returning the parsed URL if so.
func parseRemoteManifest(s string) (*url.URL, bool) {
	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return nil, false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, false
	}
	return u, true
}

// resolveManifestURL rewrites GitHub file page links into raw content links, so
// that "https://github.com/<owner>/<repo>/blob/<ref>/<path>" (or ".../raw/...")
// becomes "https://raw.githubusercontent.com/<owner>/<repo>/<ref>/<path>". Any
// other URL is returned unchanged.
func resolveManifestURL(u *url.URL) *url.URL {
	if !strings.EqualFold(u.Host, "github.com") && !strings.EqualFold(u.Host, "www.github.com") {
		return u
	}
	// Expect: /<owner>/<repo>/<blob|raw>/<ref>/<path...>
	parts := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 4)
	if len(parts) < 4 || (parts[2] != "blob" && parts[2] != "raw") {
		return u
	}
	return &url.URL{
		Scheme: "https",
		Host:   "raw.githubusercontent.com",
		Path:   "/" + parts[0] + "/" + parts[1] + "/" + parts[3],
	}
}

// isGitHubHost reports whether a GitHub token may be sent to the given host.
func isGitHubHost(host string) bool {
	switch strings.ToLower(host) {
	case "github.com", "www.github.com", "raw.githubusercontent.com":
		return true
	}
	return false
}

// openCRDManifest opens a CRD manifest from a local path or downloads it from an
// http(s) URL. For GitHub hosts, the GITHUB_TOKEN environment variable (if set) is
// sent so that manifests in private repositories can be fetched.
func openCRDManifest(ctx context.Context, location string) (io.ReadCloser, error) {
	u, ok := parseRemoteManifest(location)
	if !ok {
		f, err := os.Open(location)
		if err != nil {
			return nil, fmt.Errorf("error opening CRD manifest file: %w", err)
		}
		return f, nil
	}

	u = resolveManifestURL(u)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for CRD manifest %q: %w", location, err)
	}
	gitHub := isGitHubHost(u.Host)
	if token := os.Getenv("GITHUB_TOKEN"); token != "" && gitHub {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	logger.V(9).Infof("Downloading CRD manifest from %s", u)
	resp, err := manifestHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error downloading CRD manifest %q: %w", location, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		resp.Body.Close()
		msg := fmt.Sprintf("error downloading CRD manifest %q: %s", location, resp.Status)
		if gitHub && resp.StatusCode == http.StatusNotFound {
			msg += " (if the repository is private, set GITHUB_TOKEN)"
		}
		return nil, errors.New(msg)
	}

	return struct {
		io.Reader
		io.Closer
	}{io.LimitReader(resp.Body, maxManifestBytes), resp.Body}, nil
}

// fillDefaultNames sets the default names for the CRD if they are not specified.
// This allows the OpenAPI builder to generate the swagger specs correctly with
// the correct defaults.
func setCRDDefaults(crd *extensionv1.CustomResourceDefinition) {
	if crd.Spec.Names.Singular == "" {
		crd.Spec.Names.Singular = strings.ToLower(crd.Spec.Names.Kind)
	}
	if crd.Spec.Names.ListKind == "" {
		crd.Spec.Names.ListKind = crd.Spec.Names.Kind + "List"
	}
}

// crdToOpenAPI generates the OpenAPI specs for a given CRD manifest.
func crdToOpenAPI(crd *extensionv1.CustomResourceDefinition) ([]*spec.Swagger, error) {
	var openAPIManifests []*spec.Swagger

	setCRDDefaults(crd)

	for _, v := range crd.Spec.Versions {
		if !v.Served {
			continue
		}
		// Defaults are not pruned here, but before being served.
		sw, err := builder.BuildOpenAPIV2(
			crd,
			v.Name,
			builder.Options{V2: true, StripValueValidation: false, StripNullable: true, AllowNonStructural: true},
		)
		if err != nil {
			return nil, err
		}

		err = flattenOpenAPI(sw)
		if err != nil {
			return nil, fmt.Errorf("error flattening OpenAPI spec: %w", err)
		}

		openAPIManifests = append(openAPIManifests, sw)
	}

	return openAPIManifests, nil
}

// flattenOpenAPI recursively finds all nested objects in the OpenAPI spec and flattens them into
// top-level definitions. This includes both direct object properties and objects nested inside
// arrays (e.g., spec.listeners, spec.rules[].backendRefs).
func flattenOpenAPI(sw *spec.Swagger) error {
	// Create a stack of definition names to be processed.
	definitionStack := make([]string, 0, len(sw.Definitions))

	// Populate existing definitions into the stack.
	for defName := range sw.Definitions {
		definitionStack = append(definitionStack, defName)
	}

	for len(definitionStack) != 0 {
		// Pop the last definition from the stack.
		definitionName := definitionStack[len(definitionStack)-1]
		definitionStack = definitionStack[:len(definitionStack)-1]
		// Get the definition from the OpenAPI spec.
		definition := sw.Definitions[definitionName]

		for propertyName, propertySchema := range definition.Properties {
			// If the property is already a reference to a URL, we can skip it.
			if propertySchema.Ref.GetURL() != nil {
				continue
			}

			// Handle arrays whose items are objects — e.g., spec.listeners, spec.rules.
			if propertySchema.Type.Contains("array") {
				itemSchema := propertySchema.Items
				if itemSchema == nil || itemSchema.Schema == nil {
					continue
				}
				inner := itemSchema.Schema
				if inner.Ref.GetURL() != nil {
					continue
				}
				if !inner.Type.Contains("object") || inner.Properties == nil {
					continue
				}
				if inner.AdditionalProperties != nil {
					continue
				}

				nestedDefinitionName := definitionName + cgstrings.UppercaseFirst(propertyName)
				sw.Definitions[nestedDefinitionName] = *inner
				definitionStack = append(definitionStack, nestedDefinitionName)

				ref, err := makeDefinitionRef(nestedDefinitionName)
				if err != nil {
					return err
				}
				propertySchema.Items = &spec.SchemaOrArray{
					Schema: &spec.Schema{
						SchemaProps: spec.SchemaProps{Ref: ref},
					},
				}
				definition.Properties[propertyName] = propertySchema
				continue
			}

			// Handle direct object properties.
			if !propertySchema.Type.Contains("object") {
				continue
			}

			if propertySchema.Properties == nil {
				continue
			}

			// If the property is an object with additional properties, we can skip it. We only care about
			// nested objects that are explicitly defined.
			if propertySchema.AdditionalProperties != nil {
				continue
			}

			// Create a new definition for the nested object by joining the parent definition name and the property
			// name.
			// This is to ensure that the nested object is unique and does not conflict with other definitions.
			nestedDefinitionName := definitionName + cgstrings.UppercaseFirst(propertyName)
			sw.Definitions[nestedDefinitionName] = propertySchema
			// Add nested object to the stack to be recursively flattened.
			definitionStack = append(definitionStack, nestedDefinitionName)

			ref, err := makeDefinitionRef(nestedDefinitionName)
			if err != nil {
				return err
			}
			definition.Properties[propertyName] = spec.Schema{
				SchemaProps: spec.SchemaProps{Ref: ref},
			}
		}
	}
	return nil
}

// makeDefinitionRef builds the {$ref: "#/definitions/<name>"} pointer that
// flattenOpenAPI leaves behind when it hoists a nested object to a top-level
// definition. Helper because flattenOpenAPI does this in two places (direct
// object properties and array-of-object properties).
func makeDefinitionRef(name string) (spec.Ref, error) {
	ref, err := jsonreference.New(definitionPrefix + name)
	if err != nil {
		return spec.Ref{}, fmt.Errorf("error creating OpenAPI json reference for nested object: %w", err)
	}
	return spec.Ref{Ref: ref}, nil
}

// readCRDManifest reads the CRDs from the given local path or http(s) URL.
func readCRDManifest(ctx context.Context, location string) ([]*extensionv1.CustomResourceDefinition, error) {
	f, err := openCRDManifest(ctx, location)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewYAMLOrJSONDecoder(f, 4096)
	var crds []*extensionv1.CustomResourceDefinition

	for {
		crd := new(extensionv1.CustomResourceDefinition)
		err = decoder.Decode(crd)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("error decoding CRD manifest within %q: %w", location, err)
		}

		if crd.Kind != "CustomResourceDefinition" {
			logger.V(9).Info("Skipping document within manifest that could not be decoded as a CRD")
			continue
		}

		crds = append(crds, crd)
	}

	return crds, nil
}

// mergeSpecs merges a slice of OpenAPI specs into a single OpenAPI spec.
func mergeSpecs(specs []*spec.Swagger) (*spec.Swagger, error) {
	if len(specs) == 0 {
		return nil, errors.New("no OpenAPI specs to merge")
	}

	mergedSpecs, err := builder.MergeSpecs(specs[0], specs[1:]...)
	if err != nil {
		return nil, fmt.Errorf("error merging OpenAPI specs: %w", err)
	}

	return mergedSpecs, nil
}

// generateSchema generates the Pulumi schema with parameterization for the given OpenAPI spec.
func generateSchema(
	swagger *spec.Swagger,
	extensionName, extensionVersion, baseProvName, baseProvVersion string,
) *pulumischema.PackageSpec {
	// TODO: property names containing a dot still produce an invalid identifier.
	// Language codegen quotes hyphenated names, but has no equivalent for dots.
	marshaledOpenAPISchema, err := json.Marshal(swagger)
	if err != nil {
		log.Fatalf("error marshalling OpenAPI spec: %v", err)
	}

	unstructuredOpenAPISchema := make(map[string]any)
	err = json.Unmarshal(marshaledOpenAPISchema, &unstructuredOpenAPISchema)
	if err != nil {
		log.Fatalf("error unmarshalling OpenAPI spec: %v", err)
	}

	pSchema := gen.PulumiSchema(unstructuredOpenAPISchema,
		gen.WithParameterization(&pulumischema.ExtensionParameterizationSpec{
			BaseProvider: pulumischema.BaseProviderRefSpec{
				Name:    baseProvName,
				Version: strings.TrimPrefix(baseProvVersion, "v"),
			},
			Parameter: marshaledOpenAPISchema,
		}),
		gen.WithExtensionName(extensionName),
	)

	pSchema.Version = extensionVersion

	return &pSchema
}

// Parameterize is called by the engine when the Kubernetes provider is used for CRDs.
func (k *kubeProvider) Parameterize(
	ctx context.Context,
	req *pulumirpc.ParameterizeRequest,
) (*pulumirpc.ParameterizeResponse, error) {
	log.Println("Parameterizing CRD schemas...")
	logger.V(9).Info("Parameterizing Pulumi Kubernetes provider")

	switch p := req.Parameters.(type) {
	case *pulumirpc.ParameterizeRequest_Args:
		return k.parameterizeRequestArgs(ctx, p)

	case *pulumirpc.ParameterizeRequest_Value:
		return k.parameterizeRequestValue(p)

	default:
		return nil, fmt.Errorf("provider parameter can only be args or value type, received %T", p)

	}
}

// parameterizeRequestArgs is the implementation for the parameterization of the CRD schemas to create typed SDKs from
// CRD manifests.
func (k *kubeProvider) parameterizeRequestArgs(
	ctx context.Context,
	p *pulumirpc.ParameterizeRequest_Args,
) (*pulumirpc.ParameterizeResponse, error) {
	args, err := parseCrdArgs(p.Args.GetArgs())
	if err != nil {
		return nil, err
	}

	logger.V(9).Infof("Parameterized Pulumi Kubernetes provider with user specified args: %s", args)

	var allCRDSpecs []*spec.Swagger

	// We need to iterate through all paths/URLs provided by the user to generate the CRD schemas. Within each
	// manifest, we can also have multiple CRDs.
	for _, crdPath := range args.CRDManifestPaths {
		crds, err := readCRDManifest(ctx, crdPath)
		if err != nil {
			return nil, fmt.Errorf("error reading CRD manifest: %w", err)
		}

		for _, crd := range crds {
			crdVersionSpecs, err := crdToOpenAPI(crd)
			if err != nil {
				return nil, fmt.Errorf("error generating OpenAPI specs for given CRD %q: %w", crd.Name, err)
			}

			mergedCRDVersionSpecs, err := mergeSpecs(crdVersionSpecs)
			if err != nil {
				return nil, fmt.Errorf("error reading OpenAPI spec: %w", err)
			}

			allCRDSpecs = append(allCRDSpecs, mergedCRDVersionSpecs)
		}
	}

	mergedSpecs, err := mergeSpecs(allCRDSpecs)
	if err != nil {
		return nil, fmt.Errorf("error merging OpenAPI specs for all provided CRDs: %w", err)
	}

	// Resolve the extension name: explicit -n > derived from manifest filename > default.
	extensionName := args.ExtensionName
	if extensionName == "" {
		extensionName = deriveExtensionName(args.CRDManifestPaths)
	}

	crdsPackageSpec := generateSchema(mergedSpecs, extensionName, args.ExtensionVersion, k.name, k.version)

	if crdsPackageSpec != nil {
		k.crdSchemas.add(extensionName, args.ExtensionVersion, crdsPackageSpec)
	}

	return &pulumirpc.ParameterizeResponse{Name: extensionName, Version: args.ExtensionVersion}, nil
}

// parameterizeRequestValue reconstructs the CRD schema from the saved
// parameterization state. The engine calls this path on subsequent runs (e.g.,
// `pulumi up` after the initial `pulumi package add`) by replaying the
// parameters that were persisted in Pulumi.yaml.
//
// The Value.Value bytes contain the marshaled OpenAPI spec that was originally
// stored as the Parameter field in the ExtensionParameterizationSpec during the Args
// path. We unmarshal it, regenerate the Pulumi schema, and cache it.
func (k *kubeProvider) parameterizeRequestValue(
	p *pulumirpc.ParameterizeRequest_Value,
) (*pulumirpc.ParameterizeResponse, error) {
	v := p.Value
	if v == nil {
		return nil, errors.New("parameterize value is nil")
	}

	extensionName := v.GetName()
	extensionVersion := v.GetVersion()
	paramBytes := v.GetValue()

	if extensionName == "" {
		return nil, errors.New("parameterize value: extension name must be provided")
	}
	if extensionVersion == "" {
		return nil, errors.New("parameterize value: extension version must be provided")
	}
	if len(paramBytes) == 0 {
		return nil, errors.New("parameterize value: parameter bytes must be provided")
	}

	logger.V(9).Infof("Reconstructing CRD schema for %s@%s from saved parameters", extensionName, extensionVersion)

	// Deserialize the OpenAPI spec from the saved parameter bytes.
	var swagger spec.Swagger
	if err := json.Unmarshal(paramBytes, &swagger); err != nil {
		return nil, fmt.Errorf("parameterize value: error unmarshalling saved OpenAPI spec: %w", err)
	}

	crdsPackageSpec := generateSchema(&swagger, extensionName, extensionVersion, k.name, k.version)
	if crdsPackageSpec != nil {
		k.crdSchemas.add(extensionName, extensionVersion, crdsPackageSpec)
	}

	return &pulumirpc.ParameterizeResponse{Name: extensionName, Version: extensionVersion}, nil
}
