// Copyright 2024, Pulumi Corporation.
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
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAwaitDaemonSet(t *testing.T) {
	t.Parallel()

	assertReady := func(t *testing.T, outputs auto.OutputMap) {
		cns, ok := outputs["currentNumberScheduled"]
		require.True(t, ok)

		dns, ok := outputs["desiredNumberScheduled"]
		require.True(t, ok)

		nms, ok := outputs["numberMisscheduled"]
		require.True(t, ok)

		nr, ok := outputs["numberReady"]
		require.True(t, ok)

		assert.Greater(t, cns.Value.(float64), float64(0))
		assert.Equal(t, float64(0), nms.Value.(float64))
		assert.Equal(t, dns, cns)
		assert.Equal(t, cns, nr)
	}

	test := pulumitest.NewPulumiTest(t,
		"testdata/await/daemonset",
		opttest.SkipInstall(),
	)
	t.Cleanup(func() {
		test.Destroy()
	})

	// Create a new DS that takes a few seconds to become ready.
	up := test.Up()
	t.Log(up.Summary.Message)
	assertReady(t, up.Outputs)

	test.Refresh() // Exercise read-await logic.

	// Update the DS to use a different but valid image tag.
	test.UpdateSource("testdata/await/daemonset/step2")
	up = test.Up()
	assertReady(t, up.Outputs)

	// Update the DS to use an invalid image tag. It should never become ready.
	test.UpdateSource("testdata/await/daemonset/step3")
	_, err := test.CurrentStack().Up(context.Background())
	assert.ErrorContains(t, err, `the Kubernetes API server reported that "default/await-daemonset" failed to fully initialize or become live: timed out waiting for the condition`)
}

func TestAwaitPVC(t *testing.T) {
	t.Parallel()

	test := pulumitest.NewPulumiTest(t,
		"testdata/await/pvc",
		opttest.SkipInstall(),
	)
	t.Cleanup(func() {
		test.Destroy()
	})

	// WaitUntilFirstConsumer PVC should still be Pending.
	up := test.Up()
	assert.Equal(t, "Pending", up.Outputs["status"].Value)

	// Adding a Deployment to consume the PVC should succeed.
	test.UpdateSource("testdata/await/pvc/step2")
	up = test.Up()
}

func TestAwaitSkip(t *testing.T) {
	t.Parallel()

	test := pulumitest.NewPulumiTest(t,
		"testdata/await/skipawait",
		opttest.SkipInstall(),
	)
	t.Cleanup(func() {
		test.Destroy()
	})

	start := time.Now()
	_ = test.Up()
	took := time.Since(start)
	assert.Less(t, took, 2*time.Minute, "didn't skip pod's slow startup")

	start = time.Now()
	_ = test.Refresh()
	took = time.Since(start)
	assert.Less(t, took, 2*time.Minute, "didn't skip pod's slow read")

	test.UpdateSource("testdata/await/skipawait/step2")
	start = time.Now()
	_ = test.Refresh()
	took = time.Since(start)
	assert.Less(t, took, 2*time.Minute, "didn't skip pod's slow update")

	start = time.Now()
	_ = test.Destroy()
	took = time.Since(start)
	assert.Less(t, took, 2*time.Minute, "didn't skip config map's stuck delete")
}

type awaitedResource struct {
	Spec struct {
		SomeField string `json:"someField"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			Status string `json:"status"`
			Type   string `json:"type"`
		} `json:"conditions"`
	} `json:"status"`
}

type expectation struct {
	name            string
	someField       string
	conditionType   string
	conditionStatus string
}

// With generic wait enabled & disabled
func TestAwaitGeneric(t *testing.T) {
	t.Parallel()

	waitForCRDs := func() {
		for {
			time.Sleep(1 * time.Second)
			cmd := exec.Command("kubectl", "get", "crd/genericawaiters.test.pulumi.com")
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err == nil {
				break
			}
		}
		// Wait another second for CRs.
		time.Sleep(1 * time.Second)
	}

	touch := func(t *testing.T, dir string) {
		cmd := exec.Command(filepath.Join(dir, "make-progressing.sh"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		require.NoError(t, err)
	}

	makeReady := func(t *testing.T, dir string) {
		cmd := exec.Command(filepath.Join(dir, "make-ready.sh"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		require.NoError(t, err)
	}

	assertExpectations := func(t *testing.T, outputs auto.OutputMap, expect []expectation) {
		for _, want := range expect {
			t.Run(want.name, func(t *testing.T) {
				obj := outputs[want.name].Value
				bytes, err := json.Marshal(obj)
				require.NoError(t, err)

				var resource awaitedResource
				err = json.Unmarshal(bytes, &resource)
				require.NoError(t, err)

				assert.Equal(t, want.someField, resource.Spec.SomeField)
				if want.conditionType != "" {
					require.Len(t, resource.Status.Conditions, 1)
					assert.Equal(t, want.conditionType, resource.Status.Conditions[0].Type)
					assert.Equal(t, want.conditionStatus, resource.Status.Conditions[0].Status)
				}
			})
		}
	}

	assertReady := func(t *testing.T, outputs auto.OutputMap) {
		expect := []expectation{
			{
				name:            "wantsReady",
				someField:       "",
				conditionType:   "Ready",
				conditionStatus: "True",
			},
		}
		assertExpectations(t, outputs, expect)
	}

	assertUntouched := func(t *testing.T, outputs auto.OutputMap) {
		expect := []expectation{
			{
				name:            "wantsReady",
				someField:       "",
				conditionType:   "Ready",
				conditionStatus: "False",
			},
		}
		assertExpectations(t, outputs, expect)
	}

	t.Run("enabled", func(t *testing.T) {
		test := pulumitest.NewPulumiTest(t,
			"testdata/await/generic",
			opttest.SkipInstall(),
		)
		dir := test.Source()
		t.Cleanup(func() {
			test.Destroy()
		})
		test.Install()

		// Simulate an operator acting on our resources.
		go func() {
			waitForCRDs()

			// First apply some unrelated changes.
			touch(t, dir)

			time.Sleep(2 * time.Second)

			// Now make the resources ready.
			makeReady(t, dir)
		}()

		up := test.Up()
		assertReady(t, up.Outputs)
	})

	t.Run("disabled", func(t *testing.T) {
		// With generic await disabled, CustomResources and other types without
		// custom await logic should no-op instead of waiting for readiness.

		test := pulumitest.NewPulumiTest(t,
			"testdata/await/generic",
			opttest.SkipInstall(),
		)
		dir := test.Source()
		t.Cleanup(func() {
			test.Destroy()
		})
		test.Install()

		// Simulate an operator acting on our resources.
		go func() {
			waitForCRDs()
			// Apply some unrelated changes -- our update should already be
			// finished, so this shouldn't impact our stack.
			touch(t, dir)
		}()

		up := test.Up()
		assertUntouched(t, up.Outputs)
	})
}
