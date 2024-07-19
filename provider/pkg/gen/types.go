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

	"github.com/pulumi/pulumi/pkg/v3/codegen"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

// TemplateProperty holds information about a resource property that can be used to generate SDK overlays.
type TemplateProperty struct {
	ConstValue string // If set, the constant value of the property (e.g., "flowcontrol.apiserver.k8s.io/v1alpha1")
	Name       string // The name of the property (e.g., "FlowSchemaSpec")
	Package    string // The package path containing the property definition (e.g., "outputs.flowcontrol.v1alpha1")
}

// Type returns the property type. This could be either a constant value or the type definition.
func (tp TemplateProperty) Type() string {
	if len(tp.ConstValue) > 0 {
		return tp.ConstValue
	}
	return tp.Package
}

// TemplateResource holds information about a resource that can be used to generate SDK overlays.
type TemplateResource struct {
	Alias      string             // The optional alias to use for package imports (e.g., "flowcontrolv1alpha1")
	Name       string             // The name of the resource (e.g., "FlowSchema")
	Package    string             // The name of the package containing the resource definition (e.g., "flowcontrol.v1alpha1")
	Properties []TemplateProperty // Properties of the resource
	Token      string             // The schema token for the resource (e.g., "kubernetes:flowcontrol.apiserver.k8s.io/v1alpha1:FlowSchema")
}

// GVK returns the GroupVersionKind string for the k8s resource in the form "group/version/kind". The "core" group is
// rewritten to "", so resources in that Group look like "v1/Pod" rather than "core/v1/Pod".
func (tr TemplateResource) GVK() string {
	parts := strings.Split(tr.Token, ":")
	contract.Assertf(len(parts) == 3, "expected token to have three parts: %s", tr.Token)
	gvk := parts[1] + "/" + parts[2]
	return strings.TrimPrefix(gvk, "core/")
}

// IsListKind returns true if the resource name has the suffix "List".
func (tr TemplateResource) IsListKind() bool {
	return strings.HasSuffix(tr.Name, "List")
}

// TemplateResources holds a list of TemplateResource structs than can be filtered for codegen purposes.
type TemplateResources struct {
	Resources []TemplateResource
	Packages  []string
}

// ListKinds returns a sorted list of resources that are a List kind.
func (tr TemplateResources) ListKinds() []TemplateResource {
	var resources []TemplateResource
	for _, r := range tr.Resources {
		if r.IsListKind() {
			resources = append(resources, r)
		}
	}
	return resources
}

// ListKinds returns a sorted list of resources that are not a List kind.
func (tr TemplateResources) NonListKinds() []TemplateResource {
	var resources []TemplateResource
	for _, r := range tr.Resources {
		if !r.IsListKind() {
			resources = append(resources, r)
		}
	}
	return resources
}

// GoTemplateResources are TemplateResources specific to the Go SDK.
type GoTemplateResources struct {
	TemplateResources
}

// Imports returns a sorted list of all resource imports. This is currently used for the YAML SDK overlay.
func (tr GoTemplateResources) Imports() []string {
	imports := codegen.StringSet{}
	for _, r := range tr.Resources {
		importPath := fmt.Sprintf(`%s "%s"`, r.Alias, r.Package)
		imports.Add(importPath)
	}

	return imports.SortedValues()
}

// GVK holds sorted lists of GroupVersions and Kinds.
type GVK struct {
	GroupVersions []GroupVersion
	Kinds         []string
	PatchKinds    []string
	ListKinds     []string
}

// GroupVersion is the GroupVersion for a k8s resource.
type GroupVersion string

// GVConstName converts a GroupVersion to a name suitable for a const variable. This is used by codegen internal to
// the k8s provider.
// Example: apps/v1beta1 -> AppsV1B1
func (gv GroupVersion) GVConstName() string {
	parts := strings.Split(string(gv), "/")
	contract.Assertf(len(parts) == 2, "expected GroupVersion to have two parts: %s", gv)

	group, version := parts[0], parts[1]
	groupName := strings.Title(strings.SplitN(group, ".", 2)[0]) //nolint:staticcheck // Not unicode input.
	version = strings.Replace(version, "v", "V", -1)
	version = strings.Replace(version, "alpha", "A", -1)
	version = strings.Replace(version, "beta", "B", -1)

	return groupName + version
}
