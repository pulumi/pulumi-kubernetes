// Copyright 2016-2018, Pulumi Corporation.
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

package client

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/version"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    version.Info
		expected ServerVersion
		error    bool
	}{
		{
			input:    version.Info{Major: "1", Minor: "6"},
			expected: ServerVersion{Major: 1, Minor: 6},
		},
		{
			input:    version.Info{Major: "1", Minor: "70"},
			expected: ServerVersion{Major: 1, Minor: 70},
		},
		{
			input: version.Info{Major: "1", Minor: "6x"},
			error: true,
		},
		{
			input:    version.Info{Major: "1", Minor: "8+"},
			expected: ServerVersion{Major: 1, Minor: 8},
		},
		{
			input:    version.Info{Major: "", Minor: "", GitVersion: "v1.8.0"},
			expected: ServerVersion{Major: 1, Minor: 8},
		},
		{
			input:    version.Info{Major: "1", Minor: "", GitVersion: "v1.8.0"},
			expected: ServerVersion{Major: 1, Minor: 8},
		},
		{
			input:    version.Info{Major: "", Minor: "8", GitVersion: "v1.8.0"},
			expected: ServerVersion{Major: 1, Minor: 8},
		},
		{
			input:    version.Info{Major: "", Minor: "", GitVersion: "v1.8.8-test.0"},
			expected: ServerVersion{Major: 1, Minor: 8},
		},
		{
			input:    version.Info{Major: "1", Minor: "8", GitVersion: "v1.9.0"},
			expected: ServerVersion{Major: 1, Minor: 8},
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
			t.Errorf("Expected %v, got %v", test.expected, v)
		}
	}
}

func TestVersionCompare(t *testing.T) {
	v := ServerVersion{Major: 2, Minor: 3}
	tests := []struct {
		major, minor, result int
	}{
		{major: 1, minor: 0, result: 1},
		{major: 2, minor: 0, result: 1},
		{major: 2, minor: 2, result: 1},
		{major: 2, minor: 3, result: 0},
		{major: 2, minor: 4, result: -1},
		{major: 3, minor: 0, result: -1},
	}
	for _, test := range tests {
		res := v.Compare(test.major, test.minor)
		if res != test.result {
			t.Errorf("%d.%d => Expected %d, got %d", test.major, test.minor, test.result, res)
		}
	}
}

func TestFqName(t *testing.T) {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "tests/v1alpha1",
			"kind":       "Test",
			"metadata": map[string]interface{}{
				"name": "myname",
			},
		},
	}

	if n := FqName(obj.GetNamespace(), obj.GetName()); n != "myname" {
		t.Errorf("Got %q for %v", n, obj)
	}

	obj.SetNamespace("mynamespace")
	if n := FqName(obj.GetNamespace(), obj.GetName()); n != "mynamespace.myname" {
		t.Errorf("Got %q for %v", n, obj)
	}
}
