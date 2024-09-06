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
	"os"
	"strings"
	"sync"

	"github.com/go-openapi/jsonreference"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gen"
	"github.com/pulumi/pulumi/pkg/v3/codegen/cgstrings"
	pulumischema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"github.com/spf13/pflag"
	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/controller/openapi/builder"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

// TODO(rquitales): Remove the hardcoded package name once upstream extension parameterization is implemented and can be passed. For now, we can only "parameterize" with one package, but it is
// intended to not be a upper bound on the number of packages that can be parameterized.
const crdPackageName = "mycrd"

const definitionPrefix = "#/definitions/"

// parameterizedPackageMap is a map of packages to their respective Pulumi PackageSpecs. This is used to store the
// CRD schemas that are generated from the CRD manifests with an internal lock to prevent concurrent access.
// Note: One package can contain multiple CRD schemas within. This enables users to work with multiple CRD versions across different clusters
// in the same Pulumi program.
type parameterizedPackageMap struct {
	sync.Mutex
	crdSchemas map[string]*pulumischema.PackageSpec
}

// get retrieves the CRD schema for the given package name.
//
//nolint:unused // TODO: Will be used once extension parameterization is implemented.
func (c *parameterizedPackageMap) get(name string) *pulumischema.PackageSpec {
	c.Lock()
	defer c.Unlock()
	return c.crdSchemas[name]
}

// add adds the PackageSpec for a given parameterized package.
func (c *parameterizedPackageMap) add(name string, schema *pulumischema.PackageSpec) {
	c.Lock()
	defer c.Unlock()
	if c.crdSchemas == nil {
		c.crdSchemas = make(map[string]*pulumischema.PackageSpec)
	}

	c.crdSchemas[name] = schema
}

// ParameterizedArgs is the struct that holds the arguments for the Kubernetes Provider parameterization.
// Currently, this is only used for generating types from CRD manifests.
type ParameterizedArgs struct {
	PackageVersion   string
	CRDManifestPaths []string
}

// String returns a string representation of the ParameterizedArgs. This is used for logging purposes.
func (p ParameterizedArgs) String() string {
	return fmt.Sprintf("version: %s, crd-manifests: %v", p.PackageVersion, strings.Join(p.CRDManifestPaths, ", "))
}

// parseCrdArgs parses the user provided arguments for provider parameterization.
func parseCrdArgs(args []string) (*ParameterizedArgs, error) {
	var crdPackageVersion string
	var yamlPaths []string

	flags := pflag.NewFlagSet("crdargs", pflag.PanicOnError)
	flags.StringVarP(&crdPackageVersion, "version", "v", "", "The version of the CRD package.")
	flags.StringArrayVarP(&yamlPaths, "crd-manifests", "c", nil, "The paths to the CRD manifests.")
	err := flags.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("error parsing parameterization args: %w", err)
	}

	if crdPackageVersion == "" {
		return nil, errors.New("package version must be provided")
	}

	if len(yamlPaths) == 0 {
		return nil, errors.New("no locations of yaml files given")
	}

	return &ParameterizedArgs{
		PackageVersion:   crdPackageVersion,
		CRDManifestPaths: yamlPaths,
	}, nil
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
		sw, err := builder.BuildOpenAPIV2(crd, v.Name, builder.Options{V2: true, StripValueValidation: false, StripNullable: false, AllowNonStructural: true})
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

// flattenOpenAPI recursively finds all nested objects in the OpenAPI spec and flattens them into a single object as definitions.
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

			// If the property is not an object, we can skip it.
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

			// Create a new definition for the nested object by joining the parent definition name and the property name.
			// This is to ensure that the nested object is unique and does not conflict with other definitions.
			nestedDefinitionName := definitionName + cgstrings.UppercaseFirst(propertyName)
			sw.Definitions[nestedDefinitionName] = propertySchema
			// Add nested object to the stack to be recursively flattened.
			definitionStack = append(definitionStack, nestedDefinitionName)

			// Reset the property to be a reference to the nested object.
			refName := definitionPrefix + nestedDefinitionName
			ref, err := jsonreference.New(refName)
			if err != nil {
				return fmt.Errorf("error creating OpenAPI json reference for nested object: %w", err)
			}

			definition.Properties[propertyName] = spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: spec.Ref{
						Ref: ref,
					},
				},
			}
		}
	}
	return nil
}

