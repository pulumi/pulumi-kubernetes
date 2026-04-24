// Copyright 2016-2021, Pulumi Corporation.
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

package checker

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
)

// Result specifies the result of a Condition applied to an input object.
type Result struct {
	Ok          bool            // True if the Condition is true, false otherwise.
	Description string          // A human-readable description of the associated Condition.
	Message     logging.Message // The message to be logged after evaluating the Condition.
}

func (r Result) String() string {
	var s string
	if r.Ok {
		s = fmt.Sprintf(`["done"] %s`, r.Description)
	} else {
		s = fmt.Sprintf(`["pending"] %s`, r.Description)
	}

	if !r.Message.Empty() {
		s = fmt.Sprintf("%s -- %s", s, r.Message)
	}

	return s
}

// Results is a slice of Result objects.
type Results []Result

func (rr Results) String() string {
	s := strings.Builder{}
	for _, r := range rr {
		fmt.Fprintf(&s, "%s\n", r)
	}

	return s.String()
}

// Messages iterates the Results and returns a slice of the underlying Message objects. Note that these messages are
// not cached, so each invocation of this method will allocate memory for the slice.
func (rr Results) Messages() logging.Messages {
	var messages logging.Messages
	for _, r := range rr {
		if !r.Message.Empty() {
			messages = append(messages, r.Message)
		}
	}
	return messages
}

// Condition is a function that checks a state and returns a Result.
type Condition func(state interface{}) Result

// StateChecker holds the data required to generically implement await logic.
type StateChecker struct {
	conditions []Condition // Conditions that must be true for the state to be Ready.
}

type StateCheckerArgs struct {
	Conditions []Condition // Conditions that must be true for the state to be Ready.
}

func NewStateChecker(args *StateCheckerArgs) *StateChecker {
	return &StateChecker{
		conditions: args.Conditions,
	}
}

func (s *StateChecker) Ready(state interface{}) bool {
	ok, _ := s.readyDetails(state)
	return ok
}

func (s *StateChecker) ReadyStatus(state interface{}) (bool, Result) {
	ok, results := s.readyDetails(state)
	return ok, results[len(results)-1]
}

func (s *StateChecker) ReadyDetails(state interface{}) (bool, Results) {
	return s.readyDetails(state)
}

func (s *StateChecker) readyDetails(state interface{}) (bool, Results) {
	var results Results

	for _, condition := range s.conditions {
		result := condition(state)
		results = append(results, result)
		if !result.Ok {
			return false, results
		}
	}

	return true, results
}
