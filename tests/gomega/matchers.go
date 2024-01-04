//nolint:copylocks
package gomega

import (
	"fmt"
	"strings"

	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
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

// Alias matches an Alias by name, type, and/or parent.
// The following option types are supported:
// - resource.URN - matches the parent URN
// - tokens.Type - matches the type
// - tokens.QName - matches the name
// - string - matches the name
func Alias(opts ...any) gomegatypes.GomegaMatcher {
	m := AliasMatcher{}
	for len(opts) > 0 {
		switch v := opts[0].(type) {
		case resource.URN:
			m.ParentURN = &v
		case tokens.Type:
			m.Type = &v
		case tokens.QName:
			m.Name = &v
		case string:
			q := tokens.QName(v)
			m.Name = &q
		default:
		}
		opts = opts[1:]
	}
	return &m
}

type AliasMatcher struct {
	Name      *tokens.QName
	Type      *tokens.Type
	ParentURN *resource.URN
	NoParent  *bool
}

var _ gomegatypes.GomegaMatcher = &AliasMatcher{}

func (matcher *AliasMatcher) Match(actual interface{}) (success bool, err error) {
	if alias, ok := actual.(*pulumirpc.Alias); ok {
		if matcher.Name != nil && alias.GetSpec().GetName() != string(*matcher.Name) {
			return false, nil
		}
		if matcher.Type != nil && alias.GetSpec().GetType() != string(*matcher.Type) {
			return false, nil
		}
		if matcher.ParentURN != nil && alias.GetSpec().GetParentUrn() != string(*matcher.ParentURN) {
			return false, nil
		}
		if matcher.NoParent != nil && alias.GetSpec().GetNoParent() != *matcher.NoParent {
			return false, nil
		}
		return true, nil
	}
	return false, fmt.Errorf("aliasNameMatcher matcher expects a *pulumirpc.Alias")
}

func (matcher *AliasMatcher) FailureMessage(actual interface{}) (message string) {
	var msg strings.Builder
	fmt.Fprintf(&msg, "Expected:\n\t%#v\nto have ", actual)
	if matcher.Name != nil {
		fmt.Fprintf(&msg, "name=%q", *matcher.Name)
	}
	if matcher.Type != nil {
		fmt.Fprintf(&msg, "type=%q", *matcher.Type)
	}
	if matcher.ParentURN != nil {
		fmt.Fprintf(&msg, "parentURN=%q", *matcher.ParentURN)
	}
	if matcher.NoParent != nil {
		fmt.Fprintf(&msg, "noParent=%t", *matcher.NoParent)
	}
	return msg.String()
}

func (matcher *AliasMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	var msg strings.Builder
	fmt.Fprintf(&msg, "Expected:\n\t%#v\nto not have ", actual)
	if matcher.Name != nil {
		fmt.Fprintf(&msg, "name=%q", *matcher.Name)
	}
	if matcher.Type != nil {
		fmt.Fprintf(&msg, "type=%q", *matcher.Type)
	}
	if matcher.ParentURN != nil {
		fmt.Fprintf(&msg, "parentURN=%q", *matcher.ParentURN)
	}
	if matcher.NoParent != nil {
		fmt.Fprintf(&msg, "noParent=%t", *matcher.NoParent)
	}
	return msg.String()
}
