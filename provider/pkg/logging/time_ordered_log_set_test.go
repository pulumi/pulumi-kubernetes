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

package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

func TestOrderedStringSet_Add(t *testing.T) {
	status1 := logging.Message{S: "foo", Severity: diag.Info}
	status2 := logging.Message{S: "bar", Severity: diag.Info}
	warn1 := logging.Message{S: "boom", Severity: diag.Warning}

	type fields struct {
		exists   map[logging.Message]bool
		Messages []logging.Message
	}
	type args struct {
		msg logging.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect []logging.Message
	}{
		{
			"add a message to uninitialized struct",
			fields{},
			args{status1},
			[]logging.Message{status1},
		},
		{
			"add a new message to empty list",
			fields{map[logging.Message]bool{}, []logging.Message{}},
			args{status1},
			[]logging.Message{status1},
		},
		{
			"add a new info message to existing list",
			fields{map[logging.Message]bool{status1: true}, []logging.Message{status1}},
			args{status2},
			[]logging.Message{status1, status2},
		},
		{
			"add a new warning message to existing list",
			fields{map[logging.Message]bool{status1: true}, []logging.Message{status1}},
			args{warn1},
			[]logging.Message{status1, warn1},
		},
		{
			"add a duplicate string",
			fields{map[logging.Message]bool{status1: true}, []logging.Message{status1}},
			args{status1},
			[]logging.Message{status1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			o := &TimeOrderedLogSet{
				exists:   tt.fields.exists,
				Messages: tt.fields.Messages,
			}
			o.Add(tt.args.msg)
			assert.ObjectsAreEqual(o.Messages, tt.expect)
		})
	}
}
