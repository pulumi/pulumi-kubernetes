// Copyright 2016-2022, Pulumi Corporation.
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

package ssa

import (
	"reflect"
	"testing"
)

func Test_setRequiredFields(t *testing.T) {
	type args struct {
		live  map[string]interface{}
		obj   map[string]interface{}
		field string
	}
	tests := []struct {
		name     string
		args     args
		wantMap  map[string]interface{}
		wantBool bool
	}{
		{
			name: "nested value",
			args: args{
				live: map[string]interface{}{
					"a": map[string]interface{}{
						"b": "c", // should return this value
						"d": "e", // should be ignored
					},
				},
				obj:   map[string]interface{}{},
				field: "a.b",
			},
			wantMap: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "c",
				},
			},
			wantBool: true,
		},
		{
			name: "nested map",
			args: args{
				live: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{ // should return this map
							"c": "d",
						},
						"e": "f", // should be ignored
					},
				},
				obj:   map[string]interface{}{},
				field: "a.b",
			},
			wantMap: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": "d",
					},
				},
			},
			wantBool: true,
		},
		{
			name: "top level field",
			args: args{
				live: map[string]interface{}{
					"a": "b",
				},
				obj:   map[string]interface{}{},
				field: "a",
			},
			wantMap: map[string]interface{}{
				"a": "b",
			},
			wantBool: true,
		},
		{
			name: "sliced string value",
			args: args{
				live: map[string]interface{}{
					"a": map[string]interface{}{
						"b": []interface{}{"c", "d"}, // should return the second element of this slice
						"e": "f",                     // should be ignored
					},
				},
				obj:   map[string]interface{}{},
				field: "a.b[1]",
			},
			wantMap: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []interface{}{nil, "d"},
				},
			},
			wantBool: true,
		},
		{
			name: "sliced int value",
			args: args{
				live: map[string]interface{}{
					"a": map[string]interface{}{
						"b": []interface{}{1, 2}, // should return the second element of this slice
						"e": "f",                 // should be ignored
					},
				},
				obj:   map[string]interface{}{},
				field: "a.b[1]",
			},
			wantMap: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []interface{}{nil, 2},
				},
			},
			wantBool: true,
		},
		{
			name: "sliced map value",
			args: args{
				live: map[string]interface{}{
					"a": map[string]interface{}{
						"b": []interface{}{ // should return the first element of this slice
							map[string]interface{}{
								"c": "d",
							},
						},
						"e": "f", // should be ignored
					},
				},
				obj:   map[string]interface{}{},
				field: "a.b[0]",
			},
			wantMap: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []interface{}{
						map[string]interface{}{
							"c": "d",
						},
					},
				},
			},
			wantBool: true,
		},
		{
			name: "invalid path",
			args: args{
				live:  map[string]interface{}{},
				obj:   map[string]interface{}{},
				field: "a",
			},
			wantMap:  nil,
			wantBool: false,
		},
		{
			name: "invalid index",
			args: args{
				live: map[string]interface{}{
					"a": []interface{}{},
				},
				obj:   map[string]interface{}{},
				field: "a[1]",
			},
			wantMap:  nil,
			wantBool: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMap, gotBool := setRequiredField(tt.args.live, tt.args.obj, tt.args.field)
			if !reflect.DeepEqual(gotMap, tt.wantMap) {
				t.Errorf("setRequiredField() gotMap = %v, want %v", gotMap, tt.wantMap)
			}
			if gotBool != tt.wantBool {
				t.Errorf("setRequiredField() gotBool = %v, want %v", gotBool, tt.wantBool)
			}
		})
	}
}
