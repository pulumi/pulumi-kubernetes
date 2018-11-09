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
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/testing/integration"
)

var baseOptions = &integration.ProgramTestOptions{
	Dependencies: []string{
		filepath.Join("..", "..", "sdk", "python", "bin"),
	},
}

func TestSmoke(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")
	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join("python", "smoke-test"),
	})
	integration.ProgramTest(t, &options)
}
