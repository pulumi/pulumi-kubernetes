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

package ints

// FIXME(levi): Figure out why this test is flaky, and re-enable it in CI.
//  https://github.com/pulumi/pulumi-kubernetes/issues/1016

//import (
//	"os"
//	"testing"
//
//	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
//)
//
//func TestPerformance(t *testing.T) {
//	kubectx := os.Getenv("KUBERNETES_CONTEXT")
//
//	if kubectx == "" {
//		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
//	}
//
//	integration.ProgramTest(t, &integration.ProgramTestOptions{
//		Dir:                  "step1",
//		Dependencies:         []string{"@pulumi/kubernetes"},
//		ExpectRefreshChanges: true, // The Mutating and Validating webhooks update on refresh.
//	})
//}
