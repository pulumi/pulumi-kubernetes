// Copyright 2016-2018, Pulumi Corporation.
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

package client

import (
	"fmt"
	"strings"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

// --------------------------------------------------------------------------
// Names and namespaces.
// --------------------------------------------------------------------------

// FqObjName returns "namespace.name"
func FqObjName(o metav1.Object) string {
	return FqName(o.GetNamespace(), o.GetName())
}

// ParseFqName will parse a fully-qualified Kubernetes object name.
func ParseFqName(id string) (namespace, name string) {
	split := strings.Split(id, "/")
	if len(split) == 1 {
		return "", split[0]
	}
	namespace, name = split[0], split[1]
	return
}

// FqName returns "namespace/name"
func FqName(namespace, name string) string {
	if namespace == "" {
		return name
	}
	return fmt.Sprintf("%s/%s", namespace, name)
}

// NamespaceOrDefault returns `ns` or the the default namespace `"default"` if `ns` is empty.
func NamespaceOrDefault(ns string) string {
	if ns == "" {
		return "default"
	}
	return ns
}

// --------------------------------------------------------------------------
// Client utilities.
// --------------------------------------------------------------------------

// FromResource returns the ResourceClient for a given object
func FromResource(
	pool dynamic.ClientPool, disco discovery.ServerResourcesInterface, obj runtime.Object,
) (dynamic.ResourceInterface, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	meta, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	return FromGVK(pool, disco, gvk, NamespaceOrDefault(meta.GetNamespace()))
}

// FromGVK returns the ResourceClient for a given object
func FromGVK(
	pool dynamic.ClientPool, disco discovery.ServerResourcesInterface, gvk schema.GroupVersionKind,
	namespace string,
) (dynamic.ResourceInterface, error) {
	client, err := pool.ClientForGroupVersionKind(gvk)
	if err != nil {
		return nil, err
	}

	resource, err := serverResourceForGVK(disco, gvk)
	if err != nil {
		return nil, err
	}

	glog.V(3).Infof("Fetching client for %s namespace=%s", resource, namespace)
	rc := client.Resource(resource, namespace)
	return rc, nil
}

func serverResourceForGVK(
	disco discovery.ServerResourcesInterface, gvk schema.GroupVersionKind,
) (*metav1.APIResource, error) {
	resources, err := disco.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return nil, err
	}

	for _, r := range resources.APIResources {
		if r.Kind == gvk.Kind {
			glog.V(3).Infof("Using resource '%s' for %s", r.Name, gvk)
			return &r, nil
		}
	}

	return nil, fmt.Errorf("Server is unable to handle %s", gvk)
}
