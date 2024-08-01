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
	"strings"
	"testing"
	"time"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optrefresh"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
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
	_ = test.Up(optup.ProgressStreams(os.Stdout))
	took := time.Since(start)
	assert.Less(t, took, 2*time.Minute, "didn't skip pod's slow startup")

	start = time.Now()
	_ = test.Refresh(optrefresh.ProgressStreams(os.Stdout))
	took = time.Since(start)
	assert.Less(t, took, 2*time.Minute, "didn't skip pod's slow read")

	test.UpdateSource("testdata/await/skipawait/step2")
	start = time.Now()
	_ = test.Refresh(optrefresh.ProgressStreams(os.Stdout))
	took = time.Since(start)
	assert.Less(t, took, 2*time.Minute, "didn't skip pod's slow update")

	start = time.Now()
	_ = test.Destroy(optdestroy.ProgressStreams(os.Stdout))
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

func TestAwaitGeneric(t *testing.T) {
	// No t.Parallel because this touches environment variables.

	waitForCRDs := func(events <-chan events.EngineEvent) {
		// Wait until we see a resource with a "Waiting for readiness message"
		t.Helper()
		t.Log("waiting for a resource to start awaiting")
		for e := range events {
			if e.DiagnosticEvent == nil {
				continue
			}
			if strings.Contains(e.DiagnosticEvent.Message, "Waiting for readiness") {
				go func() {
					for range events {
						// Need to exhaust the channel otherwise things deadlock.
					}
				}()
				break
			}
		}
		// Wait an extra second to let any other resources to get applied.
		time.Sleep(1 * time.Second)
	}

	touch := func(t *testing.T, dir string) {
		t.Helper()
		t.Log("touching resources")
		cmd := exec.Command(filepath.Join(dir, "make-progressing.sh"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		require.NoError(t, err)
	}

	makeReady := func(t *testing.T, dir string) {
		touch(t, dir)
		time.Sleep(1 * time.Second)
		t.Log("marking resources ready")
		t.Helper()
		cmd := exec.Command(filepath.Join(dir, "make-ready.sh"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		require.NoError(t, err)
	}

	assertExpectations := func(t *testing.T, outputs auto.OutputMap, expect []expectation) {
		t.Helper()
		for _, want := range expect {
			obj := outputs[want.name].Value
			bytes, err := json.Marshal(obj)
			require.NoError(t, err)

			var resource awaitedResource
			err = json.Unmarshal(bytes, &resource)
			require.NoError(t, err)

			assert.Equal(t, want.someField, resource.Spec.SomeField, want.name)
			if want.conditionType != "" {
				require.Len(t, resource.Status.Conditions, 1, want.name)
				assert.Equal(t, want.conditionType, resource.Status.Conditions[0].Type, want.name)
				assert.Equal(t, want.conditionStatus, resource.Status.Conditions[0].Status, want.name)
			}
		}
	}

	assertReady := func(t *testing.T, outputs auto.OutputMap) {
		t.Helper()
		expect := []expectation{
			{
				name:            "wantsReady",
				someField:       "touched",
				conditionType:   "Ready",
				conditionStatus: "True",
			},
		}
		assertExpectations(t, outputs, expect)
	}

	assertUntouched := func(t *testing.T, outputs auto.OutputMap) {
		t.Helper()
		expect := []expectation{
			{
				name:            "wantsReady",
				someField:       "untouched",
				conditionType:   "Ready",
				conditionStatus: "False",
			},
		}
		assertExpectations(t, outputs, expect)
	}

	t.Run("enabled", func(t *testing.T) {
		t.Setenv("PULUMI_K8S_AWAIT_ALL", "true")

		test := pulumitest.NewPulumiTest(t,
			"testdata/await/generic",
			opttest.SkipInstall(),
		)
		dir := test.Source()
		t.Cleanup(func() {
			test.Destroy()
		})
		test.Install()

		// Use kubectl to simulate an operator acting on our resources.
		ch := make(chan events.EngineEvent)
		go func() {
			waitForCRDs(ch)
			makeReady(t, dir)
		}()

		// Create
		up := test.Up(optup.EventStreams(ch), optup.ProgressStreams(os.Stdout), optup.ErrorProgressStreams(os.Stderr))
		assertReady(t, up.Outputs)

		// Touch our resources and refresh in order to trigger an update later.
		touch(t, dir)

		// Read
		test.Refresh(optrefresh.ProgressStreams(os.Stdout))

		ch = make(chan events.EngineEvent)
		go func() {
			waitForCRDs(ch)
			makeReady(t, dir)
		}()

		// Update
		up = test.Up(optup.EventStreams(ch), optup.ProgressStreams(os.Stdout))
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

		// Create should return immediately.
		up := test.Up(optup.ProgressStreams(os.Stdout))
		assertUntouched(t, up.Outputs)

		// Touch the resources and refresh to pick up the new state.
		touch(t, dir)

		// Read
		refresh := test.Refresh(optrefresh.ProgressStreams(os.Stdout))
		require.NotNil(t, refresh.Summary.ResourceChanges)
		assert.Equal(t, (*refresh.Summary.ResourceChanges)["update"], 1)

		// Update should exit immediately and reflect the inputs again.
		up = test.Up(optup.ProgressStreams(os.Stdout))
		assertUntouched(t, up.Outputs)
	})
}
