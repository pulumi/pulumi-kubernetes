// Copyright 2016-2022, Pulumi Corporation.
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

package await

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/informers"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

// ------------------------------------------------------------------------------------------------

// Await logic for core/v1/Service.
//
// The goal of this code is to provide a fine-grained account of the status of a Kubernetes Service
// as it is being initialized. The idea is that if something goes wrong early, we want to alert the
// user so they can cancel the operation instead of waiting for timeout (~10 minutes).
//
// A Service can be one of several types, and the initialization behavior differs for each:
//
//   - If the type is `LoadBalancer`, the Service will be allocated a public IP address, and an
//     Endpoint object will be created, which specifies to which Pods traffic on different ports is
//     should be directed.
//   - If the type is `ClusterIP`, the Service is directly addressable only from inside the
//     cluster, so no public IP address will be allocated. An Endpoint object will still be created
//     to specify to which Pods traffic on different ports should be directed.
//
// The design of this awaiter is fundamentally an event loop on five channels:
//
//   1. The Service channel, to which the Kubernetes API server will proactively push every change
//      (additions, modifications, deletions) to any Service it knows about.
//   2. The Endpoint channel, which is the same idea as the Service channel, except it gets updates
//      to Endpoint objects.
//   3. A timeout channel, which fires after some minutes.
//   4. A cancellation channel, with which the user can signal cancellation (e.g., using SIGINT).
//   5. A "settled" channel, which is meant to fire a few seconds after any update to an Endpoint
//      object, so that we're sure they have time to "settle".
//
// The `serviceInitAwaiter` will synchronously process events from the union of all these channels.
// Any time the success conditions described above a reached, we will terminate the awaiter.
//
// The intermediate status we report tends to be related to whether endpoints are targeting > 0
// Pods. Because an external IP can take a long time to execute, we simply have to wait.
//
//
// x-refs:
//   * https://kubernetes.io/docs/tutorials/services/

// ------------------------------------------------------------------------------------------------

const (
	DefaultServiceTimeoutMins = 10
)

type serviceInitAwaiter struct {
	config           createAwaitConfig
	service          *unstructured.Unstructured
	serviceReady     bool
	endpointsReady   bool
	endpointsSettled bool
	serviceType      string
}

func makeServiceInitAwaiter(c createAwaitConfig) *serviceInitAwaiter {
	specType, _ := openapi.Pluck(c.currentOutputs.Object, "spec", "type")
	var t string
	if specTypeString, isString := specType.(string); isString {
		t = specTypeString
	} else {
		// The default value if `.spec.type` is not present.
		t = string(v1.ServiceTypeClusterIP)
	}

	return &serviceInitAwaiter{
		config:           c,
		service:          c.currentOutputs,
		serviceReady:     false,
		endpointsReady:   false,
		endpointsSettled: false,
		serviceType:      t,
	}
}

func awaitServiceInit(c createAwaitConfig) error {
	return makeServiceInitAwaiter(c).Await()
}

func awaitServiceRead(c createAwaitConfig) error {
	return makeServiceInitAwaiter(c).Read()
}

func awaitServiceUpdate(u updateAwaitConfig) error {
	return makeServiceInitAwaiter(u.createAwaitConfig).Await()
}

func (sia *serviceInitAwaiter) Await() error {
	//
	// We succeed only when all of the following are true:
	//
	//   1. Service object exists.
	//   2. Endpoint objects created. Each time we get an update, wait ~5-10 seconds
	//      after update to wait for any stragglers.
	//   3. The endpoints objects target some number of living objects.
	//   4. External IP address is allocated (if we're type `LoadBalancer`).
	//

	stopper := make(chan struct{})
	defer close(stopper)

	informerFactory := informers.NewInformerFactory(sia.config.clientSet,
		informers.WithNamespaceOrDefault(sia.config.currentOutputs.GetNamespace()))
	informerFactory.Start(stopper)

	serviceEvents := make(chan watch.Event)
	serviceInformer, err := informers.New(informerFactory, informers.ForServices(), informers.WithEventChannel(serviceEvents))
	if err != nil {
		return err
	}
	go serviceInformer.Informer().Run(stopper)

	endpointsEvents := make(chan watch.Event)
	endpointsInformer, err := informers.New(informerFactory, informers.ForEndpoints(), informers.WithEventChannel(endpointsEvents))
	if err != nil {
		return err
	}
	go endpointsInformer.Informer().Run(stopper)

	version := cluster.TryGetServerVersion(sia.config.clientSet.DiscoveryClientCached)

	timeout := sia.config.getTimeout(DefaultServiceTimeoutMins * 60)
	return sia.await(serviceEvents, endpointsEvents, time.After(timeout), make(chan struct{}), version)
}

