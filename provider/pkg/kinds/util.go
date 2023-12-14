package kinds

import (
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func suffixSearch(urnS string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(urnS, suffix) {
			return true
		}
	}
	return false
}

// IsPatchURN returns true if the URN is for a Patch resource.
func IsPatchURN(urn resource.URN) bool {
	urnS := urn.QualifiedType().String()

	// Do a simple O(1) lookup in a set of known patch types.
	if PatchQualifiedTypes.Has(urnS) {
		return true
	}

	// Do a suffix search in case the resource is a component resource containing
	// a patch resource. Eg. a component resource patch resource could have the following URN:
	// my:component:Resource$kubernetes:apps/v1:DaemonSetPatch.
	return suffixSearch(urnS, PatchQualifiedTypes.SortedValues())
}

// IsListURN returns true if the URN is for a List resource.
func IsListURN(urn resource.URN) bool {
	urnS := urn.QualifiedType().String()

	// Do a simple O(1) lookup in a set of known list types.
	if ListQualifiedTypes.Has(urnS) {
		return true
	}

	// Do a suffix search in case the resource is a component resource containing
	// a list resource.
	return suffixSearch(urnS, ListQualifiedTypes.SortedValues())
}
