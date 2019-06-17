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

	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/stretchr/testify/assert"
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
		Dir:                  filepath.Join(cwd, "yaml-test"),
		ExpectRefreshChanges: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 24, len(stackInfo.Deployment.Resources))

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
			var ns interface{}
			var namespaceName, namespace2Name string

			// Verify CRD.
			crd := stackInfo.Deployment.Resources[0]
			assert.Equal(t, tokens.Type("kubernetes:apiextensions.k8s.io/v1beta1:CustomResourceDefinition"),
				crd.URN.Type())
			name, _ = openapi.Pluck(crd.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "foos.bar.example.com"))

			// Verify CR.
			cr := stackInfo.Deployment.Resources[1]
			assert.Equal(t, tokens.Type("kubernetes:bar.example.com/v1:Foo"), cr.URN.Type())
			name, _ = openapi.Pluck(cr.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "foobar"))

			// Verify namespace1.
			namespace := stackInfo.Deployment.Resources[2]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())
			name, _ = openapi.Pluck(namespace.Outputs, "metadata", "name")
			namespaceName = name.(string)
			assert.True(t, strings.HasPrefix(namespaceName, "ns"), fmt.Sprintf("%s %s", name, "ns"))

			// Verify namespace2.
			namespace2 := stackInfo.Deployment.Resources[3]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace2.URN.Type())
			name, _ = openapi.Pluck(namespace2.Outputs, "metadata", "name")
			namespace2Name = name.(string)
			assert.True(t, strings.HasPrefix(namespace2Name, "ns2"), fmt.Sprintf("%s %s", name, "ns2"))

			// Verify Pod "bar".
			podBar := stackInfo.Deployment.Resources[4]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), podBar.URN.Type())
			name, _ = openapi.Pluck(podBar.Outputs, "metadata", "name")
			assert.True(t, strings.Contains(name.(string), "bar"), fmt.Sprintf("%s %s", name, "bar"))
			ns, _ = openapi.Pluck(podBar.Outputs, "metadata", "namespace")
			assert.Equal(t, ns, namespaceName)

			// Verify Pod "baz".
			podBaz := stackInfo.Deployment.Resources[5]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), podBaz.URN.Type())
			name, _ = openapi.Pluck(podBaz.Outputs, "metadata", "name")
			assert.True(t, strings.Contains(name.(string), "baz"), fmt.Sprintf("%s %s", name, "baz"))
			ns, _ = openapi.Pluck(podBaz.Outputs, "metadata", "namespace")
			assert.Equal(t, ns, namespaceName)

			// Verify Pod "foo".
			podFoo := stackInfo.Deployment.Resources[6]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), podFoo.URN.Type())
			name, _ = openapi.Pluck(podFoo.Outputs, "metadata", "name")
			assert.True(t, strings.Contains(name.(string), "foo"), fmt.Sprintf("%s %s", name, "foo"))
			ns, _ = openapi.Pluck(podFoo.Outputs, "metadata", "namespace")
			assert.Equal(t, ns, namespaceName)

			// Note: Skipping validation for the guestbook app in this test since it's no different from the
			// first ConfigFile.

			// Verify the provider resource.
			provRes := stackInfo.Deployment.Resources[22]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			// Verify root resource.
			stackRes := stackInfo.Deployment.Resources[23]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			// TODO[pulumi/pulumi#2782] Testing of secrets blocked on a bug in Python support for secrets.
			// // Ensure that all `Pod` have `status` marked as a `Secret`
			// for _, res := range stackInfo.Deployment.Resources {
			// 	if res.Type == tokens.Type("kubernetes:core/v1:Pod") {
			// 		spec, has := res.Outputs["apiVersion"]
			// 		assert.True(t, has)
			// 		specMap, is := spec.(map[string]interface{})
			// 		assert.True(t, is)
			// 		sigKey, has := specMap[resource.SigKey]
			// 		assert.True(t, has)
			// 		assert.Equal(t, resource.SecretSig, sigKey)
			// 	}
			// }
		},
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

func TestHelm(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")
	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join(cwd, "helm"),
	})
	integration.ProgramTest(t, &options)
}
