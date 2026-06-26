package gen

import (
	"strings"
	"testing"
)

func TestHelmV3ChartOverlayExamplesDoNotReferenceUndefinedIdentifiers(t *testing.T) {
	t.Helper()

	if strings.Contains(helmV3ChartMD, "namespaceName") {
		t.Fatal("helm v3 chart overlay contains an undefined namespaceName identifier")
	}
	if strings.Contains(helmV3ChartMD, "DependsOnInputs(chart.Ready)") {
		t.Fatal("helm v3 chart overlay contains an undefined chart identifier")
	}
}

func TestHelmV3ReleaseOverlayExamplesUseReleaseTerminology(t *testing.T) {
	t.Helper()

	if strings.Contains(helmV3ReleaseMD, "### Depend on a Chart resource") {
		t.Fatal("helm v3 release overlay still labels the example as a chart dependency")
	}
	if strings.Contains(helmV3ReleaseMD, "metadata: {namespace: namespaceName}") {
		t.Fatal("helm v3 release overlay contains an undefined namespaceName identifier")
	}
}
