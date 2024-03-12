package tests

import (
	"os/exec"
	"sort"
	"strings"
	"testing"

	"golang.org/x/exp/slices"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// SkipIfShort skips the test if the -short flag is passed to `go test`.
func SkipIfShort(t *testing.T, msg string) {
	if testing.Short() {
		t.Skip("skipping", t.Name(), "in short mode:", msg)
	}
}

func SortResourcesByURN(stackInfo integration.RuntimeValidationStackInfo) {
	sort.Slice(stackInfo.Deployment.Resources, func(i, j int) bool {
		return stackInfo.Deployment.Resources[i].URN < stackInfo.Deployment.Resources[j].URN
	})
}

func SearchResourcesByName(stackInfo integration.RuntimeValidationStackInfo, parent resource.URN, typ tokens.Type, name string) *apitype.ResourceV3 {
	i := slices.IndexFunc(stackInfo.Deployment.Resources, func(r apitype.ResourceV3) bool {
		return r.Type == typ && r.URN.Name() == name && (parent == "" || r.Parent == parent)
	})
	if i == -1 {
		return nil
	}
	return &stackInfo.Deployment.Resources[i]
}

// Kubectl is a helper function to shell out and run kubectl commands.
func Kubectl(args ...string) ([]byte, error) {
	var fmtArgs []string
	for _, arg := range args {
		fmtArgs = append(fmtArgs, strings.Fields(arg)...)
	}

	return exec.Command("kubectl", fmtArgs...).CombinedOutput()
}
