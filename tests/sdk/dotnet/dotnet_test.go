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
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
)

var baseOptions = &integration.ProgramTestOptions{
	Verbose: true,
	Dependencies: []string{
		"Pulumi.Kubernetes",
	},
	Env: []string{
		"PULUMI_K8S_CLIENT_BURST=200",
		"PULUMI_K8S_CLIENT_QPS=100",
	},
}

func TestDotnet_Basic(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  "basic",
		Quick:                true,
		ExpectRefreshChanges: true, // The CRD sometimes, but not always, has changes during refresh.
	})
	integration.ProgramTest(t, &test)
}

func TestDotnet_Guestbook(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "guestbook",
		Quick: true,
	})
	integration.ProgramTest(t, &test)
}

func TestDotnet_YamlUrl(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "yaml-url",
		Quick: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(
			t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
		) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 18, len(stackInfo.Deployment.Resources))
		},
	})
	integration.ProgramTest(t, &test)
}

func TestDotnet_YamlLocal(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "yaml-local",
		Quick: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestDotnet_Helm(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm", "step1"),
		Quick:                true,
		ExpectRefreshChanges: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Ensure that all `Services` have `status` marked as a `Secret`
			for _, res := range stackInfo.Deployment.Resources {
				if res.Type == tokens.Type("kubernetes:core/v1:Service") {
					spec, has := res.Outputs["status"]
					assert.True(t, has)
					specMap, is := spec.(map[string]any)
					assert.True(t, is)
					sigKey, has := specMap[resource.SigKey]
					assert.True(t, has)
					assert.Equal(t, resource.SecretSig, sigKey)
				}
			}
		},
	})
	integration.ProgramTest(t, &test)
}

func TestDotnet_HelmLocal(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm-local", "step1"),
		Quick:                true,
		ExpectRefreshChanges: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 12, len(stackInfo.Deployment.Resources))

			// Verify resource creation order using the Event stream. The Chart resources must be created
			// first, followed by the dependent ConfigMap. (The ConfigMap doesn't actually need the Chart, but
			// it creates almost instantly, so it's a good choice to test creation ordering)
			cmRegex := regexp.MustCompile(`ConfigMap::test-.*/nginx-server-block`)
			svcRegex := regexp.MustCompile(`Service::test-.*/nginx`)
			deployRegex := regexp.MustCompile(`Deployment::test-.*/nginx`)
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
	integration.ProgramTest(t, &test)
}

func TestDotnet_HelmApiVersions(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm-api-versions", "step1"),
		Quick:                true,
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 7, len(stackInfo.Deployment.Resources))
		},
	})
	integration.ProgramTest(t, &test)
}

func TestDotnet_HelmKubeVersion(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm-kube-version", "step1"),
		Quick:                true,
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
		},
	})
	integration.ProgramTest(t, &test)
}

func TestDotnet_HelmAllowCRDRendering(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm-skip-crd-rendering", "step1"),
		Quick:                true,
		SkipRefresh:          true,
		ExpectRefreshChanges: true,
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
	integration.ProgramTest(t, &test)
}

func TestDotnet_CustomResource(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  "custom-resource",
		Quick:                true,
		ExpectRefreshChanges: true, // The CRD sometimes, but not always, has changes during refresh.
	})
	integration.ProgramTest(t, &test)
}

func TestDotnet_Kustomize(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "kustomize",
		Quick: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestDotnet_Secrets(t *testing.T) {
	secretMessage := "secret message for testing"

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "secrets",
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
	integration.ProgramTest(t, &test)
}

func TestDotnet_ServerSideApply(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  "server-side-apply",
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		// TODO: Need to support CustomResource.Get to get the required info here.
		//ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
		//	// Validate patched CustomResource
		//	crPatched := stackInfo.Outputs["crPatched"].(map[string]interface{})
		//	fooV, ok, err := unstructured.NestedString(crPatched, "metadata", "labels", "foo")
		//	assert.True(t, ok)
		//	assert.NoError(t, err)
		//	assert.Equal(t, "foo", fooV)
		//},
	})
	integration.ProgramTest(t, &test)
}
