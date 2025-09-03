// Copyright 2021, Pulumi Corporation.  All rights reserved.

package await

import (
	"fmt"
	"testing"
	"time"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

func Test_Core_Service(t *testing.T) {
	tests := []struct {
		description   string
		serviceInput  func(namespace, name string) *unstructured.Unstructured
		do            func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time)
		version       cluster.ServerVersion
		expectedError error
	}{
		{
			description:  "Should succeed when Service is allocated an IP address and Endpoints target a Pod",
			serviceInput: serviceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes initialized service and endpoint objects back.
				services <- watchAddedEvent(initializedService("default", "foo-4setj4y6"))
				endpoints <- watchAddedEvent(initializedEndpoint("default", "foo-4setj4y6"))

				// Mark endpoint objects as having settled. Success.
				settled <- struct{}{}
			},
		},
		{
			description:  "Should succeed if Endpoints have settled when timeout occurs",
			serviceInput: serviceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes initialized service back.
				services <- watchAddedEvent(initializedService("default", "foo-4setj4y6"))

				// Pass initialized endpoint objects.
				endpoints <- watchAddedEvent(initializedEndpoint("default", "foo-4setj4y6"))

				// Time out. NOTE: the endpoint objects are not marked as settled.
				timeout <- time.Now()
			},
		},
		{
			description:  "Should fail if unrelated Service succeeds",
			serviceInput: serviceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				services <- watchAddedEvent(initializedService("default", "bar"))
				endpoints <- watchAddedEvent(initializedEndpoint("default", "bar"))

				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: serviceInput("default", "foo-4setj4y6"),
				subErrors: []string{
					"Service does not target any Pods. Selected Pods may not be ready, or " +
						"field '.spec.selector' may not match labels on any Pods",
					"Service was not allocated an IP address; does your cloud provider support this?",
				},
			},
		},
		{
			description:  "Should succeed when unrelated Service fails",
			serviceInput: serviceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				services <- watchAddedEvent(initializedService("default", "foo-4setj4y6"))
				endpoints <- watchAddedEvent(initializedEndpoint("default", "foo-4setj4y6"))

				// Unrelated Service should fail because it does not have an Endpoint.
				services <- watchAddedEvent(initializedService("default", "bar"))

				settled <- struct{}{}
				timeout <- time.Now()
			},
		},
		{
			description:  "Should report success immediately even if the next event is a failure",
			serviceInput: serviceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes initialized service and endpoint objects back.
				services <- watchAddedEvent(initializedService("default", "foo-4setj4y6"))
				endpoints <- watchAddedEvent(initializedEndpoint("default", "foo-4setj4y6"))

				// Mark endpoint objects as having settled. Success.
				settled <- struct{}{}

				endpoints <- watchAddedEvent(
					uninitializedEndpoint("default", "foo-4setj4y6"))
			},
		},
		{
			description:  "Should fail if neither the Service nor the Endpoints have initialized",
			serviceInput: serviceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// Trigger timeout.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: serviceInput("default", "foo-4setj4y6"),
				subErrors: []string{
					"Service does not target any Pods. Selected Pods may not be ready, or " +
						"field '.spec.selector' may not match labels on any Pods",
					"Service was not allocated an IP address; does your cloud provider support this?",
				},
			},
		},
		{
			description:  "Should succeed with no endpoints",
			serviceInput: serviceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes initialized service back.
				services <- watchAddedEvent(initializedService("default", "foo-4setj4y6"))

				// Pass uninitialized endpoint objects. Mark them as settled.
				endpoints <- watchAddedEvent(
					emptyEndpoint("default", "foo-4setj4y6"))
				settled <- struct{}{}
			},
		},
		{
			description:  "Should fail if Endpoints have not initialized",
			serviceInput: serviceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes initialized service back.
				services <- watchAddedEvent(initializedService("default", "foo-4setj4y6"))

				// Pass uninitialized endpoint objects. Mark them as settled.
				endpoints <- watchAddedEvent(
					uninitializedEndpoint("default", "foo-4setj4y6"))
				settled <- struct{}{}

				// Finally, time out.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: initializedService("default", "foo-4setj4y6"),
				subErrors: []string{
					"Service does not target any Pods. Selected Pods may not be ready, or " +
						"field '.spec.selector' may not match labels on any Pods",
				},
			},
		},
		{
			description:  "Should fail if Service is not allocated an IP address",
			serviceInput: serviceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				// API server passes uninitialized service back.
				services <- watchAddedEvent(serviceInput("default", "foo-4setj4y6"))

				// Pass initialized endpoint objects. Mark them as settled.
				endpoints <- watchAddedEvent(
					initializedEndpoint("default", "foo-4setj4y6"))
				settled <- struct{}{}

				// Finally, time out.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: serviceInput("default", "foo-4setj4y6"),
				subErrors: []string{
					"Service was not allocated an IP address; does your cloud provider support this?",
				},
			},
		},

		{
			description:  "Should succeed if Externalname doesn't target any Pods",
			serviceInput: externalNameService,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				services <- watchAddedEvent(externalNameService("default", "foo-4setj4y6"))

				// Finally, time out.
				timeout <- time.Now()
			},
		},
		{
			description:  "[Part 1] Should succeed if empty headless service doesn't target any Pods",
			serviceInput: headlessEmptyServiceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				services <- watchAddedEvent(headlessEmptyServiceOutput1("default", "foo-4setj4y6"))

				// Finally, time out.
				timeout <- time.Now()
			},
		},
		{
			description:  "[Part 2] Should succeed if empty headless service doesn't target any Pods",
			serviceInput: headlessEmptyServiceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				services <- watchAddedEvent(headlessEmptyServiceOutput2("default", "foo-4setj4y6"))

				// Finally, time out.
				timeout <- time.Now()
			},
		},
		{
			description:  "Should succeed if non-empty headless service targets Pods",
			serviceInput: headlessNonemptyServiceInput,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				services <- watchAddedEvent(headlessNonemptyServiceOutput("default", "foo-4setj4y6"))
				endpoints <- watchAddedEvent(initializedEndpoint("default", "foo-4setj4y6"))

				// Finally, time out.
				timeout <- time.Now()
			},
		},
		{
			description:  "Should succeed if non-empty headless service doesn't target any Pods before k8s 1.12",
			serviceInput: headlessNonemptyServiceInput,
			version:      cluster.ServerVersion{Major: 1, Minor: 11},
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				services <- watchAddedEvent(headlessNonemptyServiceOutput("default", "foo-4setj4y6"))

				// Finally, time out.
				timeout <- time.Now()
			},
		},
		{
			description:  "Should succeed if non-empty headless service doesn't target any Pods",
			serviceInput: headlessNonemptyServiceInput,
			version:      cluster.ServerVersion{Major: 1, Minor: 12},
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				services <- watchAddedEvent(headlessNonemptyServiceOutput("default", "foo-4setj4y6"))

				// Finally, time out.
				timeout <- time.Now()
			},
		},
		{
			description:  "Should succeed if no selector",
			serviceInput: serviceWithoutSelector,
			do: func(services, endpoints chan watch.Event, settled chan struct{}, timeout chan time.Time) {
				services <- watchAddedEvent(serviceWithoutSelector("default", "foo-4setj4y6"))
			},
		},
	}

	for _, test := range tests {
		awaiter := makeServiceInitAwaiter(
			mockAwaitConfig(test.serviceInput("default", "foo-4setj4y6")))

		services := make(chan watch.Event)
		endpoints := make(chan watch.Event)
		settled := make(chan struct{})
		timeout := make(chan time.Time)
		go test.do(services, endpoints, settled, timeout)

		obj, err := awaiter.await(services, endpoints, timeout, settled, test.version)
		assert.NotNil(t, obj)
		assert.Equal(t, test.expectedError, err, test.description)
	}
}

