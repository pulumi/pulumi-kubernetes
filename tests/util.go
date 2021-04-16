package tests

import (
	"sort"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

func SortResourcesByURN(stackInfo integration.RuntimeValidationStackInfo) {
	sort.Slice(stackInfo.Deployment.Resources, func(i, j int) bool {
		return stackInfo.Deployment.Resources[i].URN < stackInfo.Deployment.Resources[j].URN
	})
}
