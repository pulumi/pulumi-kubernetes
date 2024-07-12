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

package gen

import (
	"regexp"

	"github.com/pulumi/pulumi/pkg/v3/codegen/cgstrings"
)

var contiguousDigitsRegex = regexp.MustCompile(`\d+`)

// pascalCaseVersions converts idiomatic Kubernetes versions to PascalCase.
// For example, "v1beta1" becomes "V1Beta1".
func pascalCaseVersions(str string) string {
	// Ensure the first character is uppercased.
	str = cgstrings.UppercaseFirst(str)
	return modifyStringAroundDelimeter(str, contiguousDigitsRegex, cgstrings.UppercaseFirst)
}

// ModifyStringAroundDelimeter modifies the string around the delimeter. This is mostly similar to
// cgstrings.ModifyStringAroundDelimeter but allows for a regex based delimeter.
func modifyStringAroundDelimeter(str string, delimRegex *regexp.Regexp, modifyNext func(next string) string) string {
	loc := delimRegex.FindStringIndex(str)
	if loc == nil {
		return str
	}

	nextIdx := loc[1]
	if nextIdx >= len(str) {
		// Nothing left after the delimeter, it's at the end of the string.
		return str
	}

	prev := str[:nextIdx]
	next := str[nextIdx:]
	if next != "" {
		next = modifyNext(next)
	}
	return prev + modifyStringAroundDelimeter(next, delimRegex, modifyNext)
}