func Test_Core_Service_Read(t *testing.T) {
	tests := []struct {
		description       string
		serviceInput      func(namespace, name string) *unstructured.Unstructured
		service           func(namespace, name string) *unstructured.Unstructured
		endpoint          func(namespace, name string) *unstructured.Unstructured
		version           cluster.ServerVersion
		expectedSubErrors []string
	}{
		{
			description:  "Read should fail if Service does not target any Pods",
			serviceInput: serviceInput,
			service:      initializedService,
			endpoint:     uninitializedEndpoint,
			expectedSubErrors: []string{
				"Service does not target any Pods. Selected Pods may not be ready, or " +
					"field '.spec.selector' may not match labels on any Pods",
			},
		},
		{
			description:  "Read should succeed if Service does target Pods",
			serviceInput: serviceInput,
			service:      initializedService,
			endpoint:     initializedEndpoint,
		},
		{
			description:  "Read should fail if Service not allocated an IP address",
			serviceInput: serviceInput,
			service:      serviceInput,
			endpoint:     initializedEndpoint,
			expectedSubErrors: []string{
				"Service was not allocated an IP address; does your cloud provider support this?",
			},
		},
		{
			description:  "Read succeed if ExternalName doesn't target any Pods",
			serviceInput: externalNameService,
			service:      externalNameService,
		},
		{
			description:  "Read succeed if headless empty Service doesn't target any Pods",
			serviceInput: headlessEmptyServiceInput,
			service:      headlessEmptyServiceInput,
		},
		{
			description:  "Read succeed if headless non-empty Service doesn't target any Pods before k8s 1.12",
			serviceInput: headlessNonemptyServiceInput,
			service:      headlessNonemptyServiceInput,
			version:      cluster.ServerVersion{Major: 1, Minor: 11},
		},
		{
			description:  "Read succeed if headless non-empty Service doesn't target any Pods",
			serviceInput: headlessNonemptyServiceInput,
			service:      headlessNonemptyServiceInput,
			version:      cluster.ServerVersion{Major: 1, Minor: 12},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			awaiter := makeServiceInitAwaiter(mockAwaitConfig(test.serviceInput("default", "foo-4setj4y6")))
			service := test.service("default", "foo-4setj4y6")

			var err error
			if test.endpoint != nil {
				endpoint := test.endpoint("default", "foo-4setj4y6")
				err = awaiter.read(service, unstructuredList(*endpoint), test.version)
			} else {
				err = awaiter.read(service, unstructuredList(), test.version)
			}
			if test.expectedSubErrors != nil {
				assert.Equal(t, test.expectedSubErrors, err.(*initializationError).SubErrors(), test.description)
			} else {
				assert.Nil(t, err, test.description)
			}
		})
	}
}

