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

import "testing"

func TestPascalCaseVersions(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{
			"v1alpha1",
			"V1Alpha1",
		},
		{
			"v1beta1",
			"V1Beta1",
		},
		{
			"v1",
			"V1",
		},
		{
			"v2",
			"V2",
		},
		{
			"v2beta1",
			"V2Beta1",
		},
		{
			"v23gamma123",
			"V23Gamma123",
		},
		{
			"123foo321bar021baz",
			"123Foo321Bar021Baz",
		},
		{
			"123foo321bar021baz8",
			"123Foo321Bar021Baz8",
		},
		{
			"apiregistration",
			"Apiregistration",
		},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got := pascalCaseVersions(tt.version); got != tt.want {
				t.Errorf("pascalCaseVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}