func (sia *serviceInitAwaiter) Read() error {
	serviceClient, endpointsClient, err := sia.makeClients()
	if err != nil {
		return err
	}

	// Get live versions of Service and Endpoints.
	service, err := serviceClient.Get(sia.config.ctx,
		sia.config.currentOutputs.GetName(),
		metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the deployment as having been deleted.
		return err
	}

	//
	// In contrast to the case of `deployment`, an error getting the ReplicaSet or Pod lists does
	// not indicate that this resource was deleted, and we therefore should report the wrapped error
	// in a way that is useful to the user.
	//

	// Create endpoint watcher.
	endpointList, err := endpointsClient.List(sia.config.ctx, metav1.ListOptions{})
	if err != nil {
		logger.V(3).Infof("Error retrieving ReplicaSet list for Service %q: %v",
			service.GetName(), err)
		endpointList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	version := cluster.TryGetServerVersion(sia.config.clientSet.DiscoveryClientCached)

	return sia.read(service, endpointList, version)
}

func (sia *serviceInitAwaiter) read(
	service *unstructured.Unstructured, endpoints *unstructured.UnstructuredList,
	version cluster.ServerVersion,
) error {
	sia.processServiceEvent(watchAddedEvent(service))

	var err error
	settled := make(chan struct{})

	logger.V(3).Infof("Processing endpoint list: %#v", endpoints)
	err = endpoints.EachListItem(func(endpoint runtime.Object) error {
		sia.processEndpointEvent(watchAddedEvent(endpoint.(*unstructured.Unstructured)), settled)
		return nil
	})
	if err != nil {
		logger.V(3).Infof("Error iterating over ReplicaSet list for Deployment %q: %v",
			service.GetName(), err)
	}

	sia.endpointsSettled = true

	if sia.checkAndLogStatus() {
		return nil
	}

	return &initializationError{
		subErrors: sia.errorMessages(),
		object:    service,
	}
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (sia *serviceInitAwaiter) await(
	serviceEvents,
	endpointsEvents <-chan watch.Event,
	timeout <-chan time.Time,
	settled chan struct{},
	version cluster.ServerVersion,
) error {
	sia.config.logStatus(diag.Info, "[1/3] Finding Pods to direct traffic to")

	for {
		// Check whether we've succeeded.
		if sia.checkAndLogStatus() {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-sia.config.ctx.Done():
			// On cancel, check one last time if the service is ready.
			if sia.serviceReady && sia.endpointsReady {
				return nil
			}
			return &cancellationError{
				object:    sia.service,
				subErrors: sia.errorMessages(),
			}
		case <-timeout:
			// On timeout, check one last time if the service is ready.
			if sia.serviceReady && sia.endpointsReady {
				return nil
			}
			return &timeoutError{
				object:    sia.service,
				subErrors: sia.errorMessages(),
			}
		case <-settled:
			sia.endpointsSettled = true
		case event := <-serviceEvents:
			sia.processServiceEvent(event)
		case event := <-endpointsEvents:
			sia.processEndpointEvent(event, settled)
		}
	}
}

func (sia *serviceInitAwaiter) processServiceEvent(event watch.Event) {
	inputServiceName := sia.config.currentOutputs.GetName()

	service, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("Service watch received unknown object type %q",
			reflect.TypeOf(service))
		return
	}

	// Do nothing if this is not the service we're waiting for.
	if service.GetName() != inputServiceName {
		return
	}

	// Start with a blank slate.
	sia.serviceReady = false

	// Mark the service as not ready if it's deleted.
	if event.Type == watch.Deleted {
		return
	}

	sia.service = service

	if sia.serviceType == string(v1.ServiceTypeLoadBalancer) {
		// If it's type `LoadBalancer`, check whether an IP was allocated.
		lbIngress, _ := openapi.Pluck(service.Object, "status", "loadBalancer", "ingress")
		status, _ := openapi.Pluck(service.Object, "status")

		logger.V(3).Infof("Received status for service %q: %#v", inputServiceName, status)
		ing, isSlice := lbIngress.([]any)

		// Update status of service object so that we can check success.
		sia.serviceReady = isSlice && len(ing) > 0

		logger.V(3).Infof("Waiting for service %q to assign IP/hostname for a load balancer",
			inputServiceName)
	} else {
		// If it's not type `LoadBalancer`, report success.
		sia.serviceReady = true
	}
}

func (sia *serviceInitAwaiter) processEndpointEvent(event watch.Event, settledCh chan<- struct{}) {
	inputServiceName := sia.config.currentOutputs.GetName()

	// Get endpoint object.
	endpoint, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("Endpoint watch received unknown object type %q",
			reflect.TypeOf(endpoint))
		return
	}

	// Ignore if it's not one of the endpoint objects created by the service.
	//
	// NOTE: Because the client pool is per-namespace, the endpointName can be used as an
	// ID, as it's guaranteed by the API server to be unique.
	if endpoint.GetName() != inputServiceName {
		return
	}

	// Start over, prove that service is ready.
	sia.endpointsReady = false

	// Update status of endpoint objects so we can check success.
	if event.Type == watch.Added || event.Type == watch.Modified {
		subsets, hasTargets := openapi.Pluck(endpoint.Object, "subsets")
		targets, targetsIsSlice := subsets.([]any)
		endpointTargetsPod := hasTargets && targetsIsSlice && len(targets) > 0

		sia.endpointsReady = endpointTargetsPod
	} else if event.Type == watch.Deleted {
		sia.endpointsReady = false
	}

	// Every time we get an update to one of our endpoints objects, give it a few seconds
	// for them to settle.
	sia.endpointsSettled = false
	go func() {
		time.Sleep(10 * time.Second)
		settledCh <- struct{}{}
	}()
}

