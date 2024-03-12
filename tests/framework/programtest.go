package testframework

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi-kubernetes/tests/v4/clusters"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

// Wrapper for ProgramTest and ProviderTest to run against a specified cluster.
type IntegrationTest struct {
	*integration.ProgramTestOptions
	*pulumitest.PulumiTest

	kubeconfigPath string
}

func (it *IntegrationTest) KubeconfigPath() string {
	return it.kubeconfigPath
}

func NewProgramTest(baseOptions *integration.ProgramTestOptions) *IntegrationTest {
	cluster := clusters.TestClusterList.PickCluster()
	withCluster := baseOptions.With(integration.ProgramTestOptions{
		Env: []string{"KUBECONFIG=" + cluster.KubeconfigPath()},
	})

	return &IntegrationTest{
		ProgramTestOptions: &withCluster,
		kubeconfigPath:     cluster.KubeconfigPath(),
	}
}

func NewProviderTest(t *testing.T, source string, opts ...opttest.Option) *pulumitest.PulumiTest {
	cluster := clusters.TestClusterList.PickCluster()
	opts = append([]opttest.Option{opttest.Env("KUBECONFIG", cluster.KubeconfigPath())}, opts...)
	return pulumitest.NewPulumiTest(t, source, opts...)
}
