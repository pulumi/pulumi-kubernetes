//nolint:copylocks
package gomega

import (
	"fmt"

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
	return &aliasNameMatcher{Name: name}
}

// AliasByType matches an Alias by its type.
func AliasByType(typ tokens.Type) gomegatypes.GomegaMatcher {
	return &aliasTypeMatcher{Type: typ}
}

type aliasNameMatcher struct {
	Name tokens.QName
}

var _ gomegatypes.GomegaMatcher = &aliasNameMatcher{}

func (matcher *aliasNameMatcher) Match(actual interface{}) (success bool, err error) {
	if alias, ok := actual.(*pulumirpc.Alias); ok {
		return alias.GetSpec() != nil && alias.GetSpec().Name == string(matcher.Name), nil
	}
	return false, fmt.Errorf("aliasNameMatcher matcher expects a pulumirpc.Alias")
}

func (matcher *aliasNameMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected:\n\t%#v\nto have name %q", actual, matcher.Name)
}

func (matcher *aliasNameMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected:\n\t%#v\nnot to have name %q", actual, matcher.Name)
}

type aliasTypeMatcher struct {
	Type tokens.Type
}

var _ gomegatypes.GomegaMatcher = &aliasTypeMatcher{}

func (matcher *aliasTypeMatcher) Match(actual interface{}) (success bool, err error) {
	if alias, ok := actual.(*pulumirpc.Alias); ok {
		return alias.GetSpec() != nil && alias.GetSpec().Type == string(matcher.Type), nil
	}
	return false, fmt.Errorf("aliasTypeMatcher matcher expects a pulumirpc.Alias")
}

func (matcher *aliasTypeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected:\n\t%#v\nto have type %q", actual, matcher.Type)
}

func (matcher *aliasTypeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected:\n\t%#v\nnot to have type %q", actual, matcher.Type)
}
