// Copyright 2021, Pulumi Corporation.  All rights reserved.

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

	check(
		"bitnami",
		"https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami",
		"apache",
		"apache",
	)
	check("bitnami", "", "apache", "bitnami/apache")
	check("bitnami", "", "bitnami/apache", "bitnami/apache")
}
