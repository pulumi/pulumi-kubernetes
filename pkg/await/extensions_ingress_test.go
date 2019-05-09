package await

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

func Test_Extensions_Ingress(t *testing.T) {
	tests := []struct {
		description   string
		ingressInput  func(namespace, name, targetService string) *unstructured.Unstructured
		do            func(ingresses, services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time)
		expectedError error
	}{
		{
			description:  "Should succeed when Ingress is allocated an IP address and all paths match an existing Endpoint",
			ingressInput: initializedIngress,
			do: func(ingresses, services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes initialized ingress and endpoint objects back.
				ingresses <- watchAddedEvent(initializedIngress("default", "foo", "foo-4setj4y6"))
				endpoints <- watchAddedEvent(initializedEndpoint("default", "foo-4setj4y6"))

				// Mark endpoint objects as having settled. Success.
				settled <- struct{}{}
			},
		},
		{
			description:  "Should succeed when Ingress is allocated an IP address and path references an ExternalName Service",
			ingressInput: initializedIngress,
			do: func(ingresses, services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes initialized ingress and endpoint objects back.
				ingresses <- watchAddedEvent(initializedIngress("default", "foo", "foo-4setj4y6"))
				services <- watchAddedEvent(externalNameService("default", "foo-4setj4y6"))

				// Mark endpoint objects as having settled. Success.
				settled <- struct{}{}

				// Timeout, success.
				timeout <- time.Now()
			},
		},
		{
			description:  "Should fail if the Ingress does not have an IP address allocated, and not all paths match an existing Endpoint",
			ingressInput: ingressInput,
			do: func(ingresses, services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// Trigger timeout.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: ingressInput("default", "foo", "foo-4setj4y6"),
				subErrors: []string{
					"Ingress has at least one rule that does not target any Service. " +
						"Field '.spec.rules[].http.paths[].backend.serviceName' may not match any active Service",
					"Ingress .status.loadBalancer field was not updated with a hostname/IP address. " +
						"\n    for more information about this error, see https://pulumi.io/xdv72s",
				}},
		},
		{
			description:  "Should fail if not all Ingress paths match existing Endpoints",
			ingressInput: ingressInput,
			do: func(ingresses, services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes initialized ingress back.
				ingresses <- watchAddedEvent(initializedIngress("default", "foo", "foo-4setj4y6"))

				settled <- struct{}{}

				// Finally, time out.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: initializedIngress("default", "foo", "foo-4setj4y6"),
				subErrors: []string{
					"Ingress has at least one rule that does not target any Service. " +
						"Field '.spec.rules[].http.paths[].backend.serviceName' may not match any active Service"}},
		},
		{
			description:  "Should succeed for Ingress with an unspecified path.",
			ingressInput: ingressInput,
			do: func(ingresses, services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes initialized ingress back.
				ingresses <- watchAddedEvent(initializedIngressUnspecifiedPath("default", "foo", "foo-4setj4y6"))
				endpoints <- watchAddedEvent(initializedEndpoint("default", "foo-4setj4y6"))

				settled <- struct{}{}

				// Finally, time out.
				timeout <- time.Now()
			},
		},
		{
			description:  "Should fail if Ingress is not allocated an IP address",
			ingressInput: ingressInput,
			do: func(ingresses, services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes uninitialized service back.
				ingresses <- watchAddedEvent(ingressInput("default", "foo", "foo-4setj4y6"))

				// Pass initialized endpoint objects. Mark them as settled.
				endpoints <- watchAddedEvent(
					initializedEndpoint("default", "foo-4setj4y6"))
				settled <- struct{}{}

				// Finally, time out.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: ingressInput("default", "foo", "foo-4setj4y6"),
				subErrors: []string{
					"Ingress .status.loadBalancer field was not updated with a hostname/IP address. " +
						"\n    for more information about this error, see https://pulumi.io/xdv72s",
				}},
		},
	}

	for _, test := range tests {
		awaiter := makeIngressInitAwaiter(
			mockAwaitConfig(test.ingressInput("default", "foo", "foo-4setj4y6")))

		ingresses := make(chan watch.Event)
		services := make(chan watch.Event)
		endpoints := make(chan watch.Event)
		settled := make(chan struct{})
		timeout := make(chan time.Time)
		go test.do(ingresses, services, endpoints, settled, timeout)

		err := awaiter.await(
			&chanWatcher{results: ingresses}, &chanWatcher{results: services}, &chanWatcher{results: endpoints},
			settled, timeout)
		assert.Equal(t, test.expectedError, err, test.description)
	}
}

