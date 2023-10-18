package tests

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/stretchr/testify/assert"
)

// SkipIfShort skips the test if the -short flag is passed to `go test`.
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
}

func SortResourcesByURN(stackInfo integration.RuntimeValidationStackInfo) {
	sort.Slice(stackInfo.Deployment.Resources, func(i, j int) bool {
		return stackInfo.Deployment.Resources[i].URN < stackInfo.Deployment.Resources[j].URN
	})
}

// Kubectl is a helper function to shell out and run kubectl commands.
func Kubectl(args ...string) ([]byte, error) {
	var fmtArgs []string
	for _, arg := range args {
		fmtArgs = append(fmtArgs, strings.Fields(arg)...)
	}

	return exec.Command("kubectl", fmtArgs...).CombinedOutput()
}

// Matcher provides an interface for matching values.
type Matcher[T any] interface {
	Match(actual T) bool
}

// AssertEvents asserts that the stack info contains a set of events.
func AssertEvents(t assert.TestingT, stackInfo integration.RuntimeValidationStackInfo, expected ...Matcher[apitype.EngineEvent]) (success bool) {
	success = true
Expected:
	for _, m := range expected {
		for _, evt := range stackInfo.Events {
			if m.Match(evt) {
				continue Expected
			}
		}
		assert.Fail(t, fmt.Sprintf("Expected an engine event: %+v", m))
		success = false
	}
	if tt, ok := t.(*testing.T); ok && !success {
		json, err := json.MarshalIndent(stackInfo.Events, "", "  ")
		contract.AssertNoErrorf(err, "unexpected JSON error: %v", err)
		tt.Logf("Actual engine events:\n%s\n", json)
	}
	return success
}

// ResOutputsEvent matches resource output events.
type ResOutputsEvent struct {
	// Op is the operation being performed.
	Op apitype.OpType
	// Type is the resource type of the event.
	Type string
	// Name is the resource name.
	Name tokens.QName
	// Keys causing a replacement (only applicable for "create" and "replace" Ops).
	Keys []string
	// Keys that changed with this step.
	Diffs []string
}

func (e ResOutputsEvent) Match(actual apitype.EngineEvent) bool {
	if actual.ResOutputsEvent == nil {
		return false
	}
	urn, err := resource.ParseURN(actual.ResOutputsEvent.Metadata.URN)
	if err != nil {
		return false
	}
	a := ResOutputsEvent{
		Op:    actual.ResOutputsEvent.Metadata.Op,
		Type:  actual.ResOutputsEvent.Metadata.Type,
		Name:  urn.Name(),
		Keys:  actual.ResOutputsEvent.Metadata.Keys,
		Diffs: actual.ResOutputsEvent.Metadata.Diffs,
	}
	sort.Strings(a.Keys)
	sort.Strings(a.Diffs)
	return reflect.DeepEqual(e, a)
}
