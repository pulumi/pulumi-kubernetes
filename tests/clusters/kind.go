package cluster

import (
	"crypto/rand"
	"encoding/base32"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
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

// newKindCluster creates a new Kind cluster for use in testing with the specified name. The version
// parameter is variadic and defaults to the latest version of Kind if not provided.
func NewKindCluster(t *testing.T, name string, version ...KindVersion) string {
	t.Helper()

	v := Kind1_29
	if len(version) > 0 {
		v = version[0]
	}

	// Create a new tmpdir to store the kubeconfig.
	t.Logf("Creating new Kind cluster with image %q", v)
	tmpDir, err := os.MkdirTemp("", "pulumi-test-kind")
	if err != nil {
		t.Fatalf("failed to create tmpdir: %v", err)
	}
	t.Cleanup(func() {
		contract.Assertf(os.RemoveAll(tmpDir) == nil, "failed to remove tmpdir %q", tmpDir)
	})

	KubeconfigPath := filepath.Join(tmpDir, KubeconfigFilename)

	p := cluster.NewProvider()

	name = normalizeName(name + "-" + randString())
	t.Logf("Creating Kind cluster %q", name)
	teardownFn, err := createKindCluster(t, p, name, KubeconfigPath, v)
	if err != nil {
		t.Fatalf("failed to create Kind cluster: %v", err)
	}

	t.Cleanup(teardownFn)

	// return restConfig(t, KubeconfigPath), shutdown
	return KubeconfigPath

}

// normalizeName returns a normalized name for the cluster that adheres
// to the Kubernetes naming restrictions.
func normalizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, ".", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Shorten the name if necessary
	if len(name) > 63 {
		name = name[:63]
	}

	return name
}

// randString returns a random string of length 6.
func randString() string {
	c := 10
	b := make([]byte, c)
	rand.Read(b)
	length := 6
	return strings.ToLower(base32.StdEncoding.EncodeToString(b)[:length])
}

// createKindCluster attempts to create a Kind cluster with retry.
func createKindCluster(t *testing.T, p *cluster.Provider, name, kcfgPath string, version KindVersion) (func(), error) {
	var err error
	for i := 0; i < maxCreateTries; i++ {
		t.Logf("Creating Kind cluster, attempt: %d", i+1)
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

	// We failed to create the cluster maxCreateTries times, so return failure error.
	return nil, err
}

// deleteKindCluster deletes a specified kind cluster.
func deleteKindCluster(p *cluster.Provider, name, kcfgPath string) error {
	return p.Delete(name, kcfgPath)
}
