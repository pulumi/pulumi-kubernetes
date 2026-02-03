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

package gomega

import (
	. "github.com/onsi/gomega/gstruct" //nolint:golint,revive // gomega dsl
	gomegatypes "github.com/onsi/gomega/types"
)

func HaveSkipAwaitAnnotation() gomegatypes.GomegaMatcher {
	return MatchProps(IgnoreExtras, Props{
		"metadata": MatchObject(IgnoreExtras, Props{
			"annotations": MatchObject(IgnoreExtras, Props{
				"pulumi.com/skipAwait": MatchValue("true"),
			}),
		}),
	})
}
