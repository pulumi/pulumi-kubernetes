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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/pulumi-kubernetes/tests/v4"
)

// TestParameterizeGenSDK runs `pulumi package gen-sdk` against the dev
// provider with a CRD fixture and asserts the generated Go SDK declares the
// parameterized module path — not the base provider's. Uses gen-sdk rather
// than `pulumi install` because it takes the binary path directly and skips
// pulumi's plugin-catalog resolution.
//
// Requires `make k8sprovider`; will fail when gen-sdk tries to exec the
// missing binary.
func TestParameterizeGenSDK(t *testing.T) {
	tests.SkipIfShort(t, "shells out to pulumi CLI")

	providerBin, err := filepath.Abs("../../../bin/pulumi-resource-kubernetes")
	require.NoError(t, err)

	crdPath, err := filepath.Abs("parameterize-install/gateway-crd.yaml")
	require.NoError(t, err)

	tmp := t.TempDir()
	sdkOut := filepath.Join(tmp, "sdks")
	cmd := exec.Command("pulumi", "package", "gen-sdk", providerBin,
		"--language", "go",
		"--out", sdkOut,
		"--local",
		"--", "-v", "1.0.0", "-c", crdPath)
	cmd.Dir = tmp
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "pulumi package gen-sdk: %s", out)

	// Layout gen-sdk produces for Go with package name "gateway-pulumi-test":
	//   sdks/go/                             <- go.mod (module root)
	//   sdks/go/gatewaypulumitest/           <- Go package root (hyphens stripped)
	//   sdks/go/gatewaypulumitest/gateway/v1 <- generated resources
	goRoot := filepath.Join(sdkOut, "go")
	gatewayGo := filepath.Join(goRoot, "gatewaypulumitest", "gateway", "v1", "gateway.go")

	gomod, err := os.ReadFile(filepath.Join(goRoot, "go.mod"))
	require.NoError(t, err)
	assert.Contains(t, string(gomod),
		"module github.com/pulumi/pulumi-kubernetes-gateway-pulumi-test/sdk/v4/go",
		"go.mod must use the parameterized module path")
	assert.NotContains(t, string(gomod),
		"module github.com/pulumi/pulumi-kubernetes/sdk/v4/go\n",
		"go.mod must not reuse the base provider module path")

	gwBytes, err := os.ReadFile(gatewayGo)
	require.NoError(t, err)
	gw := string(gwBytes)
	assert.Contains(t, gw,
		"github.com/pulumi/pulumi-kubernetes-gateway-pulumi-test/sdk/v4/go/gatewaypulumitest",
		"gateway.go must import from the parameterized SDK path")
	for _, line := range strings.Split(gw, "\n") {
		if strings.Contains(line, `"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/`) {
			t.Errorf("gateway.go leaks base-provider import path: %q", strings.TrimSpace(line))
		}
	}
}
