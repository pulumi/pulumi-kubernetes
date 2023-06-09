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

package main

import (
	"encoding/json"

	"github.com/imdario/mergo"
)

func mergeSwaggerSpecs(legacyBytes, currentBytes []byte) []byte {

	var legacyObj, newObj map[string]any
	err := json.Unmarshal(legacyBytes, &legacyObj)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(currentBytes, &newObj)
	if err != nil {
		panic(err)
	}
	err = mergo.Merge(&legacyObj, newObj, mergo.WithOverride)
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(legacyObj)
	if err != nil {
		panic(err)
	}

	return b
}
