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
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	kubeversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/restmapper"
	kubetesting "k8s.io/client-go/testing"
	"k8s.io/kubectl/pkg/scheme"
)

var (
	DefaultServerVersion = kubeversion.Info{Major: "1", Minor: "29"}
)

type NewDynamicClientOption func(*newDynamicClientOptions)

func WithScheme(scheme *runtime.Scheme) NewDynamicClientOption {
	return func(options *newDynamicClientOptions) {
		options.Scheme = scheme
	}
}

func WithObjects(objects ...runtime.Object) NewDynamicClientOption {
	return func(options *newDynamicClientOptions) {
		options.Objects = append(options.Objects, objects...)
	}
}

func WithServerVersion(version kubeversion.Info) NewDynamicClientOption {
	return func(options *newDynamicClientOptions) {
		options.ServerVersion = version
	}
}

type newDynamicClientOptions struct {
	ServerVersion kubeversion.Info
	Scheme        *runtime.Scheme
	Objects       []runtime.Object
}

// NewSimpleDynamicClient creates a simple dynamic client for testing purposes,
// backed by fake discovery client and fake clientset as provided by the client-go library.
func NewSimpleDynamicClient(opts ...NewDynamicClientOption) (*clients.DynamicClientSet,
	*SimpleDiscovery, *SimpleRESTMapper, *SimpleDynamicClient) {
	options := newDynamicClientOptions{
		ServerVersion: DefaultServerVersion,
		Scheme:        scheme.Scheme,
		Objects:       []runtime.Object{},
	}
	for _, opt := range opts {
		opt(&options)
	}

	// make a fake discovery client that produces the core/v1 schema, and a mapper based on that.
	disco := NewSimpleDiscovery(options.ServerVersion)
	resources, err := restmapper.GetAPIGroupResources(disco)
	if err != nil {
		panic(err)
	}
	mapper := NewSimpleRESTMapper(resources)

	// make a fake clientset for testing purposes, backed by an testing.ObjectTracker with pre-populated objects.
	// see also: https://github.com/kubernetes/client-go/blob/kubernetes-1.29.0/examples/fake-client/main_test.go
	client := NewSimpleDynamicCLient(options.Scheme, options.Objects...)
	client.PrependReactor("create", "*", AdmitCreate())

	cs := &clients.DynamicClientSet{
		GenericClient:         client,
		DiscoveryClientCached: disco,
		RESTMapper:            mapper,
	}
	return cs, disco, mapper, client
}

func AdmitCreate() kubetesting.ReactionFunc {
	return func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
		if action, ok := action.(kubetesting.CreateAction); ok {
			objMeta, err := meta.Accessor(action.GetObject())
			if err != nil {
				return false, nil, err
			}

			// implement GenerateName since underlying object tracker doesn't natively support this.
			if objMeta.GetGenerateName() != "" && objMeta.GetName() == "" {
				name := objMeta.GetGenerateName() + "generated"
				objMeta.SetName(name)
			}
		}
		return false, nil, nil
	}
}
