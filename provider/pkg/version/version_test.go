package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAgent(t *testing.T) {
	assert.Equal(t, "pulumi-kubernetes/dev (darwin/arm64) client-go/unknown", UserAgent)
}
