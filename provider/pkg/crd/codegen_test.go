// Copyright 2016-2023, Pulumi Corporation.
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

package crd

import (
	_ "embed"
	"testing"

	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenSchemaFromFileWithSingleCRD(t *testing.T) {
	cs := &CodegenSettings{
		PackageName:    "cert-manager",
		PackageVersion: "1.5.3",
	}
	schema, err := GenerateFromFiles(cs, []string{"../../../tests/testdata/crd/single.yaml"})
	require.NoError(t, err)
	assert.Equal(t, 1, len(schema.Resources))
	res := schema.Resources[0]
	assert.Equal(t, "Order is a type to represent an Order with an ACME server", res.Comment)
	orderProps := res.Properties
	assert.Equal(t, 5, len(orderProps))
	orderSpec := orderProps[3] // Expecting this to be `spec`
	assert.Equal(t, "spec", orderSpec.Name)
	tp := orderSpec.Type
	switch tp := tp.(type) {
	case *pschema.ObjectType:
		assert.Equal(t, 6, len(tp.Properties))
	default:
		assert.Fail(t, "type %s not expected", tp)
	}
}

func TestGenSchemaFromFileWithMultipleCRDs(t *testing.T) {
	cs := &CodegenSettings{
		PackageName:    "cert-manager",
		PackageVersion: "1.0",
	}
	schema, err := GenerateFromFiles(cs, []string{"../../../tests/testdata/crd/multiple.yaml"})
	require.NoError(t, err)
	assert.Equal(t, 5, len(schema.Resources))
}

func TestGenSchemaFromMultipleFilesWithCRDs(t *testing.T) {
	cs := &CodegenSettings{
		PackageName:    "cert-manager",
		PackageVersion: "1.0",
	}
	_, err := GenerateFromFiles(cs, []string{
		"../../../tests/testdata/crd/single.yaml",
		"../../../tests/testdata/crd/multiple.yaml",
	})
	require.NoError(t, err)
}

func TestGenSchemaFromFolderWithCRDs(t *testing.T) {
	cs := &CodegenSettings{
		PackageName:    "cert-manager",
		PackageVersion: "1.0",
	}
	schema, err := GenerateFromFiles(cs, []string{"../../../tests/testdata/crd"})
	require.NoError(t, err)
	assert.Equal(t, 6, len(schema.Resources))
}
