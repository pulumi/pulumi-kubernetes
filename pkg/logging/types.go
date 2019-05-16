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
	"github.com/pulumi/pulumi/pkg/diag"
)

// Message stores a log string and the severity for the log message.
type Message struct {
	S        string
	Severity diag.Severity
}

func (m Message) String() string {
	return m.S
}

func (m Message) Empty() bool {
	return len(m.S) == 0 && len(m.Severity) == 0
}

func StatusMessage(msg string) Message {
	return Message{S: msg, Severity: diag.Info}
}

func WarningMessage(msg string) Message {
	return Message{S: msg, Severity: diag.Warning}
}

func ErrorMessage(msg string) Message {
	return Message{S: msg, Severity: diag.Error}
}

type Messages []Message

func (m Messages) Infos() Messages {
	var messages Messages
	for _, message := range m {
		if message.Severity == diag.Info {
			messages = append(messages, message)
		}
	}

	return messages
}

func (m Messages) Warnings() Messages {
	var messages Messages
	for _, message := range m {
		if message.Severity == diag.Warning {
			messages = append(messages, message)
		}
	}

	return messages
}

func (m Messages) Errors() Messages {
	var messages Messages
	for _, message := range m {
		if message.Severity == diag.Error {
			messages = append(messages, message)
		}
	}

	return messages
}

