package clusters

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"sigs.k8s.io/kind/pkg/cluster"
)

// KindVersion is a specific Kind version associated with a Kubernetes minor version.
type KindVersion string

const (
	// Kind versions and their associated image tags.
	Kind1_29 KindVersion = "kindest/node:v1.29.2"
	Kind1_28 KindVersion = "kindest/node:v1.28.7"
	Kind1_27 KindVersion = "kindest/node:v1.27.11"

	KubeconfigFilename = "KUBECONFIG"

	// maxCreateTries is the maximum number of times to try to create a Kind cluster.
	maxCreateTries = 5
)

type KindCluster struct {
	name           string
	kubeconfigPath string
	teardownFn     func() error
}

func (c KindCluster) Name() string {
	return c.name
}

func (c KindCluster) KubeconfigPath() string {
	return c.kubeconfigPath
}

func (c KindCluster) Connect() error {
	return nil
}

func (c KindCluster) Delete() error {
	return c.teardownFn()
}

// newKindCluster creates a new Kind cluster for use in testing with the specified name. The version
// parameter is variadic and defaults to the latest version of Kind if not provided.
func NewKindCluster(name string, version ...KindVersion) (Cluster, error) {
	v := Kind1_29
	if len(version) > 0 {
		v = version[0]
	}

	// Create a new tmpdir to store the kubeconfig.
	log.Printf("Creating new Kind cluster with image %q", v)
	tmpDir, err := os.MkdirTemp("", "pulumi-test-kind")
	if err != nil {
		return nil, err
	}

	KubeconfigPath := filepath.Join(tmpDir, KubeconfigFilename)

	p := cluster.NewProvider()

	name = normalizeName(name + "-" + randString())
	log.Printf("Creating Kind cluster %q", name)
	teardownFn, err := createKindCluster(p, name, tmpDir, KubeconfigPath, v)
	if err != nil {
		return nil, err
	}

	return &KindCluster{name: name, kubeconfigPath: KubeconfigPath, teardownFn: teardownFn}, nil

}

// createKindCluster attempts to create a KinD cluster with retry.
func createKindCluster(p *cluster.Provider, name, tmpDir, kcfgPath string, version KindVersion) (func() error, error) {
	var err error
	for i := 0; i < maxCreateTries; i++ {
		log.Printf("Creating Kind cluster, attempt: %d", i+1)
		err = p.Create(name,
			cluster.CreateWithKubeconfigPath(kcfgPath),
			cluster.CreateWithNodeImage(string(version)),
			cluster.CreateWithWaitForReady(20*time.Second),
		)
		if err == nil {
			return func() error {
				if err := deleteKindCluster(p, name, kcfgPath); err != nil {
					return fmt.Errorf("unable to delete Kind cluster %q: %v", name, err)
				}

				return os.RemoveAll(tmpDir)
			}, nil
		}
	}

	// We failed to create the cluster maxCreateTries times, so return failure error.
	return nil, err
}

// deleteKindCluster deletes a specified kind cluster.
func deleteKindCluster(p *cluster.Provider, name, kcfgPath string) error {
	return p.Delete(name, kcfgPath)
}
