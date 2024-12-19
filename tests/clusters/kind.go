package clusters

import (
	"log"
	"time"

	"sigs.k8s.io/kind/pkg/cluster"
)

// kindVersion is a specific Kind version associated with a Kubernetes minor version.
type kindVersion string

const (
	// Kind versions and their associated image tags.
	_kind1_32 kindVersion = "kindest/node:v1.32.0"
	_kind1_31 kindVersion = "kindest/node:v1.31.4"
	_kind1_30 kindVersion = "kindest/node:v1.30.8"
	_kind1_29 kindVersion = "kindest/node:v1.29.2"
	_kind1_28 kindVersion = "kindest/node:v1.28.7"
	_kind1_27 kindVersion = "kindest/node:v1.27.11"

	// _maxCreateTries is the maximum number of times to try to create a Kind cluster.
	_maxCreateTries = 5
)

type KindCluster struct {
	name       string
	teardownFn func() error
}

func (c KindCluster) Name() string {
	return c.name
}

func (c KindCluster) Connect() error {
	return nil
}

func (c KindCluster) Delete() error {
	return c.teardownFn()
}

// newKindCluster creates a new Kind cluster for use in testing with the specified name. The version
// parameter is variadic and defaults to the latest version of Kind if not provided.
func NewKindCluster(name string, version ...kindVersion) (Cluster, error) {
	v := _kind1_32
	if len(version) > 0 {
		v = version[0]
	}

	p := cluster.NewProvider()

	name = normalizeName(name + "-" + randString())
	log.Printf("Creating Kind cluster %q", name)
	teardownFn, err := createKindCluster(p, name, v)
	if err != nil {
		return nil, err
	}

	return &KindCluster{name: name, teardownFn: teardownFn}, nil

}

// createKindCluster attempts to create a KinD cluster with retry.
func createKindCluster(p *cluster.Provider, name string, version kindVersion) (func() error, error) {
	var err error
	for i := 0; i < _maxCreateTries; i++ {
		log.Printf("Creating Kind cluster, attempt: %d", i+1)
		err = p.Create(name,
			cluster.CreateWithNodeImage(string(version)),
			cluster.CreateWithWaitForReady(20*time.Second),
		)
		if err == nil {
			return func() error {
				return p.Delete(name, "")
			}, nil
		}
	}

	// We failed to create the cluster maxCreateTries times, so return failure error.
	return nil, err
}
