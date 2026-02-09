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

package fake

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host"
)

// HostClient implements host.HostClient by forwarding all requests directly to the given pulumirpc.EngineServer.
type HostClient struct {
	Engine pulumirpc.EngineServer
}

var _ host.HostClient = &HostClient{}

func (h *HostClient) EngineConn() *grpc.ClientConn {
	return nil
}

func (h *HostClient) log(
	context context.Context, sev diag.Severity, urn resource.URN, msg string, ephemeral bool,
) error {
	var rpcsev pulumirpc.LogSeverity
	switch sev {
	case diag.Debug:
		rpcsev = pulumirpc.LogSeverity_DEBUG
	case diag.Info, diag.Infoerr:
		rpcsev = pulumirpc.LogSeverity_INFO
	case diag.Warning:
		rpcsev = pulumirpc.LogSeverity_WARNING
	case diag.Error:
		rpcsev = pulumirpc.LogSeverity_ERROR
	default:
		contract.Failf("Unrecognized log severity type: %v", sev)
	}
	if h.Engine != nil {
		_, err := h.Engine.Log(context, &pulumirpc.LogRequest{
			Severity:  rpcsev,
			Message:   strings.ToValidUTF8(msg, "�"),
			Urn:       string(urn),
			Ephemeral: ephemeral,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Log logs a global message, including errors and warnings.
func (h *HostClient) Log(
	context context.Context, sev diag.Severity, urn resource.URN, msg string,
) error {
	return h.log(context, sev, urn, msg, false)
}

// LogStatus logs a global status message, including errors and warnings. Status messages will
// appear in the `Info` column of the progress display, but not in the final output.
func (h *HostClient) LogStatus(
	context context.Context, sev diag.Severity, urn resource.URN, msg string,
) error {
	return h.log(context, sev, urn, msg, true)
}
