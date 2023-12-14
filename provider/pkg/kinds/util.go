package kinds

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// IsPatchURN returns true if the URN is for a Patch resource.
func IsPatchURN(urn resource.URN) bool {
	urnS := urn.Type().String()

	return PatchQualifiedTypes.Has(urnS)
}

// IsListURN returns true if the URN is for a List resource.
func IsListURN(urn resource.URN) bool {
	urnS := urn.Type().String()

	return ListQualifiedTypes.Has(urnS)
}
