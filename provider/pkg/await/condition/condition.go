// Copyright 2024, Pulumi Corporation.
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

package condition

import (
	"context"
	"fmt"
	"io"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

// Satisfier is an Observer which evaluates the observed object against some
// criteria.
type Satisfier interface {
	Observer

	// Satisfied returns true when the criteria is met.
	Satisfied() (bool, error)

	// Object returns the last-known state of the object being observed.
	Object() *unstructured.Unstructured
}

// logger allows injecting custom log behavior.
type logger interface {
	Log(diag.Severity, string)
	LogStatus(diag.Severity, string)
}

// logbuf logs messages to an io.Writter.
type logbuf struct{ w io.Writer }

func (l logbuf) Log(sev diag.Severity, msg string) {
	fmt.Fprintln(l.w, sev, msg)
}

func (l logbuf) LogStatus(sev diag.Severity, msg string) {
	l.Log(sev, msg)
}

// objectGetter allows injecting custom client behavior for fetching objects
// from the cluster.
type objectGetter interface {
	Get(
		ctx context.Context,
		name string,
		options metav1.GetOptions,
		subresources ...string,
	) (*unstructured.Unstructured, error)
}
