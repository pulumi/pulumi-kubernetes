/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helm

import (
	"context"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

// NewFakeTool creates a new Helm tool with faked execution.
func NewFakeTool(initActionConfig InitActionConfigF, execute ExecuteF) *Tool {
	settings := cli.New()
	if initActionConfig == nil {
		initActionConfig = FakeInit("default", nil)
	}
	if execute == nil {
		execute = FakeExecute("---\n", nil)
	}
	return &Tool{
		EnvSettings:      settings,
		HelmDriver:       "memory",
		initActionConfig: initActionConfig,
		execute:          execute,
	}
}

func FakeInit(namespace string, caps *chartutil.Capabilities) InitActionConfigF {
	return func(actionConfig *action.Configuration, namespaceOverride string) error {
		if namespaceOverride == "" {
			namespaceOverride = namespace
		}
		if caps == nil {
			caps = chartutil.DefaultCapabilities
		}
		actionConfig.Capabilities = caps
		return actionConfig.Init(nil, namespaceOverride, "memory", debug)
	}
}

func FakeExecute(manifest string, err error) ExecuteF {
	return func(ctx context.Context, i *action.Install, chrt *chart.Chart, vals map[string]interface{}) (*release.Release, error) {
		r := &release.Release{
			Name:      i.ReleaseName,
			Chart:     chrt,
			Info:      &release.Info{},
			Config:    vals,
			Namespace: i.Namespace,
			Manifest:  manifest,
		}
		return r, err
	}
}