func Test_Extensions_Ingress_Read(t *testing.T) {
	tests := []struct {
		description       string
		ingressInput      func(namespace, name, targetService string) *unstructured.Unstructured
		ingress           func(namespace, name, targetService string) *unstructured.Unstructured
		endpoint          func(namespace, name string) *unstructured.Unstructured
		service           func(namespace, name string) *unstructured.Unstructured
		expectedSubErrors []string
	}{
		{
			description:  "Read should succeed when Ingress is allocated an IP address and all paths match an existing Endpoint",
			ingressInput: ingressInput,
			ingress:      initializedIngress,
			endpoint:     initializedEndpoint,
		},
		{
			description:  "Read should fail if not all Ingress paths match existing Endpoints",
			ingressInput: ingressInput,
			ingress:      initializedIngress,
			expectedSubErrors: []string{
				"Ingress has at least one rule that does not target any Service. " +
					"Field '.spec.rules[].http.paths[].backend.serviceName' may not match any active Service",
			},
		},
		{
			description:  "Read should succeed when Ingress is allocated an IP address and Service is type ExternalName",
			ingressInput: ingressInput,
			ingress:      initializedIngress,
			service:      externalNameService,
		},
		{
			description:  "Read should fail if Ingress not allocated an IP address",
			ingressInput: ingressInput,
			ingress:      ingressInput,
			endpoint:     initializedEndpoint,
			expectedSubErrors: []string{
				"Ingress .status.loadBalancer field was not updated with a hostname/IP address. " +
					"\n    for more information about this error, see https://pulumi.io/xdv72s",
			},
		},
		{
			description:  "Read should fail if the Ingress does not have an IP address allocated, and not all paths match an existing Endpoint",
			ingressInput: ingressInput,
			ingress:      ingressInput,
			expectedSubErrors: []string{
				"Ingress has at least one rule that does not target any Service. " +
					"Field '.spec.rules[].http.paths[].backend.serviceName' may not match any active Service",
				"Ingress .status.loadBalancer field was not updated with a hostname/IP address. " +
					"\n    for more information about this error, see https://pulumi.io/xdv72s",
			},
		},
	}

	for _, test := range tests {
		awaiter := makeIngressInitAwaiter(mockAwaitConfig(
			test.ingressInput("default", "foo", "foo-4setj4y6")))
		ingress := test.ingress("default", "foo", "foo-4setj4y6")

		var err error
		endpointList := unstructuredList()
		serviceList := unstructuredList()
		if test.endpoint != nil {
			endpoint := test.endpoint("default", "foo-4setj4y6")
			endpointList = unstructuredList(*endpoint)
		}
		if test.service != nil {
			service := test.service("default", "foo-4setj4y6")
			serviceList = unstructuredList(*service)
		}
		err = awaiter.read(ingress, endpointList, serviceList)

		if test.expectedSubErrors != nil {
			assert.Equal(t, test.expectedSubErrors, err.(*initializationError).SubErrors(), test.description)
		} else {
			assert.Nil(t, err, test.description)
		}
	}

}

// --------------------------------------------------------------------------

// Utility constructs.

// --------------------------------------------------------------------------

func ingressInput(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Ingress",
    "metadata": {
        "name": "%s",
        "namespace": "%s"
    },
    "spec": {
        "rules": [
            {
                "http": {
                    "paths": [
                        {
                            "backend": {
                                "serviceName": "%s",
                                "servicePort": 80
                            },
                            "path": "/nginx"
                        }
                    ]
                }
            }
        ]
    },
    "status": {
        "loadBalancer": {}
    }
}`, name, namespace, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

func initializedIngress(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Ingress",
    "metadata": {
        "name": "%s",
        "namespace": "%s"
    },
    "spec": {
        "rules": [
            {
                "http": {
                    "paths": [
                        {
                            "backend": {
                                "serviceName": "%s",
                                "servicePort": 80
                            },
                            "path": "/nginx"
                        }
                    ]
                }
            }
        ]
    },
    "status": {
        "loadBalancer": {
            "ingress": [
                {
                    "hostname": "localhost"
                }
            ]
        }
    }
}`, name, namespace, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

func initializedIngressUnspecifiedPath(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Ingress",
    "metadata": {
        "name": "%s",
        "namespace": "%s"
    },
    "spec": {
        "rules": [
            {
                "http": {
                    "paths": [
                        {
                            "backend": {
                                "serviceName": "%s",
                                "servicePort": 80
                            }
                        }
                    ]
                }
            }
        ]
    },
    "status": {
        "loadBalancer": {
            "ingress": [
                {
                    "hostname": "localhost"
                }
            ]
        }
    }
}`, name, namespace, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

func Test_expectedIngressPath(t *testing.T) {
	type args struct {
		host        string
		path        string
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "host + path", args: args{host: "foo", path: "/bar", serviceName: "baz"}, want: `"foo/bar" -> "baz"`},
		{name: "host only", args: args{host: "foo", serviceName: "baz"}, want: `"foo" -> "baz"`},
		{name: "path only", args: args{path: "/bar", serviceName: "baz"}, want: `"/bar" -> "baz"`},
		{name: "empty", args: args{serviceName: "baz"}, want: `"" (default path) -> "baz"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expectedIngressPath(tt.args.host, tt.args.path, tt.args.serviceName); got != tt.want {
				t.Errorf("expectedIngressPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
