package await

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/pkg/client"
	"github.com/pulumi/pulumi/pkg/diag"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
// The design of this awaiter is fundamentally an event loop on three channels:
//
//   1. The Ingress channel, to which the Kubernetes API server will proactively push every change
//      (additions, modifications, deletions) to any Ingress it knows about.
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
	ingressPaths     map[string]*v1beta1.HTTPIngressPath
	ingressReady     bool
	ingressClient    *dynamic.ResourceInterface
	endpointsReady   bool
	endpointsSettled bool
	endpointClient   *dynamic.ResourceInterface
}

func makeIngressInitAwaiter(c createAwaitConfig) *ingressInitAwaiter {
	return &ingressInitAwaiter{
		config:           c,
		ingress:          c.currentOutputs,
		ingressPaths:     make(map[string]*v1beta1.HTTPIngressPath),
		ingressReady:     false,
		endpointsReady:   false,
		endpointsSettled: false,
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
	//   1. Ingress object exists.
	//   2. Endpoint objects exist with matching names for each Ingress path.
	//   3. Ingress entry exists for .status.loadBalancer.ingress.
	//

	if err := iia.initializeClients(); err != nil {
		glog.V(3).Infof("Failed to initialize client(s): %v", err)
		return errors.Wrapf(err, "Failed to initialize client(s) required for Ingress awaiter")
	}

	// Create ingress watcher.
	ingressWatcher, err := (*iia.ingressClient).Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "Could not set up watch for Ingress object %q",
			iia.config.currentInputs.GetName())
	}
	defer ingressWatcher.Stop()

	endpointWatcher, err := (*iia.endpointClient).Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for Endpoint objects associated with Ingress %q",
			iia.config.currentInputs.GetName())
	}
	defer endpointWatcher.Stop()

	return iia.await(ingressWatcher, endpointWatcher, time.After(10*time.Minute), make(chan struct{}))
}

func (iia *ingressInitAwaiter) initializeClients() error {
	if iia.ingressClient == nil {
		iia.ingressClient = &iia.config.clientForResource
	}

	if iia.endpointClient == nil {
		client, err := client.FromGVK(iia.config.pool, iia.config.disco, schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Endpoints",
		}, iia.config.currentInputs.GetNamespace())
		if err != nil {
			glog.V(3).Infof("Failed to initialize Endpoints client: %v", err)
			return err
		}

		iia.endpointClient = &client
	}

	return nil
}

func (iia *ingressInitAwaiter) Read() error {
	if err := iia.initializeClients(); err != nil {
		glog.V(3).Infof("Failed to initialize client(s): %v", err)
		return errors.Wrapf(err, "Failed to initialize client(s) required for Ingress awaiter")
	}

	// Get live versions of Ingress.
	ingress, err := (*iia.ingressClient).Get(iia.config.currentInputs.GetName(), metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the deployment as having been deleted.
		return err
	}

	// Create endpoint watcher.

	endpointList, err := (*iia.endpointClient).List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Failed to list endpoints needed for Ingress awaiter: %v", err)
		endpointList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	return iia.read(ingress, endpointList.(*unstructured.UnstructuredList))
}

func (iia *ingressInitAwaiter) read(
	ingress *unstructured.Unstructured,
	endpoints *unstructured.UnstructuredList,
) error {
	iia.processIngressEvent(watchAddedEvent(ingress))

	var err error
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
func (iia *ingressInitAwaiter) await(
	ingressWatcher, endpointWatcher watch.Interface, timeout <-chan time.Time,
	settled chan struct{},
) error {
	iia.config.logStatus(diag.Info, "[1/3] Waiting on endpoints for each Ingress path")
	iia.endpointsReady = iia.checkIfEndpointsReady()

	for {
		// Check whether we've succeeded.
		if iia.checkAndLogStatus() {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-iia.config.ctx.Done():
			// On cancel, check one last time if the ingress is ready.
			if iia.ingressReady {
				return nil
			}
			return &cancellationError{
				object:    iia.ingress,
				subErrors: iia.errorMessages(),
			}
		case <-timeout:
			// On timeout, check one last time if the ingress is ready.
			if iia.ingressReady {
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
		}
	}
}

func (iia *ingressInitAwaiter) processIngressEvent(event watch.Event) {
	inputIngressName := iia.config.currentInputs.GetName()

	ingress, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Ingress watch received unknown object type '%s'",
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

	glog.V(3).Infof("Received status for ingress '%s': %#v", inputIngressName, obj.Status)

	// Update status of ingress object so that we can check success.
	iia.ingressReady = len(obj.Status.LoadBalancer.Ingress) > 0

	glog.V(3).Infof("Waiting for ingress '%q' to assign IP/hostname for a load balancer",
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
			_, err := (*iia.endpointClient).Get(path.Backend.ServiceName, metav1.GetOptions{})
			if err != nil {
				if is404(err) {
					// No matching endpoint found for ingress rule
					iia.config.logStatus(diag.Error,
						fmt.Sprintf("No matching service found for ingress rule: %q", path.Path))
					return false
				}
			}
		}
	}

	return true
}

func (iia *ingressInitAwaiter) processEndpointEvent(event watch.Event, settledCh chan<- struct{}) {
	inputIngressName := iia.config.currentInputs.GetName()

	// Get endpoint object.
	endpoint, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Endpoint watch received unknown object type %q",
			reflect.TypeOf(endpoint))
		return
	}

	// Ignore if it's not one of the endpoint objects created by the ingress.
	//
	// NOTE: Because the client pool is per-namespace, the endpointName can be used as an
	// ID, as it's guaranteed by the API server to be unique.
	if endpoint.GetName() != inputIngressName {
		return
	}

	// Start over, prove that ingress is ready.
	iia.endpointsReady = false
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
	messages := []string{}

	if !iia.endpointsReady {
		messages = append(messages,
			"Ingress has at least one rule that does not target any Service. "+
				"Field '.spec.rules[].http.paths[].backend.serviceName' may not match any active Service")
	}

	return messages
}

func (iia *ingressInitAwaiter) checkAndLogStatus() bool {
	success := iia.ingressReady && iia.endpointsReady
	if success {
		iia.config.logStatus(diag.Info, "✅ Ingress initialization complete")
	} else if iia.endpointsReady {
		iia.config.logStatus(diag.Info, "[2/3] Waiting for load balancer to assign IP/hostname")
	}

	return success
}
