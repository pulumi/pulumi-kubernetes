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
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

// ------------------------------------------------------------------------------------------------

// Await logic for extensions/v1beta1/Ingress.
//
// The goal of this code is to provide a fine-grained account of the status of a Kubernetes Ingress
// resource as it is being initialized. The idea is that if something goes wrong early, we want
// to alert the user so they can cancel the operation instead of waiting for timeout (~10 minutes).
//
// The design of this awaiter is fundamentally an event loop on four channels:
//
//   1. The Ingress channel, to which the Kubernetes API server will proactively push every change
//      (additions, modifications, deletions) to any Ingress it knows about.
//   2. The Endpoint channel, which is the same idea as the Ingress channel, except it gets updates
//      to Endpoint objects.
//   3. A timeout channel, which fires after some time.
//   4. A cancellation channel, with which the user can signal cancellation (e.g., using SIGINT).
//
// The `ingressInitAwaiter` will synchronously process events from the union of all these channels.
// Any time the success conditions described above a reached, we will terminate the awaiter.
//
// x-refs:
//   * https://github.com/nginxinc/kubernetes-ingress/blob/5847d1f3906287d2771f3767d61c15ac02522caa/docs/report-ingress-status.md

// ------------------------------------------------------------------------------------------------

const (
	DefaultIngressTimeoutMins = 10
)

type ingressInitAwaiter struct {
	config                    awaitConfig
	ingress                   *unstructured.Unstructured
	ingressReady              bool
	endpointsSettled          bool
	endpointEventsCount       uint64
	knownEndpointObjects      sets.Set[string]
	knownExternalNameServices sets.Set[string]
}

func makeIngressInitAwaiter(c awaitConfig) *ingressInitAwaiter {
	return &ingressInitAwaiter{
		config:                    c,
		ingress:                   c.currentOutputs,
		ingressReady:              false,
		endpointsSettled:          false,
		knownEndpointObjects:      sets.Set[string]{},
		knownExternalNameServices: sets.Set[string]{},
	}
}

func awaitIngressInit(c awaitConfig) error {
	return makeIngressInitAwaiter(c).Await()
}

func awaitIngressRead(c awaitConfig) error {
	return makeIngressInitAwaiter(c).Read()
}

func (iia *ingressInitAwaiter) Await() error {
	//
	// We succeed only when all of the following are true:
	//
	//   1.  Ingress object exists.
	//   2.  Endpoint objects exist with matching names for each Ingress path (except when Service
	//       type is ExternalName).
	//   3.  Ingress entry exists for .status.loadBalancer.ingress.
	//

	stopper := make(chan struct{})
	defer close(stopper)

	ingressEvents := make(chan watch.Event)
	ingressInformer, err := iia.config.factory.Subscribe(
		networkingv1.SchemeGroupVersion.WithResource("ingresses"),
		ingressEvents,
	)
	if err != nil {
		return err
	}
	defer ingressInformer.Close()

	endpointsEvents := make(chan watch.Event)
	endpointsInformer, err := iia.config.factory.Subscribe(
		corev1.SchemeGroupVersion.WithResource("endpoints"),
		endpointsEvents,
	)
	if err != nil {
		return err
	}
	defer endpointsInformer.Close()

	serviceEvents := make(chan watch.Event)
	serviceInformer, err := iia.config.factory.Subscribe(
		corev1.SchemeGroupVersion.WithResource("services"),
		serviceEvents,
	)
	if err != nil {
		return err
	}
	defer serviceInformer.Close()

	timeout := iia.config.getTimeout(DefaultIngressTimeoutMins * 60)
	return iia.await(ingressEvents, serviceEvents, endpointsEvents, make(chan struct{}), time.After(60*time.Second), time.After(timeout))
}

