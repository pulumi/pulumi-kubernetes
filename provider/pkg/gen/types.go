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
	"fmt"
	"strings"

	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
)

type TemplateProperty struct {
	ConstValue string
	Name       string
	Package    string
}

func (tp TemplateProperty) PackageOrConst() string {
	if len(tp.ConstValue) > 0 {
		return fmt.Sprintf("%q", tp.ConstValue)
	}
	return tp.Package
}

type TemplateResource struct {
	Alias      string
	Name       string
	Package    string
	Properties []TemplateProperty
	Token      string
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

type GVK struct {
	GroupVersions []GroupVersion
	Kinds         []string
}

type GroupVersion string

func (gv GroupVersion) GVConstName() string {
	parts := strings.Split(string(gv), "/")
	contract.Assert(len(parts) == 2)

	group, version := parts[0], parts[1]
	groupName := strings.Title(strings.SplitN(group, ".", 2)[0])
	version = strings.Replace(version, "v", "V", -1)
	version = strings.Replace(version, "alpha", "A", -1)
	version = strings.Replace(version, "beta", "B", -1)

	return groupName + version
}
