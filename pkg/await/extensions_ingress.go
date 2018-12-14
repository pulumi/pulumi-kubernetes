package await

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/pkg/client"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/diag"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
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

type ingressInitAwaiter struct {
	config           createAwaitConfig
	ingress          *unstructured.Unstructured
	ingressReady     bool
	endpointsReady   bool
	endpointsSettled bool
	endpointExists   map[string]bool
	externalServices map[string]bool
}

func makeIngressInitAwaiter(c createAwaitConfig) *ingressInitAwaiter {
	return &ingressInitAwaiter{
		config:           c,
		ingress:          c.currentOutputs,
		ingressReady:     false,
		endpointsReady:   false,
		endpointsSettled: false,
		endpointExists:   make(map[string]bool),
		externalServices: make(map[string]bool),
	}
}

func awaitIngressInit(c createAwaitConfig) error {
	return makeIngressInitAwaiter(c).Await()
}

func awaitIngressRead(c createAwaitConfig) error {
	return makeIngressInitAwaiter(c).Read()
}

func awaitIngressUpdate(u updateAwaitConfig) error {
	return makeIngressInitAwaiter(u.createAwaitConfig).Await()
}

func (iia *ingressInitAwaiter) Await() error {
	//
	// We succeed only when all of the following are true:
	//
	//   1.  Ingress object exists.
	//   2.  Endpoint objects exist with matching names for each Ingress path.
	//	 2.1 Alternatively, a Service with type: ExternalName must path the Ingress path.
	//   3.  Ingress entry exists for .status.loadBalancer.ingress.
	//

	// Create ingress watcher.
	ingressWatcher, err := iia.config.clientForResource.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "Could not set up watch for Ingress object %q",
			iia.config.currentInputs.GetName())
	}
	defer ingressWatcher.Stop()

	endpointsClient, err := client.FromGVK(iia.config.pool, iia.config.disco, schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Endpoints",
	}, iia.config.currentInputs.GetNamespace())
	if err != nil {
		glog.V(3).Infof("Failed to initialize Endpoints client: %v", err)
		return err
	}
	endpointWatcher, err := endpointsClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for Endpoint objects associated with Ingress %q",
			iia.config.currentInputs.GetName())
	}
	defer endpointWatcher.Stop()

	servicesClient, err := client.FromGVK(iia.config.pool, iia.config.disco, schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Service",
	}, iia.config.currentInputs.GetNamespace())
	if err != nil {
		glog.V(3).Infof("Failed to initialize Services client: %v", err)
		return err
	}
	serviceWatcher, err := servicesClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for Service objects associated with Ingress %q",
			iia.config.currentInputs.GetName())
	}

	return iia.await(ingressWatcher, serviceWatcher, endpointWatcher, make(chan struct{}), time.After(10*time.Minute))
}

func (iia *ingressInitAwaiter) Read() error {
	// Get live versions of Ingress.
	ingress, err := iia.config.clientForResource.Get(iia.config.currentInputs.GetName(), metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the deployment as having been deleted.
		return err
	}

	// Get live version of Endpoints.
	endpointsClient, err := client.FromGVK(iia.config.pool, iia.config.disco, schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Endpoints",
	}, iia.config.currentInputs.GetNamespace())
	if err != nil {
		glog.V(3).Infof("Failed to initialize Endpoints client: %v", err)
		return err
	}
	endpointList, err := endpointsClient.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Failed to list endpoints needed for Ingress awaiter: %v", err)
		endpointList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	servicesClient, err := client.FromGVK(iia.config.pool, iia.config.disco, schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Service",
	}, iia.config.currentInputs.GetNamespace())
	if err != nil {
		glog.V(3).Infof("Failed to initialize Services client: %v", err)
		return err
	}
	serviceList, err := servicesClient.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Failed to list services needed for Ingress awaiter: %v", err)
		serviceList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	return iia.read(ingress, endpointList.(*unstructured.UnstructuredList), serviceList.(*unstructured.UnstructuredList))
}