func (sia *serviceInitAwaiter) errorMessages() []string {
	messages := make([]string, 0)
	if sia.isHeadlessService() || sia.isExternalNameService() {
		return messages
	}

	if !sia.endpointsReady {
		messages = append(messages,
			"Service does not target any Pods. Selected Pods may not be ready, or "+ //nolint:goconst
				"field '.spec.selector' may not match labels on any Pods") //nolint:goconst
	}

	if sia.serviceType == string(v1.ServiceTypeLoadBalancer) && !sia.serviceReady {
		messages = append(messages,
			"Service was not allocated an IP address; does your cloud provider support this?")
	}

	return messages
}

// isHeadlessService checks if the Service has a defined .spec.clusterIP
func (sia *serviceInitAwaiter) isHeadlessService() bool {
	clusterIP, _ := openapi.Pluck(sia.service.Object, "spec", "clusterIP")
	return clusterIP == v1.ClusterIPNone
}

// isExternalNameService checks if the Service type is "ExternalName"
func (sia *serviceInitAwaiter) isExternalNameService() bool {
	return sia.serviceType == string(v1.ServiceTypeExternalName)
}

// shouldWaitForPods determines whether to wait for Pods to be ready before marking the Service ready.
func (sia *serviceInitAwaiter) shouldWaitForPods() bool {
	// For these special cases, skip the wait for Pod logic.
	if sia.isExternalNameService() || sia.isHeadlessService() {
		sia.endpointsReady = true
		return false
	}

	return true
}

func (sia *serviceInitAwaiter) checkAndLogStatus() bool {
	if !sia.shouldWaitForPods() {
		return sia.serviceReady
	}

	success := sia.serviceReady && sia.endpointsSettled && sia.endpointsReady
	if success {
		sia.config.logStatus(diag.Info,
			fmt.Sprintf("%sService initialization complete", cmdutil.EmojiOr("✅ ", "")))
	} else if sia.endpointsSettled && sia.endpointsReady {
		sia.config.logStatus(diag.Info, "[2/3] Attempting to allocate IP address to Service")
	}

	return success
}

func (sia *serviceInitAwaiter) makeClients() (
	serviceClient, endpointClient dynamic.ResourceInterface, err error,
) {
	serviceClient, err = clients.ResourceClient(
		kinds.Service, sia.config.currentOutputs.GetNamespace(), sia.config.clientSet)
	if err != nil {
		return nil, nil, errors.Wrapf(err,
			"Could not make client to read Service %q",
			sia.config.currentOutputs.GetName())
	}
	endpointClient, err = clients.ResourceClient(
		kinds.Endpoints, sia.config.currentOutputs.GetNamespace(), sia.config.clientSet)
	if err != nil {
		return nil, nil, errors.Wrapf(err,
			"Could not make client to read Endpoints associated with Service %q",
			sia.config.currentOutputs.GetName())
	}

	return serviceClient, endpointClient, nil
}
