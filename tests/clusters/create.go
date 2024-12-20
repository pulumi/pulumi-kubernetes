package clusters

import (
	"os"
	"sync"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	_kindClusterStr = "kind"
	_gkeCluster     = "gke"
)

// Ensure that a cluster is available for testing.
//
// If `~/.kube/config` already exists, as is often the case with local
// development, it will be used as the test cluster config. This can speed up
// local development, but you should still periodically reset your cluster to
// avoid cross-test contamination. The `KUBECONFIG` environment variable is
// also respected and can point at an existing path.
//
// If `~/.kube/config` doesn't exist a new KinD cluster will be created. It
// will not be torn down.
//
// TODO: A GKE cluster will be created in the case where we detect we're
// running in a CI release workflow.
var Ensure = sync.OnceValues(func() (func(), error) {
	kubeconfig := "~/.kube/config"
	if env := os.Getenv("KUBECONFIG"); env != "" {
		kubeconfig = env
	}

	abs, _ := homedir.Expand(kubeconfig)
	cfg, err := clientcmd.LoadFromFile(abs)

	// TODO: Temporarily forcing GKE for testing.
	cluster, err := NewGKECluster("pulumi-kubernetes", *cfg)
	return func() { _ = cluster.Delete() }, err

	/*
		if err == nil && len(cfg.Clusters) > 0 {
			log.Printf("Using %q for the test cluster", kubeconfig)
			return nil, nil
		}

		if os.Getenv("GITHUB_EVENT_NAME") == "push" {
			cluster, err := NewGKECluster("pulumi-kubernetes", *cfg)
			return func() { _ = cluster.Delete() }, err
		}

		_, err = NewKindCluster("pulumi-kubernetes")
		return nil, err
	*/
})