// --------------------------------------------------------------------------

// Utility constructs.

// --------------------------------------------------------------------------

func serviceInput(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "labels": {
            "app": "%s"
        },
        "name": "%s",
        "namespace": "%s"
    },
    "spec": {
        "clusterIP": "10.35.241.240",
        "externalTrafficPolicy": "Cluster",
        "ports": [
            {
                "nodePort": 32277,
                "port": 6379,
                "protocol": "TCP",
                "targetPort": 6379
            }
        ],
        "selector": {
            "app": "foo"
        },
        "sessionAffinity": "None",
        "type": "LoadBalancer"
    },
    "status": {
        "loadBalancer": {}
    }
}`, name, name, namespace))
	if err != nil {
		panic(err)
	}
	return obj
}

func initializedService(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "labels": {
            "app": "%s"
        },
        "name": "%s",
        "namespace": "%s"
    },
    "spec": {
        "clusterIP": "10.35.241.240",
        "externalTrafficPolicy": "Cluster",
        "ports": [
            {
                "nodePort": 32277,
                "port": 6379,
                "protocol": "TCP",
                "targetPort": 6379
            }
        ],
        "selector": {
            "app": "foo"
        },
        "sessionAffinity": "None",
        "type": "LoadBalancer"
    },
    "status": {
        "loadBalancer": {
            "ingress": [
                {
                    "ip": "35.184.65.22"
                }
            ]
        }
    }
}`, name, name, namespace))
	if err != nil {
		panic(err)
	}
	return obj
}

