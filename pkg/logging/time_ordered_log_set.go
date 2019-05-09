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
	"fmt"

	"github.com/pulumi/pulumi/pkg/diag"
)

// Message stores a log string and the severity for the log message.
type Message struct {
	s        string
	severity diag.Severity
}

func (m Message) String() string {
	return fmt.Sprintf("[%s] %s", m.severity, m.s)
}

// TimeOrderedLogSet stores a temporally-ordered set of log messages.
type TimeOrderedLogSet struct {
	exists   map[Message]bool
	Messages []Message
}

// Add appends a message to the time-ordered set.
func (o *TimeOrderedLogSet) Add(msg Message) {
	// Ensure memory has been allocated.
	if o.exists == nil {
		o.exists = make(map[Message]bool)
	}
	if o.Messages == nil {
		o.Messages = []Message{}
	}

	if !o.exists[msg] {
		o.Messages = append(o.Messages, msg)
		o.exists[msg] = true
	}
}
