// Copyright 2016-2023, Pulumi Corporation.
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

package crd

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var alphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// SplitGroupVersion returns the <group> and <version> field of a string in the
// format <group>/<version>
func SplitGroupVersion(groupVersion string) (string, string, error) {
	parts := strings.Split(groupVersion, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected a version string with the format <group>/<version>, but got %q", groupVersion)
	}
	return parts[0], parts[1], nil
}

// groupPrefix returns the first word in the dot-separated group string, with
// all non-alphanumeric characters removed.
func GroupPrefix(group string) (string, error) {
	if group == "" {
		return "", fmt.Errorf("group cannot be empty")
	}
	return removeNonAlphanumeric(strings.Split(group, ".")[0]), nil
}

// Capitalizes and returns the given version. For example,
// VersionToUpper("v2beta1") returns "V2Beta1".
func VersionToUpper(version string) string {
	var sb strings.Builder
	for i, r := range version {
		if unicode.IsLetter(r) && (i == 0 || !unicode.IsLetter(rune(version[i-1]))) {
			sb.WriteRune(unicode.ToUpper(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// removeNonAlphanumeric removes all non-alphanumeric characters
func removeNonAlphanumeric(input string) string {
	return alphanumericRegex.ReplaceAllString(input, "")
}
