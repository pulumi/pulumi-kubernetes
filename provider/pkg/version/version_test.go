package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAgent(t *testing.T) {
	assert.Regexp(t, "^pulumi-kubernetes/dev (.*/.*) client-go/unknown$", UserAgent)
}
