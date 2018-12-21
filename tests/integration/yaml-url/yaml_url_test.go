// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package ints

import (
	"os"
	"testing"

	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestYAMLURL(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")

	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir:          "step1",
		Dependencies: []string{"@pulumi/kubernetes"},
		Quick:        true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)

			// Assert that we've retrieved the YAML from the URL and provisioned them.
			assert.Equal(t, 9, len(stackInfo.Deployment.Resources))
		},
	})
}
