// Copyright 2016-2026, Pulumi Corporation.
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

package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pulumi/pulumi-kubernetes/tests/v4/crdsdk"
)

func TestCRDExtensionSchemaHasExpectedResources(t *testing.T) {
	for _, testCase := range crdsdk.Cases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			pkg := testCase.Schema(t)
			for _, token := range testCase.ExpectedResources {
				require.Contains(t, pkg.Resources, token)
			}
		})
	}
}
