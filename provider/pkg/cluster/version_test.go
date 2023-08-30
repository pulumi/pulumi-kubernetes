// Copyright 2016-2019, Pulumi Corporation.
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

package cluster

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/version"
)

func TestServerVersion_Compare(t *testing.T) {
	v := ServerVersion{Major: 2, Minor: 3}

	tests := []struct {
		name    string
		input   ServerVersion
		version ServerVersion
		want    int
	}{
		{"Older major", ServerVersion{1, 0}, v, 1},
		{"Older minor", ServerVersion{2, 0}, v, 1},
		{"Equal", ServerVersion{2, 3}, v, 0},
		{"Newer minor", ServerVersion{2, 4}, v, -1},
		{"Newer major", ServerVersion{3, 0}, v, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v.Compare(tt.input); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseGitVersion(t *testing.T) {
	type args struct {
		versionString string
	}
	tests := []struct {
		name    string
		args    args
		want    gitVersion
		wantErr bool
	}{
		{
			name: "Valid",
			args: args{"v1.13.1"},
			want: gitVersion{1, 13, 1},
		},
		{
			name: "Valid + suffix",
			args: args{"v1.8.8-test.0"},
			want: gitVersion{1, 8, 8},
		},
		{
			name:    "Missing v",
			args:    args{"1.13.0"},
			wantErr: true,
		},
		{
			name:    "Missing patch",
			args:    args{"1.13"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGitVersion(tt.args.versionString)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGitVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseGitVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			input:    version.Info{Major: "1", Minor: "9", GitVersion: "v1.9.1"},
			expected: ServerVersion{Major: 1, Minor: 9},
		},
		{
			input:    version.Info{Major: "1", Minor: "13", GitVersion: "v1.13.0"},
			expected: ServerVersion{Major: 1, Minor: 13},
		},
		{
			input: version.Info{Major: "", Minor: "", GitVersion: "v1.a"},
			error: true,
		},
	}

	for _, test := range tests {
		test := test
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
