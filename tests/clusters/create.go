package clusters

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

const (
	KindClusterStr = "kind"
	GKECluster     = "gke"
)

var (
	TestClusterList = new(TestClusters)
)

// RunWithClusterCreation is the entry point for the tests and provides a way to create and delete clusters for the tests.
// It is able to create a pool of clusters to randomly run tests on.
func RunWithClusterCreation(m *testing.M, clusterPrefix string) {
	createCluster := flag.Bool("create-cluster", false, "Create a cluster for the tests, default is false and to use an existing cluster")
	clusterType := flag.String("cluster-type", KindClusterStr, "The type of cluster to create for the tests, default is kind")
	numClusters := flag.Int("num-clusters", 1, "The number of clusters to create for the tests, default is 1")
	flag.Parse()

	if *createCluster {
		makeCluster(clusterPrefix, TestClusterList, *clusterType, *numClusters)
	}

	exitCode := m.Run()

	// Always attempt to tear down all the clusters after the tests have run.
	var wg sync.WaitGroup
	for _, cluster := range TestClusterList.clusters {
		cluster := cluster
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := cluster.Delete(); err != nil {
				fmt.Println("Unable to delete cluster: ", err)
				exitCode = 1
			}
		}()
	}
	wg.Wait()

	os.Exit(exitCode)
}

func makeCluster(clusterPrefix string, testClusters *TestClusters, clusterType string, numClusters int) {
	// Make the slice with the number of clusters to create so we do not need a mutex to handle write access to the slice.
	testClusters.clusters = make([]Cluster, numClusters)

	var wg sync.WaitGroup
	for i := 0; i < numClusters; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()

			var err error
			var cluster Cluster
			switch clusterType {
			case KindClusterStr:
				cluster, err = NewKindCluster(clusterPrefix)
				if err != nil {
					fmt.Println("Unable to create Kind cluster: ", err)
					os.Exit(1)
				}
			case GKECluster:
				fmt.Println("STILL TODO")
				os.Exit(1)
			default:
				fmt.Println("Unrecognized cluster type", clusterType)
				os.Exit(1)
			}

			testClusters.clusters[i] = cluster
		}()
	}
	wg.Wait()
}

type TestClusters struct {
	clusters []Cluster
	picked   int
}

// pickCluster picks a cluster to use in sequence (ie. round robbin) from the list of clusters created in TestMain.
// This ensures that we don't an uneven number of tests to run on each cluster.
func (t *TestClusters) PickCluster() Cluster {
	t.picked++
	return t.clusters[t.picked%len(t.clusters)]
}

// WrapProviderTestOptions is a convenience function that wraps the provider test options with the kubeconfig of the cluster to use.
// The cluster's kubeconfig is also returned for test cases that manually shell out to kubectl for validation purposes. This kubeconfig is
// required for the kubectl command to target the right cluster that the test is running on.
func (t *TestClusters) WrapProviderTestOptions(opts integration.ProgramTestOptions) (integration.ProgramTestOptions, string) {
	if t.clusters == nil {
		// Check if KUBECONFIG is set, if it is then use that for the tests.
		kcfg := os.Getenv("KUBECONFIG")
		if kcfg == "" {
			kcfg = os.ExpandEnv("$HOME/.kube/config")
		}
		return opts, kcfg
	}

	cluster := t.PickCluster()

	return opts.With(integration.ProgramTestOptions{
		Env: []string{"KUBECONFIG=" + cluster.KubeconfigPath()},
	}), cluster.KubeconfigPath()
}