func (iia *ingressInitAwaiter) Read() error {
	ingressClient, endpointsClient, servicesClient, err := iia.makeClients()
	if err != nil {
		return err
	}

	// Get live versions of Ingress.
	ingress, err := ingressClient.Get(iia.config.ctx, iia.config.currentOutputs.GetName(), metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the deployment as having been deleted.
		return err
	}

	// Get live version of Endpoints.
	endpointList, err := endpointsClient.List(iia.config.ctx, metav1.ListOptions{})
	if err != nil {
		logger.V(3).Infof("Failed to list endpoints needed for Ingress awaiter: %v", err)
		endpointList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	serviceList, err := servicesClient.List(iia.config.ctx, metav1.ListOptions{})
	if err != nil {
		logger.V(3).Infof("Failed to list services needed for Ingress awaiter: %v", err)
		serviceList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	return iia.read(ingress, endpointList, serviceList)
}

func (iia *ingressInitAwaiter) read(ingress *unstructured.Unstructured, endpoints *unstructured.UnstructuredList,
	services *unstructured.UnstructuredList,
) error {
	iia.processIngressEvent(watchAddedEvent(ingress))

	err := services.EachListItem(func(service runtime.Object) error {
		iia.processServiceEvent(watchAddedEvent(service.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		logger.V(3).Infof("Error iterating over endpoint list for service %q: %v", ingress.GetName(), err)
	}

	settled := make(chan struct{})

	logger.V(3).Infof("Processing endpoint list: %#v", endpoints)
	err = endpoints.EachListItem(func(endpoint runtime.Object) error {
		iia.processEndpointEvent(watchAddedEvent(endpoint.(*unstructured.Unstructured)), settled)
		return nil
	})
	if err != nil {
		logger.V(3).Infof("Error iterating over endpoint list for ingress %q: %v", ingress.GetName(), err)
	}

	iia.endpointsSettled = true
	if iia.checkAndLogStatus() {
		return nil
	}

	return &initializationError{
		subErrors: iia.errorMessages(),
		object:    ingress,
	}
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (iia *ingressInitAwaiter) await(
	ingressEvents, serviceEvents, endpointsEvents <-chan watch.Event,
	settled chan struct{},
	settlementGracePeriodExpired <-chan time.Time,
	timeout <-chan time.Time,
) error {
	iia.config.logger.LogStatus(diag.Info, "[1/3] Finding a matching service for each Ingress path")

	for {
		// Check whether we've succeeded.
		if iia.checkAndLogStatus() {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-iia.config.ctx.Done():
			// On cancel, check one last time if the ingress is ready.
			if _, ready := iia.checkIfEndpointsReady(); ready && iia.ingressReady {
				return nil
			}
			return &cancellationError{
				object:    iia.ingress,
				subErrors: iia.errorMessages(),
			}
		case <-timeout:
			// On timeout, check one last time if the ingress is ready.
			if _, ready := iia.checkIfEndpointsReady(); ready && iia.ingressReady {
				return nil
			}
			return &timeoutError{
				object:    iia.ingress,
				subErrors: iia.errorMessages(),
			}
		case <-settlementGracePeriodExpired:
			// If we don't see any endpoint events in the designated time, assume endpoints have settled.
			// This is to account for the distinct possibility of ingress using a resource reference or non-existent
			// endpoints - in which case we will never see corresponding endpoint events.
			if iia.endpointEventsCount == 0 {
				iia.endpointsSettled = true
			}
		case <-settled:
			iia.endpointsSettled = true
		case event := <-ingressEvents:
			iia.processIngressEvent(event)
		case event := <-endpointsEvents:
			iia.processEndpointEvent(event, settled)
		case event := <-serviceEvents:
			iia.processServiceEvent(event)
		}
	}
}

func (iia *ingressInitAwaiter) processServiceEvent(event watch.Event) {
	service, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("Service watch received unknown object type %q",
			reflect.TypeOf(service))
		return
	}

	name := service.GetName()

	if event.Type == watch.Deleted {
		iia.knownExternalNameServices.Delete(name)
		return
	}

	t, ok := openapi.Pluck(service.Object, "spec", "type")
	if ok && t.(string) == "ExternalName" {
		iia.knownExternalNameServices.Insert(name)
	}
}

func (iia *ingressInitAwaiter) processIngressEvent(event watch.Event) {
	inputIngressName := iia.config.currentOutputs.GetName()

	ingress, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("Ingress watch received unknown object type %q",
			reflect.TypeOf(ingress))
		return
	}

	// Do nothing if this is not the ingress we're waiting for.
	if ingress.GetName() != inputIngressName {
		return
	}

	// Start with a blank slate.
	iia.ingressReady = false

	// Mark the ingress as not ready if it's deleted.
	if event.Type == watch.Deleted {
		return
	}

	iia.ingress = ingress

	// To the best of my knowledge, this works across all known ingress api version variations.
	ingressesRaw, ok := openapi.Pluck(ingress.Object, "status", "loadBalancer", "ingress")
	if !ok {
		logger.V(3).Infof("Unable to decode Ingress object from unstructured: %#v", ingress)
		return
	}

	ingresses, ok := ingressesRaw.([]any)
	if !ok {
		logger.V(3).Infof("Unexpected ingress object structure from unstructured: %#v", ingress)
		return
	}

	// Update status of ingress object so that we can check success.
	iia.ingressReady = len(ingresses) > 0

	logger.V(3).Infof("Waiting for ingress %q to update .status.loadBalancer with hostname/IP",
		inputIngressName)
}

func decodeIngress(u *unstructured.Unstructured, to any) error {
	b, err := u.MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)
	if err != nil {
		return err
	}
	return nil
}

func (iia *ingressInitAwaiter) checkIfEndpointsReady() (string, bool) {
	apiVersion := iia.ingress.GetAPIVersion()

	switch apiVersion {
	case "extensions/v1beta1", "networking.k8s.io/v1beta1":
		var obj networkingv1beta1.Ingress

		if err := decodeIngress(iia.ingress, &obj); err != nil {
			logger.V(3).Infof("Unable to decode Ingress object from unstructured: %#v", iia.ingress)
			return apiVersion, false
		}

		for _, rule := range obj.Spec.Rules {
			var httpIngressPaths []networkingv1beta1.HTTPIngressPath

			if rule.HTTP != nil {
				httpIngressPaths = rule.HTTP.Paths
			}
			for _, path := range httpIngressPaths {
				// Ignore ExternalName services
				if path.Backend.ServiceName != "" && iia.knownExternalNameServices.Has(path.Backend.ServiceName) {
					continue
				}

				if path.Backend.ServiceName != "" && !iia.knownEndpointObjects.Has(path.Backend.ServiceName) {
					if iia.endpointsSettled {
						// We haven't seen the target endpoint emit any events within the settlement period
						// and there is a chance it may never exist.
						iia.config.logger.LogStatus(diag.Warning, fmt.Sprintf("No matching service found for ingress rule: %s",
							expectedIngressPath(rule.Host, path.Path, path.Backend.ServiceName)))
					} else {
						// We may get more endpoint events, lets wait and retry.
						return apiVersion, false
					}
				}
			}
		}
	case "networking.k8s.io/v1":
		var obj networkingv1.Ingress
		if err := decodeIngress(iia.ingress, &obj); err != nil {
			logger.V(3).Infof("Unable to decode Ingress object from unstructured: %#v", iia.ingress)
			return apiVersion, false
		}

		for _, rule := range obj.Spec.Rules {
			var httpIngressPaths []networkingv1.HTTPIngressPath

			if rule.HTTP != nil {
				httpIngressPaths = rule.HTTP.Paths
			}
			for _, path := range httpIngressPaths {
				// TODO: Should we worry about "resource" backends?
				if path.Backend.Service == nil {
					continue
				}

				// Ignore ExternalName services
				if path.Backend.Service.Name != "" && iia.knownExternalNameServices.Has(path.Backend.Service.Name) {
					continue
				}

				if path.Backend.Service.Name != "" && !iia.knownEndpointObjects.Has(path.Backend.Service.Name) {
					if iia.endpointsSettled {
						// We haven't seen the target endpoint emit any events within the settlement period
						// and there is a chance it may never exist
						// (https://github.com/pulumi/pulumi-kubernetes/issues/1810)
						iia.config.logger.LogStatus(diag.Warning, fmt.Sprintf("No matching service found for ingress rule: %s",
							expectedIngressPath(rule.Host, path.Path, path.Backend.Service.Name)))
					} else {
						// We may get more endpoint events, lets wait and retry.
						return apiVersion, false
					}
				}
			}
		}
	}

	return apiVersion, true
}

// expectedIngressPath is a helper to print a useful error message.
func expectedIngressPath(host, path, serviceName string) string {
	rulePath := path
	if host != "" {
		rulePath = host + path
	}

	// It is valid for a user not to specify either a host or path [1]. In this case, any traffic not
	// matching another rule is routed to the specified Service for this rule. Print
	// `"" (default path)` to make this expectation clear to users.
	//
	// [1] https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#httpingresspath-v1beta1-extensions
	if rulePath == "" {
		rulePath = `"" (default path)`
	} else {
		rulePath = fmt.Sprintf("%q", rulePath)
	}

	// [host][path] -> serviceName
	return fmt.Sprintf("%s -> %q", rulePath, serviceName)
}

func (iia *ingressInitAwaiter) processEndpointEvent(event watch.Event, settledCh chan<- struct{}) {
	// Get endpoint object.
	endpoint, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("Endpoint watch received unknown object type %q",
			reflect.TypeOf(endpoint))
		return
	}

	iia.endpointEventsCount++

	name := endpoint.GetName()
	switch event.Type {
	case watch.Added, watch.Modified:
		iia.knownEndpointObjects.Insert(name)
	case watch.Deleted:
		iia.knownEndpointObjects.Delete(name)
		// NOTE: Unlike `processServiceEvent` don't return; we still want to set
		// `iia.endpointsSettled` to `false`.
	}

	// Every time we get an update to one of our endpoints objects, give it a few seconds
	// for them to settle.
	iia.endpointsSettled = false
	go func() {
		time.Sleep(10 * time.Second)
		settledCh <- struct{}{}
	}()
}

func (iia *ingressInitAwaiter) errorMessages() []string {
	messages := make([]string, 0)

	if _, ready := iia.checkIfEndpointsReady(); !ready {
		messages = append(messages, "Ingress has at least one rule with an unavailable target endpoint.")
	}

	if !iia.ingressReady {
		messages = append(messages,
			"Ingress .status.loadBalancer field was not updated with a hostname/IP address. "+ //nolint:goconst
				"\n    for more information about this error, see https://pulumi.io/xdv72s") //nolint:goconst
	}

	return messages
}

func (iia *ingressInitAwaiter) checkAndLogStatus() bool {
	_, ready := iia.checkIfEndpointsReady()
	success := iia.ingressReady && ready
	if success {
		iia.config.logger.LogStatus(diag.Info,
			fmt.Sprintf("%sIngress initialization complete", cmdutil.EmojiOr("✅ ", "")))
	} else if ready {
		iia.config.logger.LogStatus(diag.Info, "[2/3] Waiting for update of .status.loadBalancer with hostname/IP")
	}

	return success
}

func (iia *ingressInitAwaiter) makeClients() (
	ingressClient, endpointsClient, servicesClient dynamic.ResourceInterface, err error,
) {
	ingressClient, err = clients.ResourceClient(
		kinds.Ingress, iia.config.currentOutputs.GetNamespace(), iia.config.clientSet)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not make client to watch Ingress %q: %w",
			iia.config.currentOutputs.GetName(), err)
	}
	endpointsClient, err = clients.ResourceClient(
		kinds.Endpoints, iia.config.currentOutputs.GetNamespace(), iia.config.clientSet)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not make client to watch Endpoints associated with Ingress %q: %w",
			iia.config.currentOutputs.GetName(), err)
	}
	servicesClient, err = clients.ResourceClient(
		kinds.Service, iia.config.currentOutputs.GetNamespace(), iia.config.clientSet)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not make client to watch Services associated with Ingress %q: %w",
			iia.config.currentOutputs.GetName(), err)
	}

	return
}
