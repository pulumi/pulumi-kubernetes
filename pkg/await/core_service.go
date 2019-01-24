package await

import (
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/diag"
	"k8s.io/api/core/v1"
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

	serviceClient, endpointsClient, err := sia.makeClients()
	if err != nil {
		return err
	}

	// Create service watcher.
	serviceWatcher, err := serviceClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "Could set up watch for Service object '%s'",
			sia.config.currentInputs.GetName())
	}
	defer serviceWatcher.Stop()

	// Create endpoint watcher.
	endpointWatcher, err := endpointsClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for Endpoint objects associated with Service %q",
			sia.config.currentInputs.GetName())
	}
	defer endpointWatcher.Stop()

	version := ServerVersion(sia.config.clientSet.DiscoveryClientCached)

	return sia.await(serviceWatcher, endpointWatcher, time.After(10*time.Minute), make(chan struct{}), version)
}

func (sia *serviceInitAwaiter) Read() error {
	serviceClient, endpointsClient, err := sia.makeClients()
	if err != nil {
		return err
	}

	// Get live versions of Service and Endpoints.
	service, err := serviceClient.Get(sia.config.currentOutputs.GetName(),
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
	endpointList, err := endpointsClient.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Error retrieving ReplicaSet list for Service %q: %v",
			service.GetName(), err)
		endpointList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	version := ServerVersion(sia.config.clientSet.DiscoveryClientCached)

	return sia.read(service, endpointList, version)
}

func (sia *serviceInitAwaiter) read(service *unstructured.Unstructured, endpoints *unstructured.UnstructuredList,
	version serverVersion,
) error {
	sia.processServiceEvent(watchAddedEvent(service))

	var err error
	settled := make(chan struct{})

	glog.V(3).Infof("Processing endpoint list: %#v", endpoints)
	err = endpoints.EachListItem(func(endpoint runtime.Object) error {
		sia.processEndpointEvent(watchAddedEvent(endpoint.(*unstructured.Unstructured)), settled)
		return nil
	})
	if err != nil {
		glog.V(3).Infof("Error iterating over ReplicaSet list for Deployment %q: %v",
			service.GetName(), err)
	}

	sia.endpointsSettled = true

	if sia.checkAndLogStatus(version) {
		return nil
	}

	return &initializationError{
		subErrors: sia.errorMessages(),
		object:    service,
	}
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (sia *serviceInitAwaiter) await(
	serviceWatcher, endpointWatcher watch.Interface, timeout <-chan time.Time,
	settled chan struct{}, version serverVersion,
) error {
	sia.config.logStatus(diag.Info, "[1/3] Finding Pods to direct traffic to")

	for {
		// Check whether we've succeeded.
		if sia.checkAndLogStatus(version) {
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
		case event := <-serviceWatcher.ResultChan():
			sia.processServiceEvent(event)
		case event := <-endpointWatcher.ResultChan():
			sia.processEndpointEvent(event, settled)
		}
	}
}

func (sia *serviceInitAwaiter) processServiceEvent(event watch.Event) {
	inputServiceName := sia.config.currentInputs.GetName()

	service, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Service watch received unknown object type %q",
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

		glog.V(3).Infof("Received status for service %q: %#v", inputServiceName, status)
		ing, isSlice := lbIngress.([]interface{})

		// Update status of service object so that we can check success.
		sia.serviceReady = isSlice && len(ing) > 0

		glog.V(3).Infof("Waiting for service %q to assign IP/hostname for a load balancer",
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
		glog.V(3).Infof("Endpoint watch received unknown object type %q",
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
		targets, targetsIsSlice := subsets.([]interface{})
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
	if sia.emptyHeadlessOrExternalName() {
		return messages
	}

	if !sia.endpointsReady {
		messages = append(messages,
			"Service does not target any Pods. Selected Pods may not be ready, or "+
				"field '.spec.selector' may not match labels on any Pods")
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

// emptyHeadlessOrExternalName checks whether the current `Service` is either an "empty" headless
// `Service`[1] (i.e., it targets 0 `Pod`s) or a `Service` with `.spec.type: ExternalName` (which
// also targets 0 `Pod`s). This is useful to know when deciding whether to wait for a `Service` to
// target some number of `Pod`s.
//
// [1]: https://kubernetes.io/docs/concepts/services-networking/service/#headless-services
func (sia *serviceInitAwaiter) emptyHeadlessOrExternalName() bool {
	selectorI, _ := openapi.Pluck(sia.service.Object, "spec", "selector")
	selector, _ := selectorI.(map[string]interface{})

	headlessEmpty := len(selector) == 0 && sia.isHeadlessService()
	return headlessEmpty || sia.isExternalNameService()

}

// hasHeadlessServicePortBug checks whether the current `Service` is affected by a bug [1][2]
// that prevents endpoints from being populated when ports are not specified for a headless
// or external name Service.
//
// [1]: https://github.com/kubernetes/dns/issues/174
// [2]: https://github.com/kubernetes/kubernetes/commit/1c0137252465574519baf99252df8d75048f1304
func (sia *serviceInitAwaiter) hasHeadlessServicePortBug(version serverVersion) bool {
	// This bug only affects headless or external name Services.
	if !sia.isHeadlessService() && !sia.isExternalNameService() {
		return false
	}

	// k8s versions < 1.12.0 have the bug.
	if version.Compare(1, 12) < 0 {
		portsI, _ := openapi.Pluck(sia.service.Object, "spec", "ports")
		ports, _ := portsI.([]map[string]interface{})
		hasPorts := len(ports) > 0

		// The bug affects Services with no specified ports.
		if !hasPorts {
			return true
		}
	}

	return false
}

// shouldWaitForPods determines whether to wait for Pods to be ready before marking the Service ready.
func (sia *serviceInitAwaiter) shouldWaitForPods(version serverVersion) bool {
	// For these special cases, skip the wait for Pod logic.
	if sia.emptyHeadlessOrExternalName() || sia.hasHeadlessServicePortBug(version) {
		sia.endpointsReady = true
		return false
	}

	return true
}

func (sia *serviceInitAwaiter) checkAndLogStatus(version serverVersion) bool {
	if !sia.shouldWaitForPods(version) {
		return sia.serviceReady
	}

	success := sia.serviceReady && sia.endpointsSettled && sia.endpointsReady
	if success {
		sia.config.logStatus(diag.Info, "✅ Service initialization complete")
	} else if sia.endpointsSettled && sia.endpointsReady {
		sia.config.logStatus(diag.Info, "[2/3] Attempting to allocate IP address to Service")
	}

	return success
}

func (sia *serviceInitAwaiter) makeClients() (
	serviceClient, endpointClient dynamic.ResourceInterface, err error,
) {
	serviceClient, err = clients.ResourceClient(
		clients.Service, sia.config.currentOutputs.GetNamespace(), sia.config.clientSet)
	if err != nil {
		return nil, nil, errors.Wrapf(err,
			"Could not make client to watch Service %q",
			sia.config.currentOutputs.GetName())
	}
	endpointClient, err = clients.ResourceClient(
		clients.Endpoints, sia.config.currentOutputs.GetNamespace(), sia.config.clientSet)
	if err != nil {
		return nil, nil, errors.Wrapf(err,
			"Could not make client to watch Endpoints associated with Service %q",
			sia.config.currentOutputs.GetName())
	}

	return serviceClient, endpointClient, nil
}
