package clusters

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

// Cluster is an interface to interact with a Kubernetes cluster in a test.
type Cluster interface {
	Name() string
	Connect() error
	Delete() error
	KubeconfigPath() string
}

// normalizeName returns a normalized name for the cluster that adheres
// to the Kubernetes naming restrictions.
// Ref: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/
func normalizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, ".", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Trauncate names that are too long.
	if len(name) > 63 {
		name = name[:63]
	}

	// Remove any leading numeric characters.
	for i, c := range name {
		if c >= '0' && c <= '9' {
			continue
		}
		name = name[i:]
		break
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
