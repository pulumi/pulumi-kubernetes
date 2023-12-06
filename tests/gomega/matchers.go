//nolint:copylocks
package gomega

import (
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	structpb "google.golang.org/protobuf/types/known/structpb"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// ProtobufStruct matches a protobuf struct by decoding it to a map and then applying the given matcher.
func ProtobufStruct(matcher gomegatypes.GomegaMatcher) gomegatypes.GomegaMatcher {
	return WithTransform(func(actual structpb.Struct) (map[string]interface{}, error) {
		m := actual.AsMap()
		return m, nil
	}, matcher)
}

// Alias matches an Alias by its name.
func Alias(name tokens.QName) gomegatypes.GomegaMatcher {
	return WithTransform(func(actual *pulumirpc.Alias) string {
		if actual.GetSpec() == nil {
			return ""
		}
		return actual.GetSpec().Name
	}, BeEquivalentTo(name))
}

// AliasByType matches an Alias by its type.
func AliasByType(typ tokens.Type) gomegatypes.GomegaMatcher {
	return WithTransform(func(actual *pulumirpc.Alias) string {
		if actual.GetSpec() == nil {
			return ""
		}
		return actual.GetSpec().Type
	}, BeEquivalentTo(typ))
}
