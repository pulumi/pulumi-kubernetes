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

package logging

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/stretchr/testify/assert"
)

func TestOrderedStringSet_Add(t *testing.T) {
	status1 := Message{"foo", diag.Info}
	status2 := Message{"bar", diag.Info}
	warn1 := Message{"boom", diag.Warning}

	type fields struct {
		exists   map[Message]bool
		Messages []Message
	}
	type args struct {
		msg Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect []Message
	}{
		{
			"add a message to uninitialized struct",
			fields{},
			args{status1},
			[]Message{status1},
		},
		{
			"add a new message to empty list",
			fields{map[Message]bool{}, []Message{}},
			args{status1},
			[]Message{status1},
		},
		{
			"add a new info message to existing list",
			fields{map[Message]bool{status1: true}, []Message{status1}},
			args{status2},
			[]Message{status1, status2},
		},
		{
			"add a new warning message to existing list",
			fields{map[Message]bool{status1: true}, []Message{status1}},
			args{warn1},
			[]Message{status1, warn1},
		},
		{
			"add a duplicate string",
			fields{map[Message]bool{status1: true}, []Message{status1}},
			args{status1},
			[]Message{status1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &TimeOrderedLogSet{
				exists:   tt.fields.exists,
				Messages: tt.fields.Messages,
			}
			o.Add(tt.args.msg)
			assert.ObjectsAreEqual(o.Messages, tt.expect)
		})
	}
}
