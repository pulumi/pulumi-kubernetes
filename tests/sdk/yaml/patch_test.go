// Copyright 2016-2024, Pulumi Corporation.
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

//nolint:goconst // yaml
package test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/pulumi-kubernetes/tests/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

// TestPatchResources tests that patching resources works as expected. The Patch variants of the following resources are
// tested:
// a) Namespace
// b) Deployment
// c) CustomResource (TestResource)
// d) CustomResource with Kind ending in "Patch" (TestResourcePatch)
// We ensure that the Patch variants are properly handled, and deletion of these resources should unset the fields that
// were set by the patch.
// This test currently runs against yaml, golang and nodejs languages since CustomResources are overlays.
func TestPatchResources(t *testing.T) {
	t.Parallel()
	const (
		testFolder                  = "testdata/patch-resources"
		customResourceKindPlain     = "TestPatchResource"
		customResourceKindWithPatch = "TestPatchResourcePatch"
		patchAnnotationField        = "pulumi.com/testPatchAnnotation"
		patchAnnotationValue        = "patched"
	)

	// rsc is a struct used to represent a resource in the cluster.
	type rsc struct {
		name, kind string
	}

	// obj is a struct used to unmarshal the output of `kubectl get` commands to make it easier to assert on the output.
	type obj struct {
		Metadata struct {
			Annotations map[string]string `json:"annotations"`
		} `json:"metadata"`
		Spec struct {
			Foo string `json:"foo"`
		} `json:"spec"`
	}

	// getClusterObj is a helper function to get the object metadata from the cluster using `kubectl`.
	getClusterObj := func(kind, name, namespace string) (*obj, error) {
		outB, err := tests.Kubectl("get", kind, name, "-n", namespace, "-o", "json")
		if err != nil {
			return nil, fmt.Errorf("failed to get %s %s: %w", kind, name, err)
		}
		objOutput := new(obj)
		if err := json.Unmarshal(outB, &objOutput); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s %s output: %w", kind, name, err)
		}
		return objOutput, nil
	}

	// Apply the required CRDs for this test to the cluster. We manage this separately from the Pulumi program since it
	// is cluster scoped,
	// and can be used by all sub-tests.
	_, err := tests.Kubectl("apply", "-f", filepath.Join(testFolder, "crds.yaml"))
	if err != nil {
		t.Fatalf("failed to apply CRDs: %v", err)
	}
	t.Cleanup(func() {
		_, err := tests.Kubectl("delete", "-f", filepath.Join(testFolder, "crds.yaml"))
		contract.AssertNoErrorf(err, "failed to delete CRDs during cleanup")
	})

	// testRunnerF is a is the acutal test runner to run the same test case for every language. This ensures that the
	// CustomResourcePatch overlay resources
	// works in all languages.
	testRunnerF := func(t *testing.T, language string) {
		// 1. Create the resources.
		test := pulumitest.NewPulumiTest(t, filepath.Join(testFolder, language, "step1"))
		if language != "yaml" {
			test.Install(t)
		}
		t.Logf("into %s", test.WorkingDir())
		t.Cleanup(func() {
			test.Destroy(t)
		})

		test.Preview(t)
		outputs := test.Up(t).Outputs
		// Validate the resources do not contain the patch fields. Get the program outputs for use with `kubectl`.
		namespace := outputs["nsName"].Value.(string)
		resources := []rsc{
			{namespace, "namespace"},
			{outputs["depName"].Value.(string), "deployment"},
		}

		// Add custom resources if language is not Yaml.
		if language != "yaml" {
			resources = append(resources,
				rsc{
					outputs["plainCRName"].Value.(string), customResourceKindPlain,
				},
				rsc{
					outputs["patchCRName"].Value.(string), customResourceKindWithPatch,
				})
		}

		for _, resource := range resources {
			objOutput, err := getClusterObj(resource.kind, resource.name, namespace)
			if err != nil {
				t.Errorf("failed to get created object in cluster: %v", err)
				continue
			}

			if annoVal, ok := objOutput.Metadata.Annotations[patchAnnotationField]; ok {
				t.Errorf(
					"expected %s %s to not have annotation %s, but found %q",
					resource.kind,
					resource.name,
					patchAnnotationField,
					annoVal,
				)
			}

			if resource.kind == customResourceKindPlain || resource.kind == customResourceKindWithPatch {
				if objOutput.Spec.Foo != "bar" {
					t.Errorf(
						"expected %s %s to have spec.foo=%q, but found %q",
						resource.kind,
						resource.name,
						"bar",
						objOutput.Spec.Foo,
					)
				}
			}
		}

		if t.Failed() {
			t.FailNow()
		}

		// 2. Patch the resources.
		test.UpdateSource(t, testFolder, language, "step2")
		if language != "yaml" {
			test.Install(t)
		}
		test.Preview(t)
		test.Up(t)

		// Validate the resources contain the patch fields.
		for _, resource := range resources {
			objOutput, err := getClusterObj(resource.kind, resource.name, namespace)
			if err != nil {
				t.Errorf("failed to get patched object in cluster: %v", err)
				continue
			}

			if annoVal, ok := objOutput.Metadata.Annotations[patchAnnotationField]; !ok ||
				annoVal != patchAnnotationValue {
				t.Errorf(
					"expected %s %s to have annotation %s=%q, but found %q",
					resource.kind,
					resource.name,
					patchAnnotationField,
					patchAnnotationValue,
					annoVal,
				)
			}

			if resource.kind == customResourceKindPlain || resource.kind == customResourceKindWithPatch {
				if objOutput.Spec.Foo != "bar" {
					t.Errorf(
						"expected %s %s to have spec.foo=%q, but found %q",
						resource.kind,
						resource.name,
						"bar",
						objOutput.Spec.Foo,
					)
				}
			}
		}

		if t.Failed() {
			t.FailNow()
		}

		// 3. Delete the Patch resources by reverting to the Pulumi program in step 1.
		test.UpdateSource(t, testFolder, language, "step1")
		if language != "yaml" {
			test.Install(t)
		}

		test.Preview(t)
		test.Up(t)

		// Validate the resources do not contain the patch fields, and the object has not been deleted from cluster.
		for _, resource := range resources {
			objOutput, err := getClusterObj(resource.kind, resource.name, namespace)
			if err != nil {
				t.Errorf("failed to get unpatched object in cluster: %v", err)
				continue
			}

			if annoVal, ok := objOutput.Metadata.Annotations[patchAnnotationField]; ok {
				t.Errorf(
					"expected %s %s to not have annotation %s, but found %q",
					resource.kind,
					resource.name,
					patchAnnotationField,
					annoVal,
				)
			}

			if resource.kind == customResourceKindPlain || resource.kind == customResourceKindWithPatch {
				if objOutput.Spec.Foo != "bar" {
					t.Errorf(
						"expected %s %s to have spec.foo=%q, but found %q",
						resource.kind,
						resource.name,
						"bar",
						objOutput.Spec.Foo,
					)
				}
			}
		}
	}

	// Run the test for each language.
	for _, language := range []string{"nodejs", "golang", "yaml"} {
		language := language

		t.Run(language, func(t *testing.T) {
			t.Parallel()
			testRunnerF(t, language)
		})
	}
}
