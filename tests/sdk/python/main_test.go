package python

import (
	"testing"

	"github.com/pulumi/pulumi-kubernetes/tests/v4/clusters"
)

var testClusters = new(clusters.TestClusters)

func TestMain(m *testing.M) {
	clusters.RunWithClusterCreation(m, "python-test", testClusters)
}
