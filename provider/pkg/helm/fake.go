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
func NewFakeTool(initActionConfig InitActionConfigF, locateChart LocateChartF, execute ExecuteF) *Tool {
	settings := cli.New()
	if initActionConfig == nil {
		initActionConfig = FakeInitActionConfig("default", nil)
	}
	if locateChart == nil {
		locateChart = NewFakeLocator("./chart", nil).LocateChart
	}
	if execute == nil {
		execute = NewFakeExecutor().Execute
	}
	return &Tool{
		EnvSettings:      settings,
		HelmDriver:       "memory",
		initActionConfig: initActionConfig,
		locateChart:      locateChart,
		execute:          execute,
	}
}

func FakeInitActionConfig(namespace string, caps *chartutil.Capabilities) InitActionConfigF {
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

type FakeLocator struct {
	Path string
	Err  error

	action   *action.Install
	options  action.ChartPathOptions
	name     string
	settings *cli.EnvSettings
}

func (f *FakeLocator) Action() *action.Install {
	return f.action
}

func (f *FakeLocator) Name() string {
	return f.name
}

func (f *FakeLocator) Settings() *cli.EnvSettings {
	return f.settings
}

func (f *FakeLocator) LocateChart(i *action.Install, name string, settings *cli.EnvSettings) (string, error) {
	f.action = i
	f.name = name
	f.settings = settings
	return f.Path, f.Err
}

func NewFakeLocator(path string, err error) *FakeLocator {
	return &FakeLocator{
		Path: path,
		Err:  err,
	}
}

type FakeExecutor struct {
	action *action.Install
	chart  *chart.Chart
	values map[string]interface{}
}

func NewFakeExecutor() *FakeExecutor {
	return &FakeExecutor{}
}

func (f *FakeExecutor) Action() *action.Install {
	return f.action
}

func (f *FakeExecutor) Chart() *chart.Chart {
	return f.chart
}

func (f *FakeExecutor) Values() map[string]interface{} {
	return f.values
}

func (f *FakeExecutor) Execute(ctx context.Context, i *action.Install, chrt *chart.Chart, vals map[string]interface{}) (*release.Release, error) {
	f.action = i
	f.chart = chrt
	f.values = vals

	// force client-only mode
	oldDryRun := i.DryRun
	oldDryRunOption := i.DryRunOption
	oldClientOnly := i.ClientOnly
	defer func() {
		i.DryRun = oldDryRun
		i.DryRunOption = oldDryRunOption
		i.ClientOnly = oldClientOnly
	}()
	i.DryRun = true
	i.DryRunOption = "client"
	i.ClientOnly = true
	return i.RunWithContext(ctx, chrt, vals)
}
