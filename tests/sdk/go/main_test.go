package test

import (
	"testing"

	"github.com/pulumi/pulumi-kubernetes/tests/v4/clusters"
)

func TestMain(m *testing.M) {
	clusters.RunWithClusterCreation(m, "nodejs-test")
}
