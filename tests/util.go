package tests

import (
	"os/exec"
	"sort"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
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
