// Copyright 2016-2019, Pulumi Corporation.
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

	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
)

// TestGo runs Go SDK tests sequentially to avoid OOM errors in CI
func TestGo(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		integration.ProgramTest(t, &integration.ProgramTestOptions{
			Dir: "basic",
			Dependencies: []string{
				"github.com/pulumi/pulumi-kubernetes/sdk/v2",
			},
			Quick: true,
		})
	})

	t.Run("YAML", func(t *testing.T) {
		integration.ProgramTest(t, &integration.ProgramTestOptions{
			Dir: "yaml",
			Dependencies: []string{
				"github.com/pulumi/pulumi-kubernetes/sdk/v2",
			},
			Quick:                true,
			ExpectRefreshChanges: true,
		})
	})

	t.Run("Helm Local", func(t *testing.T) {
		integration.ProgramTest(t, &integration.ProgramTestOptions{
			Dir: "helm-local",
			Dependencies: []string{
				"github.com/pulumi/pulumi-kubernetes/sdk/v2",
			},
			Quick: true,
		})
	})

	t.Run("Helm Remote", func(t *testing.T) {
		integration.ProgramTest(t, &integration.ProgramTestOptions{
			Dir: "helm",
			Dependencies: []string{
				"github.com/pulumi/pulumi-kubernetes/sdk/v2",
			},
			Quick: true,
		})
	})

	t.Run("Kustomize", func(t *testing.T) {
		integration.ProgramTest(t, &integration.ProgramTestOptions{
			Dir: "kustomize",
			Dependencies: []string{
				"github.com/pulumi/pulumi-kubernetes/sdk/v2",
			},
			Quick: true,
		})
	})
}
