package kinds

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// IsPatchResource returns true if it is a Patch resource and also ensures
// that we don't return true for CustomResources that end with "Patch".
func IsPatchResource(urn resource.URN, kind string) bool {
	resourceName := urn.Type().Name().String()

	return kind+"Patch" == resourceName
}

// IsListURN returns true if the URN is for a List resource.
func IsListURN(urn resource.URN) bool {
	urnS := urn.Type().String()

	return ListQualifiedTypes.Has(urnS)
}
