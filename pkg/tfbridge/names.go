// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"unicode"

	"github.com/pulumi/pulumi-fabric/pkg/resource"
	"github.com/pulumi/pulumi-fabric/pkg/util/contract"
)

// RandomHexSuffixLength is the length of the suffix added AutoName properties by default.
const RandomHexSuffixLength = 8

// LumiToTerraformName performs a standard transformation on the given name string, from Lumi's PascalCasing or
// camelCasing, to Terraform's underscore_casing.
func LumiToTerraformName(name string) string {
	var result string
	for i, c := range name {
		if c >= 'A' && c <= 'Z' {
			// if upper case, add an underscore (if it's not #1), and then the lower case version.
			if i != 0 {
				result += "_"
			}
			result += string(unicode.ToLower(c))
		} else {
			result += string(c)
		}
	}
	return result
}

// TerraformToLumiName performs a standard transformation on the given name string, from Terraform's underscore_casing
// to Lumi's PascalCasing (if upper is true) or camelCasing (if upper is false).
func TerraformToLumiName(name string, upper bool) string {
	var result string
	var nextCap bool
	var prev rune
	for i, c := range name {
		if c == '_' {
			// skip underscores and make sure the next one is capitalized.
			contract.Assertf(!nextCap, "Unexpected duplicate underscore: %v", name)
			contract.Assertf(i != 0, "Unexpected underscore as 1st character: %v", name)
			nextCap = true
		} else {
			if ((i == 0 && upper) || nextCap) && (c >= 'a' && c <= 'z') {
				// if we're at the start and upper was requested, or the next is meant to be a cap, capitalize it.
				result += string(unicode.ToUpper(c))
			} else {
				result += string(c)
			}
			nextCap = false
		}
		prev = c
	}
	if prev == '_' {
		// we had a next cap, but it wasn't realized.  propagate the _ after all.
		result += "_"
	}
	return result
}

// IsBuiltinLumiProperty returns true if the property name s is a special Lumi builtin property.
func IsBuiltinLumiProperty(s string) bool {
	return s == string(resource.IDProperty) ||
		s == string(resource.URNProperty) ||
		s == string(resource.URNNameProperty)
}

// AutoName creates custom schema for a Terraform name property which is automatically populated from the
// resource's URN name, with a random suffix and maximum length of maxlen.  This makes it easy to propagate the Lumi
// resource's URN name part as the Terraform name as a convenient default, while still permitting it to be overridden.
func AutoName(name string, maxlen int) *SchemaInfo {
	return AutoNameTransform(name, maxlen, nil)
}

// AutoNameTransform creates custom schema for a Terraform name property which is automatically populated from the
// resource's URN name, with a random suffix, maximum length maxlen, and optional transformation function. This makes it
// easy to propagate the Lumi resource's URN name part as the Terraform name as a convenient default, while still
// permitting it to be overridden.
func AutoNameTransform(name string, maxlen int, transform func(string) string) *SchemaInfo {
	contract.Assert(name != string(resource.URNNamePropertyKey))
	return &SchemaInfo{
		Name: name,
		Default: &DefaultInfo{
			From: string(resource.URNNamePropertyKey),
			FromTransform: func(v interface{}) interface{} {
				vs := v.(string)
				if transform != nil {
					vs = transform(vs)
				}
				return resource.NewUniqueHex(vs+"-", maxlen, RandomHexSuffixLength)
			},
		},
	}
}
