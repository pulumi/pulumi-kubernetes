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

package metadata

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var dns1123Alphabet = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// AssignNameIfAutonamable generates a name for an object. Uses DNS-1123-compliant characters.
// All auto-named resources get the annotation `pulumi.com/autonamed` for tooling purposes.
func AssignNameIfAutonamable(obj *unstructured.Unstructured, base tokens.QName) {
	contract.Assert(base != "")
	if obj.GetName() == "" {
		obj.SetName(fmt.Sprintf("%s-%s", base, randString(8)))
		SetAnnotationTrue(obj, AnnotationAutonamed)
	}
}

// AdoptOldAutonameIfUnnamed checks if `newObj` has a name, and if not, "adopts" the name of `oldObj`
// instead. If `oldObj` was autonamed, then we mark `newObj` as autonamed, too.
func AdoptOldAutonameIfUnnamed(newObj, oldObj *unstructured.Unstructured) {
	contract.Assert(oldObj.GetName() != "")
	if newObj.GetName() == "" && IsAutonamed(oldObj) {
		newObj.SetName(oldObj.GetName())
		SetAnnotationTrue(newObj, AnnotationAutonamed)
	}
}

func IsAutonamed(obj *unstructured.Unstructured) bool {
	return IsAnnotationTrue(obj, AnnotationAutonamed)
}

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		// nolint:gosec
		b[i] = dns1123Alphabet[rand.Intn(len(dns1123Alphabet))]
	}
	return string(b)
}

// Seed RNG to get different random names at each suffix.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
