package provider

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/pulumi/pulumi/pkg/util/contract"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const annotationInternalPrefix = "pulumi.com/"
const annotationInternalAutonamed = "pulumi.com/autonamed"

var dns1123Alphabet = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// assignName generates a name for an object. Uses DNS-1123-compliant characters. All auto-named
// resources get the annotation `pulumi.com/autonamed` for tooling purposes.
func assignNameIfAutonamable(obj *unstructured.Unstructured, base tokens.QName) {
	contract.Assert(base != "")
	if obj.GetName() == "" {
		obj.SetName(fmt.Sprintf("%s-%s", base, randString(8)))
		setAutonameAnnotation(obj)
	}
}

// adoptOldNameIfUnnamed checks if `newObj` has a name, and if not, "adopts" the name of `oldObj`
// instead. If `oldObj` was autonamed, then we mark `newObj` as autonamed, too.
func adoptOldNameIfUnnamed(newObj, oldObj *unstructured.Unstructured) {
	contract.Assert(oldObj.GetName() != "")
	if newObj.GetName() == "" {
		newObj.SetName(oldObj.GetName())
		if isAutonamed(oldObj) {
			setAutonameAnnotation(newObj)
		}
	}
}

func setAutonameAnnotation(obj *unstructured.Unstructured) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[annotationInternalAutonamed] = "true"
	obj.SetAnnotations(annotations)
}

func isAutonamed(obj *unstructured.Unstructured) bool {
	annotations := obj.GetAnnotations()
	autonamed := annotations[annotationInternalAutonamed]
	return autonamed == "true"
}

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = dns1123Alphabet[rand.Intn(len(dns1123Alphabet))]
	}
	return string(b)
}

// Seed RNG to get different random names at each suffix.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
