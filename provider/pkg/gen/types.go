// Copyright 2016-2020, Pulumi Corporation.
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
	"strings"

	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
)

type TemplateResource struct {
	Alias   string
	Name    string
	Package string
	Token   string
}

func (tr TemplateResource) GVK() string {
	parts := strings.Split(tr.Token, ":")
	contract.Assert(len(parts) == 3)
	gvk := parts[1] + "/" + parts[2]
	return strings.TrimPrefix(gvk, "core/")
}

func (tr TemplateResource) IsListKind() bool {
	return strings.HasSuffix(tr.Name, "List")
}

type TemplateResources struct {
	Resources []TemplateResource
	Imports   []string
}

func (tr TemplateResources) ListKinds() []TemplateResource {
	var resources []TemplateResource
	for _, r := range tr.Resources {
		if r.IsListKind() {
			resources = append(resources, r)
		}
	}
	return resources
}

func (tr TemplateResources) NonListKinds() []TemplateResource {
	var resources []TemplateResource
	for _, r := range tr.Resources {
		if !r.IsListKind() {
			resources = append(resources, r)
		}
	}
	return resources
}
