package provider

import (
	"testing"
)

func TestNormalizeChartRef(t *testing.T) {
	check := func(repoName string, repoUrl string, originalChartRef string, expect string) {
		actual := normalizeChartRef(repoName, repoUrl, originalChartRef)
		if actual != expect {
			t.Errorf("Expected normalizeChartRef(%s, %s, %s) to be %s but got %s",
				repoName, repoUrl, originalChartRef, expect, actual)
		}
	}

	check("bitnami", "https://charts.bitnami.com/bitnami", "apache", "apache")
	check("bitnami", "", "apache", "bitnami/apache")
	check("bitnami", "", "bitnami/apache", "bitnami/apache")
}
