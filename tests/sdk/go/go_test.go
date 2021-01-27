// Copyright 2016-2021, Pulumi Corporation.
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
// +build go all

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

var baseOptions = &integration.ProgramTestOptions{
	Verbose: true,
	Dependencies: []string{
		"github.com/pulumi/pulumi-kubernetes/sdk/v3/go",
	},
}

// TestGo runs Go SDK tests sequentially to avoid OOM errors in CI
func TestGo(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	t.Run("Basic", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "basic"),
			ExpectRefreshChanges: true,
			Quick:                true,
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("YAML", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "yaml"),
			Quick:                true,
			ExpectRefreshChanges: true,
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Local", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:   filepath.Join(cwd, "helm-local", "step1"),
			Quick: true,
			EditDirs: []integration.EditDir{
				{
					Dir:             filepath.Join(cwd, "helm-local", "step2"),
					Additive:        true,
					ExpectNoChanges: true,
				},
			},
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Remote", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:   filepath.Join(cwd, "helm", "step1"),
			Quick: true,
			EditDirs: []integration.EditDir{
				{
					Dir:             filepath.Join("helm", "step2"),
					Additive:        true,
					ExpectNoChanges: true,
				},
			},
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm API Versions", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:   filepath.Join(cwd, "helm-api-versions", "step1"),
			Quick: true,
			EditDirs: []integration.EditDir{
				{
					Dir:             filepath.Join("helm-api-versions", "step2"),
					Additive:        true,
					ExpectNoChanges: true,
				},
			},
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Kustomize", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:   filepath.Join(cwd, "kustomize"),
			Quick: true,
		})
		integration.ProgramTest(t, &options)
	})
}
