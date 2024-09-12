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
	"testing"

	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/equality"
)

func TestSetCRDDefaults(t *testing.T) {
	tests := []struct {
		name     string
		crd      extensionv1.CustomResourceDefinition
		expected extensionv1.CustomResourceDefinition
	}{
		{
			"No defaults need to be set",
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foo",
						ListKind: "fooList",
						Kind:     "foo",
					},
				},
			},
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foo",
						ListKind: "fooList",
						Kind:     "foo",
					},
				},
			},
		},
		{
			"Need to set singular name",
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						ListKind: "fooCustomList",
						Kind:     "foo",
					},
				},
			},
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foo",
						ListKind: "fooCustomList",
						Kind:     "foo",
					},
				},
			},
		},
		{
			"Need to set list name",
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foocustomsingular",
						Kind:     "foo",
					},
				},
			},
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foocustomsingular",
						ListKind: "fooList",
						Kind:     "foo",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setCRDDefaults(&tt.crd) //nolint:gosec // This is a false positive on older versions of golangci-lint. We are already using Go v1.22+

			if !equality.Semantic.DeepEqual(tt.crd, tt.expected) {
				t.Errorf("setCRDDefaults() got = %v, want %v", tt.crd, tt.expected)
			}
		})
	}
}
