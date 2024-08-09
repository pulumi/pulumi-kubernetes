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

func TestSkipAwait(t *testing.T) {
	t.Parallel()

	test := pulumitest.NewPulumiTest(t,
		"testdata/await/skipawait",
		opttest.SkipInstall(),
	)

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
