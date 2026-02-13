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
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/require"
)

func TestTerraformConvert(t *testing.T) {
	// Test that we can convert some terraform code to pulumi code with this kubernetes provider.

	// We're going to convert the terraform code in testdata into YAML.
	tmp := t.TempDir()
	cmd := exec.Command("pulumi", "convert", "--from=terraform", "--language=yaml", "--out", tmp)
	cmd.Dir = "testdata/convert"
	cmd.Env = append(os.Environ(), "PATH=../../../bin:"+os.Getenv("PATH"))

	out, err := cmd.CombinedOutput()
	t.Logf("output: %s", out)
	require.NoError(t, err)

	// Check that the output is what we expect.
	expected, err := os.ReadFile("testdata/convert/Main.yaml")
	require.NoError(t, err)
	actual, err := os.ReadFile(filepath.Join(tmp, "Main.yaml"))
	require.NoError(t, err)
	require.Equal(t, string(expected), string(actual))
}
