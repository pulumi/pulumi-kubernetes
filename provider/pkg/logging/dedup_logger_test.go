package logging

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/urn"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockhost struct {
	status strings.Builder
	perm   strings.Builder
}

func (h *mockhost) Log(_ context.Context, sev diag.Severity, urn resource.URN, msg string) error {
	_, _ = h.perm.WriteString(fmt.Sprintf("%s (%s): %s\n", urn, sev, msg))
	return nil
}

func (h *mockhost) LogStatus(_ context.Context, sev diag.Severity, urn resource.URN, msg string) error {
	_, _ = h.status.WriteString(fmt.Sprintf("%s (%s): %s\n", urn, sev, msg))
	return nil
}

func (*mockhost) EngineConn() *grpc.ClientConn { return nil }

func TestDedupeLogger(t *testing.T) {
	h := &mockhost{}
	l := NewLogger(context.Background(), h, urn.New("stack", "proj", "", "type", "name"))

	l.Log(diag.Warning, "first message")
	l.Log(diag.Info, "second message")
	l.Log(diag.Info, "second message")
	l.Log(diag.Warning, "third message")
	l.Log(diag.Warning, "first message")
	l.Log(diag.Info, "second message")

	l.LogStatus(diag.Warning, "first status message")
	l.LogStatus(diag.Info, "second status message")
	l.LogStatus(diag.Info, "second status message")
	l.LogStatus(diag.Warning, "third status message")
	l.LogStatus(diag.Warning, "first status message")
	l.LogStatus(diag.Info, "second status message")

	want := `urn:pulumi:stack::proj::type::name (warning): first message
urn:pulumi:stack::proj::type::name (info): second message
urn:pulumi:stack::proj::type::name (warning): third message
urn:pulumi:stack::proj::type::name (warning): first status message
urn:pulumi:stack::proj::type::name (info): second status message
urn:pulumi:stack::proj::type::name (warning): third status message
`

	assert.Equal(t, want, h.perm.String()+h.status.String())
}
