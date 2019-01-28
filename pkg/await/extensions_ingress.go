package await

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/diag"
	"k8s.io/api/extensions/v1beta1"
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

type ingressInitAwaiter struct {
	config                    createAwaitConfig
	ingress                   *unstructured.Unstructured
	ingressReady              bool
	endpointsSettled          bool
	knownEndpointObjects      sets.String
	knownExternalNameServices sets.String
}

func makeIngressInitAwaiter(c createAwaitConfig) *ingressInitAwaiter {
	return &ingressInitAwaiter{
		config:                    c,
		ingress:                   c.currentOutputs,
		ingressReady:              false,
		endpointsSettled:          false,
		knownEndpointObjects:      sets.NewString(),
		knownExternalNameServices: sets.NewString(),
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
	//   2.  Endpoint objects exist with matching names for each Ingress path (except when Service
	//       type is ExternalName).
	//   3.  Ingress entry exists for .status.loadBalancer.ingress.
	//

	ingressClient, endpointsClient, servicesClient, err := iia.makeClients()
	if err != nil {
		return err
	}

	// Create ingress watcher.
	ingressWatcher, err := ingressClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "Could not set up watch for Ingress object %q",
			iia.config.currentInputs.GetName())
	}
	defer ingressWatcher.Stop()

	endpointWatcher, err := endpointsClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for Endpoint objects associated with Ingress %q",
			iia.config.currentInputs.GetName())
	}
	defer endpointWatcher.Stop()

	serviceWatcher, err := servicesClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for Service objects associated with Ingress %q",
			iia.config.currentInputs.GetName())
	}

	return iia.await(ingressWatcher, serviceWatcher, endpointWatcher, make(chan struct{}), time.After(10*time.Minute))
}

func (iia *ingressInitAwaiter) Read() error {
	ingressClient, endpointsClient, servicesClient, err := iia.makeClients()
	if err != nil {
		return err
	}

	// Get live versions of Ingress.
	ingress, err := ingressClient.Get(iia.config.currentInputs.GetName(), metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the deployment as having been deleted.
		return err
	}

	// Get live version of Endpoints.
	endpointList, err := endpointsClient.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Failed to list endpoints needed for Ingress awaiter: %v", err)
		endpointList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	serviceList, err := servicesClient.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Failed to list services needed for Ingress awaiter: %v", err)
		serviceList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	return iia.read(ingress, endpointList, serviceList)
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
			if iia.ingressReady && iia.checkIfEndpointsReady() {
				return nil
			}
			return &cancellationError{
				object:    iia.ingress,
				subErrors: iia.errorMessages(),
			}
		case <-timeout:
			// On timeout, check one last time if the ingress is ready.
			if iia.ingressReady && iia.checkIfEndpointsReady() {
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

	glog.V(3).Infof("Received status for ingress %q: %#v", inputIngressName, obj.Status)

	// Update status of ingress object so that we can check success.
	iia.ingressReady = len(obj.Status.LoadBalancer.Ingress) > 0

	glog.V(3).Infof("Waiting for ingress %q to update .status.loadBalancer with hostname/IP",
		inputIngressName)
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
			if iia.knownExternalNameServices.Has(path.Backend.ServiceName) {
				continue
			}

			if !iia.knownEndpointObjects.Has(path.Backend.ServiceName) {
				iia.config.logStatus(diag.Info, fmt.Sprintf("No matching service found for ingress rule: %s",
					expectedIngressPath(rule.Host, path.Path, path.Backend.ServiceName)))

				return false
			}
		}
	}

	return true
}

// expectedIngressPath is a helper to print a useful error message.
func expectedIngressPath(host, path, serviceName string) string {
	rulePath := path
	if host != "" {
		rulePath = host + path
	}

	// It is valid for a user not to specify either a host or path [1]. In this case, any traffic not
	// matching another rule is routed to the specified Service for this rule. Print <default> to make
	// this expectation clear to users.
	//
	// [1] https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#httpingresspath-v1beta1-extensions
	if rulePath == "" {
		rulePath = "<default>"
	}

	// [host][path] -> serviceName
	return fmt.Sprintf("%q -> %q", rulePath, serviceName)
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

	if !iia.checkIfEndpointsReady() {
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
	success := iia.ingressReady && iia.checkIfEndpointsReady()
	if success {
		iia.config.logStatus(diag.Info, "✅ Ingress initialization complete")
	} else if iia.checkIfEndpointsReady() {
		iia.config.logStatus(diag.Info, "[2/3] Waiting for update of .status.loadBalancer with hostname/IP")
	}

	return success
}

func (iia *ingressInitAwaiter) makeClients() (
	ingressClient, endpointsClient, servicesClient dynamic.ResourceInterface, err error,
) {
	ingressClient, err = clients.ResourceClient(
		clients.Ingress, iia.config.currentInputs.GetNamespace(), iia.config.clientSet)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err,
			"Could not make client to watch Ingress %q",
			iia.config.currentInputs.GetName())
	}
	endpointsClient, err = clients.ResourceClient(
		clients.Endpoints, iia.config.currentInputs.GetNamespace(), iia.config.clientSet)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err,
			"Could not make client to watch Endpoints associated with Ingress %q",
			iia.config.currentInputs.GetName())
	}
	servicesClient, err = clients.ResourceClient(
		clients.Service, iia.config.currentInputs.GetNamespace(), iia.config.clientSet)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err,
			"Could not make client to watch Services associated with Ingress %q",
			iia.config.currentInputs.GetName())
	}

	return
}
