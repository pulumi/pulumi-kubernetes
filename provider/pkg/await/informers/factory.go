// Copyright 2021, Pulumi Corporation.  All rights reserved.

package informers

import (
	"time"

	"github.com/pulumi/pulumi-kubernetes/v4/provider/pkg/clients"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

const defaultResyncInterval = 60 * time.Second

type informerFactoryOptions struct {
	namespace        string
	resyncInterval   time.Duration
	tweakListOptions dynamicinformer.TweakListOptionsFunc
}

type InformerFactoryOption interface {
	apply(*informerFactoryOptions)
}

type applyInformerFactoryOptionFunc func(*informerFactoryOptions)

func (a applyInformerFactoryOptionFunc) apply(o *informerFactoryOptions) {
	a(o)
}

// WithNamespace configures the namespace for the informer factory created by NewInformerFactory.
func WithNamespace(namespace string) InformerFactoryOption {
	return applyInformerFactoryOptionFunc(func(o *informerFactoryOptions) {
		o.namespace = namespace
	})
}

// WithNamespaceOrDefault configures the namespace for the informer factory similar to WithNamespace,
// except the empty ("") namespace is interpreted as the "default" namespace.
func WithNamespaceOrDefault(namespace string) InformerFactoryOption {
	if namespace == "" {
		return WithNamespace(metav1.NamespaceDefault)
	}
	return WithNamespace(namespace)
}

// WithResyncInterval overrides the default resync interval of 60 seconds. Setting this to 0 will disable resync.
// Refer to cache.NewSharedIndexInformer for caveats on how resync intervals are honored.
func WithResyncInterval(interval time.Duration) InformerFactoryOption {
	return applyInformerFactoryOptionFunc(func(o *informerFactoryOptions) {
		o.resyncInterval = interval
	})
}

// WithTweakListOptionsFunc allows customizing options used by the informer's underlying cache.Lister.
// By default, no customizations are performed.
func WithTweakListOptionsFunc(tweakListOptionsFunc dynamicinformer.TweakListOptionsFunc) InformerFactoryOption {
	return applyInformerFactoryOptionFunc(func(o *informerFactoryOptions) {
		o.tweakListOptions = tweakListOptionsFunc
	})
}

// NewInformerFactory is a convenient wrapper around initializing an informer factory for the purposes of
// awaiting Kubernetes resources for the provider.
// By default, the informer is configured for the "default" namespace. This can be overridden by WithNamespace.
// By default, the resync interval is configured for 60 seconds. This can be overridden by WithResyncInterval.
// By default, no constraints are placed on the listers. This can be overridden through WithTweakListOptionsFunc.
func NewInformerFactory(
	clientSet *clients.DynamicClientSet,
	opts ...InformerFactoryOption,
) dynamicinformer.DynamicSharedInformerFactory {
	o := informerFactoryOptions{
		namespace:      metav1.NamespaceDefault,
		resyncInterval: defaultResyncInterval,
	}

	for _, opt := range opts {
		opt.apply(&o)
	}

	return dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		clientSet.GenericClient, o.resyncInterval, o.namespace, o.tweakListOptions)
}
