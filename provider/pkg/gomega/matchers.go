// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gomega

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/onsi/gomega" //nolint:golint // dot-imports
	"github.com/onsi/gomega/gcustom"
	. "github.com/onsi/gomega/gstruct" //nolint:golint // dot-imports
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// ProtobufStruct matches a protobuf struct by decoding it to a map and then applying the given matcher.
func ProtobufStruct(matcher gomegatypes.GomegaMatcher) gomegatypes.GomegaMatcher {
	return WithTransform(func(actual structpb.Struct) (map[string]interface{}, error) { //nolint:govet // copylocks
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

func MatchValue(v any) gomegatypes.GomegaMatcher {
	var matcher gomegatypes.GomegaMatcher
	switch v := v.(type) {
	case gomegatypes.GomegaMatcher:
		matcher = v
	default:
		matcher = Equal(v)
	}
	return WithTransform(func(v resource.PropertyValue) (any, error) {
		return v.V, nil
	}, matcher)
}

func BeComputed() gomegatypes.GomegaMatcher {
	return Equal(resource.MakeComputed(resource.NewStringProperty("")))
}

type Props map[resource.PropertyKey]gomegatypes.GomegaMatcher

// MatchProps succeeds if the actual value is a resource.PropertyMap and all of the expected properties match.
// Options can be used to ignore extra properties or missing properties.
func MatchProps(options Options, props Props) gomegatypes.GomegaMatcher {
	keys := make(Keys, len(props))
	for p, v := range props {
		keys[p] = v
	}
	return &KeysMatcher{
		Keys:          keys,
		IgnoreExtras:  options&IgnoreExtras != 0,
		IgnoreMissing: options&IgnoreMissing != 0,
	}
}

func BeObject(matcher ...gomegatypes.GomegaMatcher) gomegatypes.GomegaMatcher {
	return WithTransform(func(v resource.PropertyValue) (resource.PropertyMap, error) {
		if !v.IsObject() {
			return nil, errors.New("expected property value of type 'object'")
		}
		return v.ObjectValue(), nil
	}, And(matcher...))
}

func MatchObject(options Options, props Props) gomegatypes.GomegaMatcher {
	return BeObject(MatchProps(options, props))
}

func BeSecret(matcher ...gomegatypes.GomegaMatcher) gomegatypes.GomegaMatcher {
	return WithTransform(func(v resource.PropertyValue) (resource.PropertyValue, error) {
		if !v.IsSecret() {
			return resource.PropertyValue{}, errors.New("expected property value of type 'secret'")
		}
		return v.SecretValue().Element, nil
	}, And(matcher...))
}

func MatchSecret(e gomegatypes.GomegaMatcher) gomegatypes.GomegaMatcher {
	return BeSecret(e)
}

func BeArray(matcher ...gomegatypes.GomegaMatcher) gomegatypes.GomegaMatcher {
	return WithTransform(func(v resource.PropertyValue) ([]resource.PropertyValue, error) {
		if !v.IsArray() {
			return nil, errors.New("expected property value of type 'array'")
		}
		return v.ArrayValue(), nil
	}, And(matcher...))
}

func MatchArrayValue(matcher gomegatypes.GomegaMatcher) gomegatypes.GomegaMatcher {
	return BeArray(matcher)
}

// MatchResourceReferenceValue succeeds if the actual value is a resource.PropertyValue of type ResourceReference
// and the URN and ID match the expected values.
func MatchResourceReferenceValue(urn resource.URN, id string) gomegatypes.GomegaMatcher {
	return gcustom.MakeMatcher(func(v resource.PropertyValue) (bool, error) {
		if !v.IsResourceReference() {
			return false, errors.New("expected property value of type 'ResourceReference'")
		}
		rrv := v.ResourceReferenceValue()
		if rrv.URN != urn {
			return false, nil
		}
		if !rrv.ID.IsString() || rrv.ID.StringValue() != id {
			return false, nil
		}
		return true, nil
	})
}
