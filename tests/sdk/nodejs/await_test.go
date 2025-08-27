package test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optrefresh"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type awaitedResource struct {
	Spec struct {
		SomeField string `json:"someField"`
	} `json:"spec"`
	Status struct {
		ObservedGeneration int `json:"observedGeneration"`
		Conditions         []struct {
			Status string `json:"status"`
			Type   string `json:"type"`
		} `json:"conditions"`
	} `json:"status"`
}

type expectation struct {
	name               string
	someField          string
	conditionType      string
	conditionStatus    string
	observedGeneration int
}

func TestAwaitGeneric(t *testing.T) {
	// No t.Parallel because this touches environment variables.

	waitForCRDs := func(events <-chan events.EngineEvent) {
		// Wait until we see a resource with a "Waiting for readiness" message
		t.Helper()
		t.Log("waiting for a resource to start awaiting")
		for e := range events {
			if e.DiagnosticEvent == nil {
				continue
			}
			msg := e.DiagnosticEvent.Message
			if strings.Contains(msg, "Waiting for") || strings.Contains(msg, "Missing") {
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
		cmd := exec.Command(filepath.Join(dir, "make-progressing.sh")) // #nosec G204
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
		cmd := exec.Command(filepath.Join(dir, "make-ready.sh")) // #nosec G204
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

			if want.someField != "" {
				assert.Equal(t, want.someField, resource.Spec.SomeField, want.name)
			}
			if want.conditionType != "" {
				require.Len(t, resource.Status.Conditions, 1, want.name)
				assert.Equal(t, want.conditionType, resource.Status.Conditions[0].Type, want.name)
				assert.Equal(t, want.conditionStatus, resource.Status.Conditions[0].Status, want.name)
			}
			if want.observedGeneration != 0 {
				assert.Equal(t, want.observedGeneration, resource.Status.ObservedGeneration)
			}
		}
	}

	// assert that the custom waitFor annotation worked.
	assertWaitForResourcesReady := func(t *testing.T, outputs auto.OutputMap) {
		t.Helper()
		expect := []expectation{
			{
				name:            "wantsCondition",
				someField:       "not-needed",
				conditionType:   "Foo",
				conditionStatus: "True",
			},
			{
				name:      "wantsField",
				someField: "foo",
			},
			{
				name:            "wantsFieldAndCondition",
				someField:       "expected",
				conditionType:   "Foo",
				conditionStatus: "True",
			},
		}
		assertExpectations(t, outputs, expect)
	}

	assertAllResourcesReady := func(t *testing.T, outputs auto.OutputMap) {
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
		assertWaitForResourcesReady(t, outputs)
	}

	assertGenericResourceUntouched := func(t *testing.T, outputs auto.OutputMap) {
		t.Helper()
		expect := []expectation{
			{
				name:            "wantsReady",
				someField:       "untouched",
				conditionType:   "Ready",
				conditionStatus: "False",
			},
			{
				name:               "wantsGenerationIncrement",
				someField:          "untouched",
				conditionType:      "Ready",
				conditionStatus:    "True",
				observedGeneration: 1,
			},
		}
		assertExpectations(t, outputs, expect)
	}

	t.Run("enabled", func(t *testing.T) {
		t.Setenv("PULUMI_K8S_AWAIT_ALL", "true")

		test := pulumitest.NewPulumiTest(t,
			"testdata/await/generic",
		)
		dir := test.WorkingDir()
		t.Cleanup(func() {
			test.Destroy(t)
		})

		// Use kubectl to simulate an operator acting on our resources.
		ch := make(chan events.EngineEvent)
		go func() {
			waitForCRDs(ch)
			makeReady(t, dir)
		}()

		// Create
		up := test.Up(t,
			optup.EventStreams(ch),
			optup.ProgressStreams(os.Stdout),
			optup.ErrorProgressStreams(os.Stderr),
		)
		assertAllResourcesReady(t, up.Outputs)

		// Touch our resources and refresh in order to trigger an update later.
		touch(t, dir)

		// Read
		refresh := test.Refresh(t, optrefresh.ProgressStreams(os.Stdout))
		require.NotNil(t, refresh.Summary.ResourceChanges)
		assert.Equal(t, 5, (*refresh.Summary.ResourceChanges)["update"])

		ch = make(chan events.EngineEvent)
		go func() {
			waitForCRDs(ch)
			makeReady(t, dir)
		}()

		// Update
		up = test.Up(t, optup.EventStreams(ch), optup.ProgressStreams(os.Stdout), optup.ErrorProgressStreams(os.Stderr))
		assertAllResourcesReady(t, up.Outputs)
	})

	t.Run("disabled", func(t *testing.T) {
		// With generic await disabled, CustomResources and other types without
		// custom await logic should no-op instead of waiting for readiness.
		// However custom await logic via the waitFor annotation still applies.

		test := pulumitest.NewPulumiTest(t,
			"testdata/await/generic",
		)
		dir := test.WorkingDir()
		t.Cleanup(func() {
			test.Destroy(t)
		})

		// Use kubectl to simulate an operator acting on our resources.
		ch := make(chan events.EngineEvent)
		go func() {
			waitForCRDs(ch)
			makeReady(t, dir)
		}()

		// Create
		up := test.Up(t, optup.EventStreams(ch), optup.ProgressStreams(os.Stdout), optup.ErrorProgressStreams(os.Stderr))
		assertGenericResourceUntouched(t, up.Outputs)
		assertWaitForResourcesReady(t, up.Outputs)

		// Touch our resources and refresh in order to trigger an update later.
		touch(t, dir)

		// Read
		refresh := test.Refresh(t, optrefresh.ProgressStreams(os.Stdout))
		require.NotNil(t, refresh.Summary.ResourceChanges)
		assert.Equal(t, 5, (*refresh.Summary.ResourceChanges)["update"])

		ch = make(chan events.EngineEvent)
		go func() {
			waitForCRDs(ch)
			makeReady(t, dir)
		}()

		// Update
		up = test.Up(t, optup.EventStreams(ch), optup.ProgressStreams(os.Stdout), optup.ErrorProgressStreams(os.Stderr))
		assertGenericResourceUntouched(t, up.Outputs)
		assertWaitForResourcesReady(t, up.Outputs)
	})
}

func TestAwaitCertManager(t *testing.T) {
	t.Setenv("PULUMI_K8S_AWAIT_ALL", "true")

	test := pulumitest.NewPulumiTest(t,
		"testdata/await/cert-manager",
	)
	t.Cleanup(func() {
		test.Destroy(t)
	})

	test.Up(t, optup.ProgressStreams(os.Stdout))
	test.Destroy(t, optdestroy.ProgressStreams(os.Stdout))
}
