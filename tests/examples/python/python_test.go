// Copyright 2016-2018, Pulumi Corporation.
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
package python

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/pulumi/pulumi/pkg/tokens"
)

var baseOptions = &integration.ProgramTestOptions{
	Verbose: true,
	Dependencies: []string{
		filepath.Join("..", "..", "..", "sdk", "python", "bin"),
	},
}

func TestSmoke(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")
	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join(cwd, "smoke-test"),
	})
	integration.ProgramTest(t, &options)
}

func TestYaml(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")
	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join(cwd, "yaml-test"),
	})
	integration.ProgramTest(t, &options)
}

func TestGuestbook(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")
	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join(cwd, "guestbook"),
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 9, len(stackInfo.Deployment.Resources))

			sort.Slice(stackInfo.Deployment.Resources, func(i, j int) bool {
				ri := stackInfo.Deployment.Resources[i]
				rj := stackInfo.Deployment.Resources[j]
				riname, _ := openapi.Pluck(ri.Outputs, "metadata", "name")
				rinamespace, _ := openapi.Pluck(ri.Outputs, "metadata", "namespace")
				rjname, _ := openapi.Pluck(rj.Outputs, "metadata", "name")
				rjnamespace, _ := openapi.Pluck(rj.Outputs, "metadata", "namespace")
				return fmt.Sprintf("%s/%s/%s", ri.URN.Type(), rinamespace, riname) <
					fmt.Sprintf("%s/%s/%s", rj.URN.Type(), rjnamespace, rjname)
			})

			var name interface{}
			var status interface{}

			// Verify frontend deployment.
			frontendDepl := stackInfo.Deployment.Resources[0]
			assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), frontendDepl.URN.Type())
			name, _ = openapi.Pluck(frontendDepl.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "frontend"))
			status, _ = openapi.Pluck(frontendDepl.Outputs, "status", "readyReplicas")
			assert.Equal(t, float64(3), status)

			// Verify redis-follower deployment.
			redisFollowerDepl := stackInfo.Deployment.Resources[1]
			assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), redisFollowerDepl.URN.Type())
			name, _ = openapi.Pluck(redisFollowerDepl.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "redis-follower"), fmt.Sprintf("%s %s", name, "redis-slave"))
			status, _ = openapi.Pluck(redisFollowerDepl.Outputs, "status", "readyReplicas")
			assert.Equal(t, float64(1), status)

			// Verify redis-leader deployment.
			redisLeaderDepl := stackInfo.Deployment.Resources[2]
			assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), redisLeaderDepl.URN.Type())
			name, _ = openapi.Pluck(redisLeaderDepl.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "redis-leader"), fmt.Sprintf("%s %s", name, "redis-master"))
			status, _ = openapi.Pluck(redisLeaderDepl.Outputs, "status", "readyReplicas")
			assert.Equal(t, float64(1), status)

			// Verify test namespace.
			namespace := stackInfo.Deployment.Resources[3]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())
			name, _ = openapi.Pluck(namespace.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "test"), fmt.Sprintf("%s %s", name, "test"))

			// Verify frontend service.
			frontendService := stackInfo.Deployment.Resources[4]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), frontendService.URN.Type())
			name, _ = openapi.Pluck(frontendService.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "frontend"), fmt.Sprintf("%s %s", name, "frontend"))
			status, _ = openapi.Pluck(frontendService.Outputs, "spec", "clusterIP")
			assert.True(t, len(status.(string)) > 1)

			// Verify redis-follower service.
			redisFollowerService := stackInfo.Deployment.Resources[5]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), redisFollowerService.URN.Type())
			name, _ = openapi.Pluck(redisFollowerService.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "redis-follower"), fmt.Sprintf("%s %s", name, "redis-slave"))
			status, _ = openapi.Pluck(redisFollowerService.Outputs, "spec", "clusterIP")
			assert.True(t, len(status.(string)) > 1)

			// Verify redis-leader service.
			redisLeaderService := stackInfo.Deployment.Resources[6]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), redisLeaderService.URN.Type())
			name, _ = openapi.Pluck(redisLeaderService.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "redis-leader"), fmt.Sprintf("%s %s", name, "redis-master"))
			status, _ = openapi.Pluck(redisLeaderService.Outputs, "spec", "clusterIP")
			assert.True(t, len(status.(string)) > 1)

			// Verify the provider resource.
			provRes := stackInfo.Deployment.Resources[7]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			// Verify root resource.
			stackRes := stackInfo.Deployment.Resources[8]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
		},
	})
	integration.ProgramTest(t, &options)
}

// Smoke test for first-class Kubernetes providers.
func TestProvider(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")
	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join(cwd, "provider"),
	})
	integration.ProgramTest(t, &options)
}