func serviceWithoutSelector(namespace, name string) *unstructured.Unstructured {
	obj, _ := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "name": "%s",
        "namespace": "%s"
    },
    "spec": {
        "ports": [
            {
                "port": 6379
            }
        ]
      }
	}`, name, namespace))
	return obj
}

func emptyEndpoint(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(
		fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Endpoints",
    "metadata": {
        "labels": {
            "app": "%s"
        },
        "name": "%s",
        "namespace": "%s"
    }
}`, name, name, namespace))
	if err != nil {
		panic(err)
	}
	return obj
}

func uninitializedEndpoint(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(
		fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Endpoints",
    "metadata": {
        "labels": {
            "app": "%s"
        },
        "name": "%s",
        "namespace": "%s"
    },
    "subsets": [
        {
            "notReadyAddresses": [
                {
                    "ip": "35.192.99.34"
                }
            ],
            "ports": [
                {
                    "name": "https",
                    "port": 443,
                    "protocol": "TCP"
                }
            ]
        }
    ]
}`, name, name, namespace))
	if err != nil {
		panic(err)
	}
	return obj
}

func initializedEndpoint(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(
		fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Endpoints",
    "metadata": {
        "labels": {
            "app": "%s"
        },
        "name": "%s",
        "namespace": "%s"
    },
    "subsets": [
        {
            "addresses": [
                {
                    "ip": "35.192.99.34"
                }
            ],
            "ports": [
                {
                    "name": "https",
                    "port": 443,
                    "protocol": "TCP"
                }
            ]
        }
    ]
}`, name, name, namespace))
	if err != nil {
		panic(err)
	}
	return obj
}

func externalNameService(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(
		fmt.Sprintf(`{
			"apiVersion": "v1",
			"kind": "Service",
			"metadata": {
				"name": "%s",
				"namespace": "%s"
			},
			"spec": {
				"externalName": "google.com",
				"sessionAffinity": "None",
				"type": "ExternalName"
			},
			"status": {
				"loadBalancer": {}
			}
		}`, name, namespace))
	if err != nil {
		panic(err)
	}
	return obj
}

func headlessEmptyServiceInput(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
  "kind": "Service",
  "apiVersion": "v1",
  "metadata": {
    "namespace": "%s",
    "name": "%s"
  },
  "spec": {
    "clusterIP": "None",
    "selector": {},
    "ports": [
      {
        "protocol": "TCP",
        "port": 80,
        "targetPort": 9376
      }
    ]
  }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func headlessEmptyServiceOutput1(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "clusterIP": "None",
        "ports": [
            {
                "port": 80,
                "protocol": "TCP",
                "targetPort": 9376
            }
        ],
        "sessionAffinity": "None",
        "type": "ClusterIP"
    },
    "status": {
        "loadBalancer": {}
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func headlessEmptyServiceOutput2(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "clusterIP": "None",
        "selector": {},
        "ports": [
            {
                "port": 80,
                "protocol": "TCP",
                "targetPort": 9376
            }
        ],
        "sessionAffinity": "None",
        "type": "ClusterIP"
    },
    "status": {
        "loadBalancer": {}
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func headlessNonemptyServiceInput(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
  "kind": "Service",
  "apiVersion": "v1",
  "metadata": {
    "namespace": "%s",
    "name": "%s"
  },
  "spec": {
    "clusterIP": "None",
    "selector": {
      "app": "foo"
    },
    "ports": [
      {
        "protocol": "TCP",
        "port": 80,
        "targetPort": 9376
      }
    ]
  }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func headlessNonemptyServiceOutput(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "clusterIP": "None",
        "ports": [
            {
                "port": 80,
                "protocol": "TCP",
                "targetPort": 9376
            }
        ],
        "selector": {
            "app": "foo"
        },
        "sessionAffinity": "None",
        "type": "ClusterIP"
    },
    "status": {
        "loadBalancer": {}
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func unstructuredList(us ...unstructured.Unstructured) *unstructured.UnstructuredList {
	return &unstructured.UnstructuredList{Items: us}
}