func (iia *ingressInitAwaiter) read(ingress *unstructured.Unstructured, endpoints *unstructured.UnstructuredList,
	services *unstructured.UnstructuredList) error {
	iia.processIngressEvent(watchAddedEvent(ingress))

	err := services.EachListItem(func(service runtime.Object) error {
		iia.processServiceEvent(watchAddedEvent(service.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		glog.V(3).Infof("Error iterating over endpoint list for service %q: %v", ingress.GetName(), err)
	}

	settled := make(chan struct{})

	glog.V(3).Infof("Processing endpoint list: %#v", endpoints)
	err = endpoints.EachListItem(func(endpoint runtime.Object) error {
		iia.processEndpointEvent(watchAddedEvent(endpoint.(*unstructured.Unstructured)), settled)
		return nil
	})
	if err != nil {
		glog.V(3).Infof("Error iterating over endpoint list for ingress %q: %v", ingress.GetName(), err)
	}

	iia.endpointsReady = iia.checkIfEndpointsReady()
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
func (iia *ingressInitAwaiter) await(ingressWatcher, serviceWatcher, endpointWatcher watch.Interface,
	settled chan struct{}, timeout <-chan time.Time) error {
	iia.config.logStatus(diag.Info, "[1/3] Finding a matching service for each Ingress path")

	for {
		// Check whether we've succeeded.
		if iia.checkAndLogStatus() {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-iia.config.ctx.Done():
			// On cancel, check one last time if the ingress is ready.
			if iia.ingressReady && iia.endpointsReady {
				return nil
			}
			return &cancellationError{
				object:    iia.ingress,
				subErrors: iia.errorMessages(),
			}
		case <-timeout:
			// On timeout, check one last time if the ingress is ready.
			iia.endpointsReady = iia.checkIfEndpointsReady()
			if iia.ingressReady && iia.endpointsReady {
				return nil
			}
			return &timeoutError{
				object:    iia.ingress,
				subErrors: iia.errorMessages(),
			}
		case <-settled:
			iia.endpointsSettled = true
		case event := <-ingressWatcher.ResultChan():
			iia.processIngressEvent(event)
		case event := <-endpointWatcher.ResultChan():
			iia.processEndpointEvent(event, settled)
		case event := <-serviceWatcher.ResultChan():
			iia.processServiceEvent(event)
		}
	}
}

func (iia *ingressInitAwaiter) processServiceEvent(event watch.Event) {
	service, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Service watch received unknown object type %q",
			reflect.TypeOf(service))
		return
	}

	name := service.GetName()

	t, ok := openapi.Pluck(service.Object, "spec", "type")
	if ok && t.(string) == "ExternalName" {
		iia.externalServices[name] = true
	} else {
		iia.externalServices[name] = false
	}
}

func (iia *ingressInitAwaiter) processIngressEvent(event watch.Event) {
	inputIngressName := iia.config.currentInputs.GetName()

	ingress, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Ingress watch received unknown object type %q",
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
	obj, err := decodeIngress(ingress)
	if err != nil {
		glog.V(3).Infof("Unable to decode Ingress object from unstructured: %#v", ingress)
		return
	}

	var serviceNames []string
	for _, rule := range obj.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			serviceNames = append(serviceNames, path.Backend.ServiceName)
		}
	}
	iia.ignoreExternalNameServices(serviceNames)

	iia.endpointsReady = iia.checkIfEndpointsReady()

	glog.V(3).Infof("Received status for ingress %q: %#v", inputIngressName, obj.Status)

	// Update status of ingress object so that we can check success.
	iia.ingressReady = len(obj.Status.LoadBalancer.Ingress) > 0

	glog.V(3).Infof("Waiting for ingress %q to update .status.loadBalancer with hostname/IP",
		inputIngressName)
}

func (iia *ingressInitAwaiter) ignoreExternalNameServices(names []string) {
	// Services with type: ExternalName do not have associated Pods/Endpoints to wait for, so mark as ready.
	for _, name := range names {
		if iia.externalServices[name] {
			iia.endpointExists[name] = true
		}
	}
}

func decodeIngress(u *unstructured.Unstructured) (*v1beta1.Ingress, error) {
	b, err := u.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var obj v1beta1.Ingress
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func (iia *ingressInitAwaiter) checkIfEndpointsReady() bool {
	obj, err := decodeIngress(iia.ingress)
	if err != nil {
		glog.V(3).Infof("Unable to decode Ingress object from unstructured: %#v", iia.ingress)
		return false
	}

	for _, rule := range obj.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			// Ignore ExternalName services
			if iia.externalServices[path.Backend.ServiceName] {
				continue
			}

			if !iia.endpointExists[path.Backend.ServiceName] {
				iia.config.logStatus(diag.Error,
					fmt.Sprintf("No matching service found for ingress rule: %q", path.Path))
				return false
			}
		}
	}

	return true
}

func (iia *ingressInitAwaiter) processEndpointEvent(event watch.Event, settledCh chan<- struct{}) {
	// Get endpoint object.
	endpoint, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Endpoint watch received unknown object type %q",
			reflect.TypeOf(endpoint))
		return
	}

	name := endpoint.GetName()
	switch event.Type {
	case watch.Added, watch.Modified:
		iia.endpointExists[name] = true
	case watch.Deleted:
		iia.endpointExists[name] = false
	}

	// Start over, prove that endpoints are ready.
	iia.endpointsReady = iia.checkIfEndpointsReady()

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

	if !iia.endpointsReady {
		messages = append(messages,
			"Ingress has at least one rule that does not target any Service. "+
				"Field '.spec.rules[].http.paths[].backend.serviceName' may not match any active Service")
	}

	if !iia.ingressReady {
		messages = append(messages,
			"Ingress .status.loadBalancer field was not updated with a hostname/IP address. "+
				"\n    for more information about this error, see https://pulumi.io/xdv72s")
	}

	return messages
}

func (iia *ingressInitAwaiter) checkAndLogStatus() bool {
	success := iia.ingressReady && iia.endpointsReady
	if success {
		iia.config.logStatus(diag.Info, "✅ Ingress initialization complete")
	} else if iia.endpointsReady {
		iia.config.logStatus(diag.Info, "[2/3] Waiting for update of .status.loadBalancer with hostname/IP")
	}

	return success
}
