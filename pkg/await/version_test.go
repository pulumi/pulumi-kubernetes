// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package await

import (
	"testing"

	"k8s.io/apimachinery/pkg/version"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    version.Info
		expected serverVersion
		error    bool
	}{
		{
			input:    version.Info{Major: "1", Minor: "6"},
			expected: serverVersion{Major: 1, Minor: 6},
		},
		{
			input:    version.Info{Major: "1", Minor: "70"},
			expected: serverVersion{Major: 1, Minor: 70},
		},
		{
			input: version.Info{Major: "1", Minor: "6x"},
			error: true,
		},
		{
			input:    version.Info{Major: "1", Minor: "8+"},
			expected: serverVersion{Major: 1, Minor: 8},
		},
		{
			input:    version.Info{Major: "", Minor: "", GitVersion: "v1.8.0"},
			expected: serverVersion{Major: 1, Minor: 8},
		},
		{
			input:    version.Info{Major: "1", Minor: "", GitVersion: "v1.8.0"},
			expected: serverVersion{Major: 1, Minor: 8},
		},
		{
			input:    version.Info{Major: "", Minor: "8", GitVersion: "v1.8.0"},
			expected: serverVersion{Major: 1, Minor: 8},
		},
		{
			input:    version.Info{Major: "", Minor: "", GitVersion: "v1.8.8-test.0"},
			expected: serverVersion{Major: 1, Minor: 8, Patch: 8},
		},
		{
			input:    version.Info{Major: "1", Minor: "8", GitVersion: "v1.9.0"},
			expected: serverVersion{Major: 1, Minor: 9},
		},
		{
			input:    version.Info{Major: "1", Minor: "9", GitVersion: "v1.9.1"},
			expected: serverVersion{Major: 1, Minor: 9, Patch: 1},
		},
		{
			input:    version.Info{Major: "1", Minor: "13", GitVersion: "v1.13.0"},
			expected: serverVersion{Major: 1, Minor: 13, Patch: 0},
		},
		{
			input: version.Info{Major: "", Minor: "", GitVersion: "v1.a"},
			error: true,
		},
	}

	for _, test := range tests {
		v, err := parseVersion(&test.input)
		if test.error {
			if err == nil {
				t.Errorf("test %s should have failed and did not", test.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("test %v failed: %v", test.input, err)
			continue
		}
		if v != test.expected {
			t.Errorf("Expected %#v, got %#v", test.expected, v)
		}
	}
}

func TestVersionCompare(t *testing.T) {
	v := serverVersion{Major: 2, Minor: 3, Patch: 0}
	tests := []struct {
		major, minor, patch, result int
	}{
		{major: 1, minor: 0, result: 1},
		{major: 2, minor: 0, result: 1},
		{major: 2, minor: 2, result: 1},
		{major: 2, minor: 2, result: 1},
		{major: 2, minor: 2, patch: 2, result: 1},
		{major: 2, minor: 3, result: 0},
		{major: 2, minor: 3, patch: 0, result: 0},
		{major: 2, minor: 4, result: -1},
		{major: 3, minor: 0, result: -1},
		{major: 2, minor: 3, patch: 1, result: -1},
	}
	for _, test := range tests {
		res := v.Compare(test.major, test.minor, test.patch)
		if res != test.result {
			t.Errorf("%d.%d.%d => Expected %d, got %d", test.major, test.minor, test.patch, test.result, res)
		}
	}
}