// readCRDManifestFile reads the CRD manifest from the given file path.
func readCRDManifestFile(filepath string) ([]*extensionv1.CustomResourceDefinition, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening CRD manifest file: %w", err)
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

			return nil, fmt.Errorf("error decoding CRD manifest within %q: %w", filepath, err)
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
func generateSchema(swagger *spec.Swagger, baseProvName, baseProvVersion string) *pulumischema.PackageSpec {
	// TODO(rquitales): We need to handle field name normalization here so that we can generate typed SDKs that contain valid field names,
	// for example, not allowing hyphens.
	marshaledOpenAPISchema, err := json.Marshal(swagger)
	if err != nil {
		log.Fatalf("error marshalling OpenAPI spec: %v", err)
	}

	unstructuredOpenAPISchema := make(map[string]any)
	err = json.Unmarshal(marshaledOpenAPISchema, &unstructuredOpenAPISchema)
	if err != nil {
		log.Fatalf("error unmarshalling OpenAPI spec: %v", err)
	}

	pSchema := gen.PulumiSchema(unstructuredOpenAPISchema, gen.WithParameterization(&pulumischema.ParameterizationSpec{
		BaseProvider: pulumischema.BaseProviderSpec{
			Name:    baseProvName,
			Version: baseProvVersion,
		},
		Parameter: marshaledOpenAPISchema,
	}))

	return &pSchema
}

// Parameterize is called by the engine when the Kubernetes provider is used for CRDs.
func (k *kubeProvider) Parameterize(ctx context.Context, req *pulumirpc.ParameterizeRequest) (*pulumirpc.ParameterizeResponse, error) {
	log.Println("Parameterizing CRD schemas...")
	logger.V(9).Info("Parameterizing Pulumi Kubernetes provider")

	switch p := req.Parameters.(type) {
	case *pulumirpc.ParameterizeRequest_Args:
		return k.parameterizeRequest_Args(p)

	case *pulumirpc.ParameterizeRequest_Value:
		return k.parameterizeRequest_Value(p)

	default:
		return nil, fmt.Errorf("provider parameter can only be args or value type, received %T", p)

	}
}

// parameterizeRequest_Args is the implementation for the parameterization of the CRD schemas to create typed SDKs from CRD manifests.
func (k *kubeProvider) parameterizeRequest_Args(p *pulumirpc.ParameterizeRequest_Args) (*pulumirpc.ParameterizeResponse, error) {
	args, err := parseCrdArgs(p.Args.GetArgs())
	if err != nil {
		return nil, err
	}

	logger.V(9).Infof("Parameterized Pulumi Kubernetes provider with user specified args: %s", args)

	var allCRDSpecs []*spec.Swagger

	// We need to iterate through all filepaths provided by the user to generate the CRD schemas. Within each file, we can also have multiple CRDs.
	for _, crdPath := range args.CRDManifestPaths {
		crds, err := readCRDManifestFile(crdPath)
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

	crdsPackageSpec := generateSchema(mergedSpecs, k.name, k.version)

	if crdsPackageSpec != nil {
		k.crdSchemas.add(crdPackageName, crdsPackageSpec)
	}

	return &pulumirpc.ParameterizeResponse{Name: crdPackageName, Version: args.PackageVersion}, nil
}

// parameterizeRequest_Value is a placeholder for the extension parameterization implementation. This allows the provider to reconstruct the necessary types for the CRD schemas
// generated from the CRD manifests. This is where we handle field name denormalization and other necessary transformations to be able to translate the typed SDKs back to the original
// CR schema.
func (k *kubeProvider) parameterizeRequest_Value(_ *pulumirpc.ParameterizeRequest_Value) (*pulumirpc.ParameterizeResponse, error) {
	// TODO(rquitales): Implement the logic to generate the CRD schema from the CRD manifests once extension parameterization is implemented.
	// We will need to handle the mapping of normalized field names (to conform to language requirements) to the original k8s field names.
	return nil, nil
}
