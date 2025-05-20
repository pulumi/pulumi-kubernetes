package version

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAgent(t *testing.T) {
	assert.Regexp(t, "^pulumi-kubernetes/dev (.*/.*) client-go/unknown$", UserAgent)
}

func TestKubernetesMinorVersion(t *testing.T) {

	re, err := regexp.Compile(`KUBE_VERSION\s+\?=\s+v(?P<version>.*)`)
	require.NoError(t, err)

	mf, err := os.ReadFile("../../../Makefile")
	require.NoError(t, err)

	matches := re.FindStringSubmatch(string(mf))
	require.Len(t, matches, 2)

	version, err := semver.ParseTolerant(matches[1])
	require.NoError(t, err)

	mod, err := os.ReadFile("../../go.mod")
	require.NoError(t, err)

	want := fmt.Sprintf("k8s.io/api v0.%d", version.Minor)

	assert.Contains(t, string(mod), want, "KUBE_VERSION=v%s doesn't match go.mod's minor version", version)
}
