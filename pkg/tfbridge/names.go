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
	return &SchemaInfo{
		Name: name,
		Default: &DefaultInfo{
			From: FromName(true, maxlen, transform),
		},
	}
}

// FromName automatically propagates a resource's URN onto the resulting default info.
func FromName(rand bool, randmaxlen int, transform func(string) string) func(res *LumiResource) interface{} {
	return func(res *LumiResource) interface{} {
		// Take the URN name part, transform it if required, and then append some unique characters.
		vs := string(res.URN.Name())
		if transform != nil {
			vs = transform(vs)
		}
		if rand {
			return resource.NewUniqueHex(vs+"-", randmaxlen, RandomHexSuffixLength)
		}
		return vs
	}
}
