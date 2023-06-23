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

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var baseOptions = &integration.ProgramTestOptions{
	Verbose: true,
	Dependencies: []string{
		"github.com/pulumi/pulumi-kubernetes/sdk/v4",
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
			OrderedConfig: []integration.ConfigValue{
				{
					Key:   "pulumi:disable-default-providers[0]",
					Value: "kubernetes",
					Path:  true,
				},
			},
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Local", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "helm-local", "step1"),
			Quick:                true,
			ExpectRefreshChanges: true,
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
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Import", func(t *testing.T) {
		baseDir := filepath.Join(cwd, "helm-release-import", "step1")
		namespace := getRandomNamespace("importtest")
		require.NoError(t, createRelease("mynginx", namespace, baseDir, true))
		defer func() {
			contract.IgnoreError(deleteRelease("mynginx", namespace))
		}()
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir: baseDir,
			Config: map[string]string{
				"namespace": namespace,
			},
			ExpectRefreshChanges: true,
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				assert.NotEmpty(t, stack.Outputs["svc_ip"])
			},
			NoParallel: true,
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Import Deployment Created by Helm", func(t *testing.T) {
		baseDir := filepath.Join(cwd, "helm-import-deployment", "step1")
		namespace := getRandomNamespace("importdepl")
		require.NoError(t, createRelease("mynginx", namespace, baseDir, true))
		defer func() {
			contract.IgnoreError(deleteRelease("mynginx", namespace))
		}()
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir: baseDir,
			Config: map[string]string{
				"namespace": namespace,
			},
			ExpectRefreshChanges: true,
			NoParallel:           true,
			Verbose:              true,
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Remote", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "helm", "step1"),
			Quick:                true,
			ExpectRefreshChanges: true,
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Release", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "helm-release", "step1"),
			Quick:                true,
			ExpectRefreshChanges: true,
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Release With Empty Local Folder", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "helm-release", "step1"),
			Quick:                true,
			ExpectRefreshChanges: true,
			PrePulumiCommand: func(verb string) (func(err error) error, error) {
				// Create an empty folder to test that the Helm provider doesn't fail when the folder is empty, and we should
				// be fetching from remote.
				emptyDir := filepath.Join(cwd, "helm-release", "step1", "nginx")
				if err := os.MkdirAll(emptyDir, 0700); err != nil {
					return nil, err
				}
				return nil, nil
			},
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Release Local", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "helm-release-local", "step1"),
			Quick:                true,
			ExpectRefreshChanges: true,
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Release Local Compressed", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "helm-release-local-tar", "step1"),
			Quick:                true,
			ExpectRefreshChanges: true,
		})
		integration.ProgramTest(t, &options)
	})

	t.Run("Helm Release Partial Error", func(t *testing.T) {
		// Validate that we only see a single release in the namespace - success or failure.
		validation := func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			var namespace string
			for _, res := range stack.Deployment.Resources {
				if res.Type == "kubernetes:helm.sh/v3:Release" {
					ns, found := res.Outputs["namespace"]
					assert.True(t, found)
					namespace = ns.(string)
				}
			}
			assert.NotEmpty(t, namespace)
			releases, err := listReleases(namespace)
			assert.NoError(t, err)
			assert.Len(t, releases, 1)
		}

		test := baseOptions.With(integration.ProgramTestOptions{
			Dir:                    filepath.Join("helm-partial-error", "step1"),
			SkipRefresh:            false,
			SkipEmptyPreviewUpdate: true,
			SkipPreview:            true,
			Verbose:                true,
			ExpectFailure:          true,
			ExtraRuntimeValidation: validation,
			EditDirs: []integration.EditDir{
				{
					Dir:                    filepath.Join("helm-partial-error", "step2"),
					Additive:               true,
					ExtraRuntimeValidation: validation,
					ExpectFailure:          false,
				},
			},
		})
		integration.ProgramTest(t, &test)
	})

	t.Run("Helm API Versions", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "helm-api-versions", "step1"),
			Quick:                true,
			ExpectRefreshChanges: true,
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
			OrderedConfig: []integration.ConfigValue{
				{
					Key:   "pulumi:disable-default-providers[0]",
					Value: "kubernetes",
					Path:  true,
				},
			},
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

	t.Run("ServerSideApply", func(t *testing.T) {
		options := baseOptions.With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(cwd, "server-side-apply"),
			ExpectRefreshChanges: true,
			OrderedConfig: []integration.ConfigValue{
				{
					Key:   "pulumi:disable-default-providers[0]",
					Value: "kubernetes",
					Path:  true,
				},
			},
			EditDirs: []integration.EditDir{
				{
					Dir:      filepath.Join("server-side-apply", "step2"),
					Additive: true,
					ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
						// Validate patched CustomResource
						crPatchedLabels := stackInfo.Outputs["crPatchedLabels"].(map[string]any)
						fooV, ok, err := unstructured.NestedString(crPatchedLabels, "foo")
						assert.True(t, ok)
						assert.NoError(t, err)
						assert.Equal(t, "foo", fooV)
					},
				},
			},
		})
		integration.ProgramTest(t, &options)
	})
}
