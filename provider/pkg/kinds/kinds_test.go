package kinds

import (
	"testing"
)

func TestKind_Namespaced(t *testing.T) {
	tests := []struct {
		name           string
		k              Kind
		wantKnown      bool
		wantNamespaced bool
	}{
		{
			"Standard kind",
			Deployment,
			true,
			true,
		},
		{
			"List kind",
			DeploymentList,
			true,
			true,
		},
		{
			"Unknown kind",
			Kind("TokenReview"),
			false,
			false,
		},
		{
			"Unknown list kind",
			Kind("TokenReviewList"),
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKnown, gotNamespaced := tt.k.Namespaced()
			if gotKnown != tt.wantKnown {
				t.Errorf("Namespaced() gotKnown = %v, want %v", gotKnown, tt.wantKnown)
			}
			if gotNamespaced != tt.wantNamespaced {
				t.Errorf("Namespaced() gotNamespaced = %v, want %v", gotNamespaced, tt.wantNamespaced)
			}
		})
	}
}
