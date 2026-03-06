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
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

// Message stores a log string and the severity for the log message.
type Message struct {
	S        string
	Severity diag.Severity
}

func (m Message) String() string {
	return m.S
}

// Empty returns true if the Message is uninitialized, false otherwise.
func (m Message) Empty() bool {
	return len(m.S) == 0 && len(m.Severity) == 0
}

// StatusMessage creates a Message with Severity set to Info.
func StatusMessage(msg string) Message {
	return Message{S: msg, Severity: diag.Info}
}

// WarningMessage creates a Message with Severity set to Warning.
func WarningMessage(msg string) Message {
	return Message{S: msg, Severity: diag.Warning}
}

// ErrorMessage creates a Message with Severity set to Error.
func ErrorMessage(msg string) Message {
	return Message{S: msg, Severity: diag.Error}
}

// Messages is a slice of Message types.
type Messages []Message

// Infos returns Messages with Info severity.
func (m Messages) Infos() Messages {
	var messages Messages
	for _, message := range m {
		if message.Severity == diag.Info {
			messages = append(messages, message)
		}
	}

	return messages
}

// Warnings returns Messages with Warning severity.
func (m Messages) Warnings() Messages {
	var messages Messages
	for _, message := range m {
		if message.Severity == diag.Warning {
			messages = append(messages, message)
		}
	}

	return messages
}

// Errors returns Messages with Error severity.
func (m Messages) Errors() Messages {
	var messages Messages
	for _, message := range m {
		if message.Severity == diag.Error {
			messages = append(messages, message)
		}
	}

	return messages
}

// MessagesWithSeverity returns Messages matching any of the provided Severity levels.
func (m Messages) MessagesWithSeverity(sev ...diag.Severity) Messages {
	var messages Messages
	for _, message := range m {
		for _, s := range sev {
			if message.Severity == s {
				messages = append(messages, message)
			}
		}
	}

	return messages
}
