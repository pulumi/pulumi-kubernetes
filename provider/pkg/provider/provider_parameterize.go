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

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gen"
	pulumischema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"github.com/spf13/pflag"
	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/controller/openapi/builder"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

// crdSchemaMap is a map of CRD schemas that are generated from the CRD manifests with an
// internal lock to prevent concurrent access.
type crdSchemaMap struct {
	sync.Mutex
	crdSchemas map[string]*pulumischema.PackageSpec
}

func (c *crdSchemaMap) getSchema(name string) *pulumischema.PackageSpec {
	c.Lock()
	defer c.Unlock()
	return c.crdSchemas[name]
}

func (c *crdSchemaMap) setCRDSchema(name string, schema *pulumischema.PackageSpec) {
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
	var crdSpecs []*spec.Swagger

	setCRDDefaults(crd)

	for _, v := range crd.Spec.Versions {
		if !v.Served {
			continue
		}
		// Defaults are not pruned here, but before being served.
		sw, err := builder.BuildOpenAPIV2(crd, v.Name, builder.Options{V2: true, StripValueValidation: false, StripNullable: false, AllowNonStructural: false})
		if err != nil {
			return nil, err
		}
		crdSpecs = append(crdSpecs, sw)
	}

	return crdSpecs, nil
}

// readCRDManifest reads the CRD manifest from the given file path.
func readCRDManifest(filepath string) (*extensionv1.CustomResourceDefinition, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening CRD manifest file: %w", err)
	}
	defer f.Close()

	decoder := yaml.NewYAMLOrJSONDecoder(f, 4096)
	crd := new(extensionv1.CustomResourceDefinition)
	for {
		err = decoder.Decode(crd)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("error decoding CRD manifest: %w", err)
		}
	}

	return crd, nil
}

// mergeSpecs merges the given OpenAPI specs into a single OpenAPI spec.
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

func generateSchema(swagger *spec.Swagger) *pulumischema.PackageSpec {
	b, err := json.Marshal(swagger)
	if err != nil {
		log.Fatalf("error marshalling OpenAPI spec: %v", err)
	}

	var openAPISchema = map[string]any{}
	err = json.Unmarshal(b, &openAPISchema)
	if err != nil {
		log.Fatalf("error unmarshalling OpenAPI spec: %v", err)
	}

	gen.PascalCaseMapping.Add("stable", "Stable")
	pSchema := gen.PulumiSchema(openAPISchema)

	return &pSchema
}

func openAsOpenAPI(filepath string) (*spec.Swagger, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening OpenAPI spec file: %w", err)
	}
	defer f.Close()

	decoder := yaml.NewYAMLOrJSONDecoder(f, 4096)
	sw := new(spec.Swagger)

	for {
		err = decoder.Decode(sw)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("error decoding OpenAPI spec: %w", err)
		}

	}

	return sw, nil
}

// Parameterize is called by the engine when the Kubernetes provider is used for CRDs.
func (k *kubeProvider) Parameterize(ctx context.Context, req *pulumirpc.ParameterizeRequest) (*pulumirpc.ParameterizeResponse, error) {
	log.Println("Parameterizing CRD schemas...")
	crdPackageName := "mycrd"
	var crdPackage *pulumischema.PackageSpec

	switch p := req.Parameters.(type) {
	case *pulumirpc.ParameterizeRequest_Args:
		args, err := parseCrdArgs(p.Args.GetArgs())
		if err != nil {
			return nil, err
		}

		crd, err := readCRDManifest(args.CRDManifestPaths[0])
		if err != nil {
			return nil, fmt.Errorf("error reading CRD manifest: %w", err)
		}

		crdSpecs, err := crdToOpenAPI(crd)
		if err != nil {
			return nil, fmt.Errorf("error generating OpenAPI specs for given CRD %q: %w", crd.Name, err)
		}

		mergedSpecs, err := mergeSpecs(crdSpecs)
		if err != nil {
			// return nil, fmt.Errorf("error merging OpenAPI specs for given CRD %q: %w", crd.Name, err)

			// Try reading file as OpenAPI spec.
			mergedSpecs, err = openAsOpenAPI(args.CRDManifestPaths[0])
			if err != nil {
				return nil, fmt.Errorf("error reading OpenAPI spec: %w", err)
			}
		}

		crdPackage = generateSchema(mergedSpecs)

		k.crdSchemas.setCRDSchema(crdPackageName, crdPackage)

	case *pulumirpc.ParameterizeRequest_Value:
		// TODO: Implement the logic to generate the CRD schema from the CRD manifests.

	default:
		return nil, fmt.Errorf("provider parameter can only be args or value type, received %T", p)

	}

	if crdPackage != nil {
		k.crdSchemas.setCRDSchema(crdPackageName, crdPackage)
	}

	return &pulumirpc.ParameterizeResponse{}, nil
}
