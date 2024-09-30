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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaseMapping_Add(t *testing.T) {
	cm := CaseMapping{mapping: make(map[string]string)}

	err := cm.Add("test", "Test")
	assert.NoError(t, err)

	err = cm.Add("test", "Test")
	assert.Error(t, err)
	assert.Equal(t, "case mapping for \"test\" already exists", err.Error())
}

func TestCaseMapping_Get(t *testing.T) {
	cm := CaseMapping{mapping: make(map[string]string)}

	cm.Add("test", "Test")
	assert.Equal(t, "Test", cm.Get("test"))

	assert.Equal(t, "Unknown", cm.Get("unknown"))
}

func TestPascalCaseMapping(t *testing.T) {
	assert.Equal(t, "AdmissionRegistration", PascalCaseMapping.Get("admissionregistration"))
	assert.Equal(t, "Unknown", PascalCaseMapping.Get("unknown"))
}
