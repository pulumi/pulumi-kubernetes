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

package host

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"golang.org/x/net/context"
)

// HostClient is the interface that the provider uses to communicate with the Pulumi engine.
//
//nolint:golint
type HostClient interface {
	// Log logs a global message, including errors and warnings.
	Log(context context.Context, sev diag.Severity, urn resource.URN, msg string) error

	// LogStatus logs a global status message, including errors and warnings. Status messages will
	// appear in the `Info` column of the progress display, but not in the final output.
	LogStatus(context context.Context, sev diag.Severity, urn resource.URN, msg string) error
}
