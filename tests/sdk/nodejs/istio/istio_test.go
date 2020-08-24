// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package ints

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
)

func TestIstio(t *testing.T) {
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir:          "step1",
		Dependencies: []string{"@pulumi/kubernetes"},
		Quick:        true,
		SkipRefresh:  true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			frontend := stackInfo.Outputs["frontendIp"].(string)

			// Retry the GET on the Istio gateway repeatedly. Istio doesn't publish `.status` on any
			// of its CRDs, so this is as reliable as we can be right now.
			for i := 1; i < 10; i++ {
				req, err := http.Get(fmt.Sprintf("http://%s", frontend))
				if err != nil {
					fmt.Printf("Request to Istio gateway failed: %v\n", err)
					time.Sleep(time.Second * 10)
				} else if req.StatusCode == 200 {
					return
				}
			}

			assert.Fail(t, "Maximum Istio gateway request retries exceeded")
		},
	})
}
