package cluster

import (
	"encoding/base32"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"sigs.k8s.io/kind/pkg/cluster"
)

// KindVersion is a specific Kind version associated with a Kubernetes minor version.
type KindVersion string

const (
	Kind1_27 KindVersion = "kindest/node:v1.27.3"

	// Kubeconfig is the filename of the KUBECONFIG file.
	Kubeconfig = "KUBECONFIG"

	// maxCreateTries is the maximum number of times to try to create a Kind cluster.
	maxCreateTries = 6
)

// newKindCluster creates a new Kind cluster for use in testing with the specified name.
// func NewKindCluster(t *testing.T, name string) (*rest.Config, func()) {
func NewKindCluster(t *testing.T, name string) string {
	t.Helper()

	// Create a new tmpdir to store the kubeconfig.
	t.Log("Creating new Kind cluster")
	tmpDir, err := os.MkdirTemp("", "pulumi-test-kind")
	if err != nil {
		t.Fatalf("failed to create tmpdir: %v", err)
	}
	KubeconfigPath := filepath.Join(tmpDir, Kubeconfig)
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	t.Logf("KubeconfigPath: %s\n", KubeconfigPath)

	p := cluster.NewProvider()

	name += "-" + randString()
	name = normalizeName(name)
	t.Logf("Creating cluster %q", name)
	shutdown, err := createKindCluster(t, p, name, KubeconfigPath, Kind1_27)
	if err != nil {
		t.Fatalf("failed to create kind cluster: %v", err)
	}

	t.Cleanup(shutdown)

	// return restConfig(t, KubeconfigPath), shutdown
	return KubeconfigPath

}

// normalizeName returns a normalized name for the cluster that adheres
// to the Kubernetes naming restrictions.
func normalizeName(name string) string {
	return strings.ToLower(name)
}

func randString() string {
	rand.Seed(time.Now().UnixNano())
	c := 10
	b := make([]byte, c)
	rand.Read(b)
	length := 6
	return strings.ToLower(base32.StdEncoding.EncodeToString(b)[:length])
}

// createKindCluster attempts to create a KinD cluster with retry.
func createKindCluster(t *testing.T, p *cluster.Provider, name, kcfgPath string, version KindVersion) (func(), error) {
	var err error
	for i := 0; i < maxCreateTries; i++ {
		t.Logf("Creating KinD cluster, attempt: %d", i+1)
		err = p.Create(name,
			cluster.CreateWithKubeconfigPath(kcfgPath),
			cluster.CreateWithNodeImage(string(version)),
			cluster.CreateWithWaitForReady(20*time.Second),
		)
		if err == nil {
			return func() {
				deleteKindCluster(p, name, kcfgPath)
			}, nil
		}
	}

	// We failed to create the cluster maxKindTries times, so fail out.
	return nil, err
}

// deleteKindCluster deletes a specified kind cluster.
func deleteKindCluster(p *cluster.Provider, name, kcfgPath string) error {
	return p.Delete(name, kcfgPath)
}
