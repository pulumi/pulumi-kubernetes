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

package test

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

var baseOptions = &integration.ProgramTestOptions{
	Verbose: true,
	Dependencies: []string{
		"github.com/pulumi/pulumi-kubernetes/sdk/v3",
	},
	Env: []string{
		"PULUMI_K8S_CLIENT_BURST=200",
		"PULUMI_K8S_CLIENT_QPS=100",
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
			ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
				// Verify resource creation order using the Event stream. The Chart resources must be created
				// first, followed by the dependent ConfigMap. (The ConfigMap doesn't actually need the Chart, but
				// it creates almost instantly, so it's a good choice to test creation ordering)
				cmRegex := regexp.MustCompile(`ConfigMap::nginx-server-block`)
				svcRegex := regexp.MustCompile(`Service::nginx`)
				deployRegex := regexp.MustCompile(`Deployment::nginx`)
				dependentRegex := regexp.MustCompile(`ConfigMap::foo`)

				var configmapFound, serviceFound, deploymentFound, dependentFound bool
				for _, e := range stackInfo.Events {
					if e.ResOutputsEvent != nil {
						switch {
						case cmRegex.MatchString(e.ResOutputsEvent.Metadata.URN):
							configmapFound = true
						case svcRegex.MatchString(e.ResOutputsEvent.Metadata.URN):
							serviceFound = true
						case deployRegex.MatchString(e.ResOutputsEvent.Metadata.URN):
							deploymentFound = true
						case dependentRegex.MatchString(e.ResOutputsEvent.Metadata.URN):
							dependentFound = true
						}
						assert.Falsef(t, dependentFound && !(configmapFound && serviceFound && deploymentFound),
							"dependent ConfigMap created before all chart resources were ready")
						fmt.Println(e.ResOutputsEvent.Metadata.URN)
					}
				}
			},
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

	t.Run("Helm Release", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:   filepath.Join(cwd, "helm-release", "step1"),
			Quick: true,
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Release Local", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:   filepath.Join(cwd, "helm-release-local", "step1"),
			Quick: true,
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

	t.Run("Helm Skip CRD Rendering", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:         filepath.Join("helm-skip-crd-rendering", "step1"),
			Quick:       true,
			SkipRefresh: true,
			ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
				assert.NotNil(t, stackInfo.Deployment)
				assert.Equal(t, 8, len(stackInfo.Deployment.Resources))

				for _, res := range stackInfo.Deployment.Resources {
					if res.Type == "kubernetes:core/v1:Pod" {
						annotations, ok := openapi.Pluck(res.Inputs, "metadata", "annotations")
						if strings.Contains(res.ID.String(), "skip-crd") {
							assert.False(t, ok)
						} else {
							assert.True(t, ok)
							assert.Contains(t, annotations, "pulumi.com/skipAwait")
						}
					}
				}
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

	t.Run("Secrets", func(t *testing.T) {
		secretMessage := "secret message for testing"

		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:   filepath.Join(cwd, "secrets"),
			Quick: true,
			Config: map[string]string{
				"message": secretMessage,
			},
			ExpectRefreshChanges: true,
			ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
				assert.NotNil(t, stackInfo.Deployment)
				state, err := json.Marshal(stackInfo.Deployment)
				assert.NoError(t, err)

				assert.NotContains(t, string(state), secretMessage)

				// The program converts the secret message to base64, to make a ConfigMap from it, so the state
				// should also not contain the base64 encoding of secret message.
				assert.NotContains(t, string(state), b64.StdEncoding.EncodeToString([]byte(secretMessage)))
			},
		})
		integration.ProgramTest(t, &options)
	})
}
