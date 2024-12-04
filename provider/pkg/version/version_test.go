package version

import (
	"fmt"
	"os"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAgent(t *testing.T) {
	assert.Regexp(t, "^pulumi-kubernetes/dev (.*/.*) client-go/unknown$", UserAgent)
}

func TestKubernetesMinorVersion(t *testing.T) {
	v := os.Getenv("KUBE_VERSION")
	if v == "" {
		t.Skip("KUBE_VERSION isn't set")
	}

	version, err := semver.ParseTolerant(v)
	require.NoError(t, err)

	mod, err := os.ReadFile("../../go.mod")
	require.NoError(t, err)

	want := fmt.Sprintf("k8s.io/api v0.%d", version.Minor)

	assert.Contains(t, string(mod), want, "KUBE_VERSION=%s doesn't match go.mod's minor version", v)

}
