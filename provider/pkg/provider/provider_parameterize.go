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
	"errors"
	"fmt"
	"log"
	"sync"

	pulumischema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"github.com/spf13/pflag"
)

// crdSchemaMap is a map of CRD schemas that are generated from the CRD manifests with an
// internal lock to prevent concurrent access.
type crdSchemaMap struct {
	sync.Mutex
	crdSchemas map[string]*pulumischema.Package
}

func (c *crdSchemaMap) getSchema(name string) *pulumischema.Package {
	c.Lock()
	defer c.Unlock()
	return c.crdSchemas[name]
}

func (c *crdSchemaMap) setCRDSchema(name string, schema *pulumischema.Package) {
	c.Lock()
	defer c.Unlock()
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

// Parameterize is called by the engine when the Kubernetes provider is used for CRDs.
func (k *kubeProvider) Parameterize(ctx context.Context, req *pulumirpc.ParameterizeRequest) (*pulumirpc.ParameterizeResponse, error) {
	log.Println("Parameterizing CRD schemas...")
	crdPackageName := "mycrd"
	var crdPackage *pulumischema.Package

	switch p := req.Parameters.(type) {
	case *pulumirpc.ParameterizeRequest_Args:
		_, err := parseCrdArgs(p.Args.GetArgs())
		if err != nil {
			return nil, err
		}

		// TODO: Implement the logic to generate the CRD schema from the CRD manifests.

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
