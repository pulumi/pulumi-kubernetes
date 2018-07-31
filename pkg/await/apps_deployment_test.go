package await

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	inputNamespace          = "default"
	deploymentInputName     = "foo-4setj4y6"
	replicaSetGeneratedName = "foo-4setj4y6-7cdf7ddc54"
	revision1               = "1"
)

func Test_Apps_Deployment(t *testing.T) {
	tests := []struct {
		description   string
		do            func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time)
		expectedError error
	}{
		{
			description: "Should succeed after creating first ReplicaSet",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully creates and initializes Deployment and ReplicaSet
				// objects.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressingUnavailable(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))

				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "Should succeed even if ReplicaSet becomes available before Deployment repots it",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully creates and initializes Deployment and ReplicaSet
				// objects.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressingUnavailable(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "Should succeed if update has rolled out",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// Deployment is updated by the user. The controller creates and successfully
				// initializes the ReplicaSet.
				deployments <- watchAddedEvent(
					deploymentUpdated(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentUpdatedReplicaSetProgressing(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentUpdatedReplicaSetProgressed(inputNamespace, deploymentInputName, revision1))

				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "Should fail if unrelated Deployment succeeds",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				deployments <- watchAddedEvent(deploymentRolloutComplete(inputNamespace, "bar", revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, "bar-ablksd", "bar", revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{objectName: deploymentInputName, subErrors: []string{}},
		},
		{
			description: "Should succeed when unrelated deployment fails",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				deployments <- watchAddedEvent(deploymentAdded(inputNamespace, "bar", revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, "bar-ablksd", "bar", revision1))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "Should report success immediately even if the next event is a failure",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressingInvalidContainer(inputNamespace, deploymentInputName, revision1))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "Should fail if timeout occurs before Deployment controller progresses",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller hasn't created the ReplicaSet when we time
				// out.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: deploymentInputName, subErrors: []string{}},
		},
		{
			description: "Should fail if timeout occurs before ReplicaSet is created",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller creates ReplicaSet, but the replication
				// controller does not start initializing it before it errors out.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: deploymentInputName, subErrors: []string{"Updated ReplicaSet was never created"}},
		},
		{
			description: "Should fail if timeout occurs before ReplicaSet becomes available",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller creates ReplicaSet, but it's still
				// unavailable when we time out.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressingUnavailable(inputNamespace, deploymentInputName, revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: deploymentInputName,
				subErrors: []string{
					"[MinimumReplicasUnavailable] Deployment does not have minimum availability.",
					"Updated ReplicaSet was never created"}},
		},
		{
			description: "Should fail if new ReplicaSet isn't created after an update",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// Deployment is updated by the user. The controller does not create a new
				// ReplicaSet before we time out.
				deployments <- watchAddedEvent(
					deploymentUpdated(inputNamespace, deploymentInputName, revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: deploymentInputName, subErrors: []string{"Updated ReplicaSet was never created"}},
		},
		{
			description: "Should fail if timeout before new ReplicaSet becomes available",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// Deployment is updated by the user. The controller creates the ReplicaSet, but we
				// time out before it can complete initializing.
				deployments <- watchAddedEvent(
					deploymentUpdated(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentUpdatedReplicaSetProgressing(inputNamespace, deploymentInputName, revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: deploymentInputName, subErrors: []string{"Updated ReplicaSet was never created"}},
		},
		{
			description: "Should fail if Deployment not progressing",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller creates ReplicaSet, and it tries to
				// progress, but it fails.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentNotProgressing(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: deploymentInputName,
				subErrors: []string{
					`[ProgressDeadlineExceeded] ReplicaSet "foo-13y9rdnu-b94df86d6" has timed ` +
						`out progressing.`}},
		},
		{
			description: "Should fail if Deployment is progressing but new ReplicaSet not available",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller creates ReplicaSet, and it tries to
				// progress, but it will not, because it is using an invalid container image.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressingInvalidContainer(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: deploymentInputName, subErrors: []string{}},
		},
		{
			description: "Failure should only report Pods from active ReplicaSet, part 1",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				readyPodName := "foo-4setj4y6-7cdf7ddc54-kvh2w"

				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Ready Pod should generate no errors.
				pods <- watchAddedEvent(deployedReadyPod(inputNamespace, readyPodName, replicaSetGeneratedName))

				// Pod belonging to some other ReplicaSet should not show up in the errors.
				pods <- watchAddedEvent(deployedFailedPod(inputNamespace, readyPodName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{objectName: deploymentInputName, subErrors: []string{}},
		},
		{
			description: "Failure should only report Pods from active ReplicaSet, part 2",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				readyPodName := "foo-4setj4y6-7cdf7ddc54-kvh2w"

				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Failed Pod should show up in the errors.
				pods <- watchAddedEvent(deployedFailedPod(inputNamespace, readyPodName, replicaSetGeneratedName))

				// // Unrelated successful Pod should generate no errors.
				pods <- watchAddedEvent(deployedReadyPod(inputNamespace, readyPodName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{objectName: deploymentInputName, subErrors: []string{}},
		},
		{
			description: "Should fail if ReplicaSet generations do not match",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, "2"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{objectName: deploymentInputName, subErrors: []string{
				"Updated ReplicaSet was never created"}},
		},
	}

	for _, test := range tests {
		awaiter := makeDeploymentInitAwaiter(
			updateAwaitConfig{
				createAwaitConfig: mockAwaitConfig(deploymentInput(inputNamespace, deploymentInputName)),
			})
		deployments := make(chan watch.Event)
		replicaSets := make(chan watch.Event)
		pods := make(chan watch.Event)

		timeout := make(chan time.Time)
		period := make(chan time.Time)
		go test.do(deployments, replicaSets, pods, timeout)

		err := awaiter.await(&mockWatcher{results: deployments}, &mockWatcher{results: replicaSets},
			&mockWatcher{results: pods}, timeout, period)
		assert.Equal(t, test.expectedError, err, test.description)
	}
}

// --------------------------------------------------------------------------

// Deployment objects.

// --------------------------------------------------------------------------

func deploymentInput(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentAdded(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {}
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentProgressing(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "conditions": [
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:04Z",
                "lastTransitionTime": "2018-07-31T21:49:04Z",
                "reason": "NewReplicaSetCreated",
                "message": "Created new replica set \"foo-lobqxn87-546cb87d96\""
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentNotProgressing(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "generation": 3,
        "labels": {
            "app": "foo"
        },
        "namespace": "%s",
        "name": "%s",
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "sdkjsdjkljklds",
                        "name": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 1,
        "conditions": [
            {
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            },
            {
                "lastTransitionTime": "2018-08-01T02:46:31Z",
                "lastUpdateTime": "2018-08-01T02:46:31Z",
                "message": "ReplicaSet \"foo-13y9rdnu-b94df86d6\" has timed out progressing.",
                "reason": "ProgressDeadlineExceeded",
                "status": "False",
                "type": "Progressing"
            }
        ],
        "observedGeneration": 3,
        "readyReplicas": 1,
        "replicas": 2,
        "unavailableReplicas": 1,
        "updatedReplicas": 1
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentProgressingInvalidContainer(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "generation": 4,
        "labels": {
            "app": "foo"
        },
        "namespace": "%s",
        "name": "%s",
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "sdkjlsdlkj",
                        "imagePullPolicy": "Always",
                        "name": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 1,
        "conditions": [
            {
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            },
            {
                "lastTransitionTime": "2018-08-01T03:04:50Z",
                "lastUpdateTime": "2018-08-01T03:04:50Z",
                "message": "ReplicaSet \"foo-13y9rdnu-58ddf8f46\" is progressing.",
                "reason": "ReplicaSetUpdated",
                "status": "True",
                "type": "Progressing"
            }
        ],
        "observedGeneration": 4,
        "readyReplicas": 1,
        "replicas": 2,
        "unavailableReplicas": 1,
        "updatedReplicas": 1
    }
}
`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentProgressingUnavailable(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 1,
        "replicas": 1,
        "updatedReplicas": 1,
        "unavailableReplicas": 1,
        "conditions": [
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:04Z",
                "lastTransitionTime": "2018-07-31T21:49:04Z",
                "reason": "NewReplicaSetCreated",
                "message": "Created new replica set \"foo-lobqxn87-546cb87d96\""
            },
            {
                "type": "Available",
                "status": "False",
                "lastUpdateTime": "2018-07-31T21:49:04Z",
                "lastTransitionTime": "2018-07-31T21:49:04Z",
                "reason": "MinimumReplicasUnavailable",
                "message": "Deployment does not have minimum availability."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentRolloutComplete(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 1,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:11Z",
                "lastTransitionTime": "2018-07-31T21:49:11Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            },
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:11Z",
                "lastTransitionTime": "2018-07-31T21:49:04Z",
                "reason": "NewReplicaSetAvailable",
                "message": "ReplicaSet \"foo-lobqxn87-546cb87d96\" has successfully progressed."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentUpdated(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 2,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx:1.15-alpine"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 1,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            },
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "lastTransitionTime": "2018-07-31T23:42:19Z",
                "reason": "NewReplicaSetAvailable",
                "message": "ReplicaSet \"foo-13y9rdnu-546cb87d96\" has successfully progressed."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentUpdatedReplicaSetProgressing(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 2,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx:1.15-alpine"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 2,
        "replicas": 2,
        "updatedReplicas": 1,
        "readyReplicas": 2,
        "availableReplicas": 2,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            },
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:43:18Z",
                "lastTransitionTime": "2018-07-31T23:42:19Z",
                "reason": "ReplicaSetUpdated",
                "message": "ReplicaSet \"foo-13y9rdnu-5694b49bf5\" is progressing."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentUpdatedReplicaSetProgressed(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 2,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx:1.15-alpine"
                    }
                ],
                "restartPolicy": "Always",
                "terminationGracePeriodSeconds": 30,
                "dnsPolicy": "ClusterFirst",
                "securityContext": {},
                "schedulerName": "default-scheduler"
            }
        }
    },
    "status": {
        "observedGeneration": 2,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            },
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:43:18Z",
                "lastTransitionTime": "2018-07-31T23:42:19Z",
                "reason": "NewReplicaSetAvailable",
                "message": "ReplicaSet \"foo-13y9rdnu-5694b49bf5\" has successfully progressed."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

// --------------------------------------------------------------------------

// ReplicaSet objects.

// --------------------------------------------------------------------------

func availableReplicaSet(namespace, name, deploymentName, revision string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "ReplicaSet",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/desired-replicas": "3",
            "deployment.kubernetes.io/max-replicas": "4",
            "deployment.kubernetes.io/revision": "%s",
            "deployment.kubernetes.io/revision-history": "3",
            "moolumi.com/metricsChecked": "true",
            "pulumi.com/autonamed": "true"
        },
        "creationTimestamp": "2018-08-03T05:03:53Z",
        "generation": 1,
        "labels": {
            "app": "foo",
            "pod-template-hash": "3789388710"
        },
        "namespace": "%s",
        "name": "%s",
        "ownerReferences": [
            {
                "apiVersion": "extensions/v1beta1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "Deployment",
                "name": "%s",
                "uid": "e4a728af-96d9-11e8-9050-080027bd9056"
            }
        ]
    },
    "spec": {
        "replicas": 3,
        "selector": {
            "matchLabels": {
                "app": "foo",
                "pod-template-hash": "3789388710"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo",
                    "pod-template-hash": "3789388710"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "nginx:1.15-alpine",
                        "imagePullPolicy": "Always",
                        "name": "nginx",
                        "resources": {},
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File",
                        "volumeMounts": [
                            {
                                "mountPath": "/etc/config",
                                "name": "config-volume"
                            }
                        ]
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30,
                "volumes": [
                    {
                        "configMap": {
                            "defaultMode": 420,
                            "name": "configmap-rollout-mfonkaf3"
                        },
                        "name": "config-volume"
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 3,
        "fullyLabeledReplicas": 3,
        "observedGeneration": 3,
        "readyReplicas": 3,
        "replicas": 3
    }
}
`, revision, namespace, name, deploymentName))
	if err != nil {
		panic(err)
	}
	return obj
}

// --------------------------------------------------------------------------

// Pod objects.

// --------------------------------------------------------------------------

func deployedReadyPod(namespace, name, replicaSetName string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "annotations": {
            "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"default\",\"name\":\"%s\",\"uid\":\"9e300c56-96da-11e8-9050-080027bd9056\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"813941\"}}\n"
        },
        "creationTimestamp": "2018-08-03T05:04:10Z",
        "generateName": "%s-",
        "labels": {
            "app": "foo",
            "pod-template-hash": "3789388710"
        },
        "namespace": "%s",
        "name": "%s",
        "ownerReferences": [
            {
                "apiVersion": "extensions/v1beta1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "ReplicaSet",
                "name": "%s",
                "uid": "9e300c56-96da-11e8-9050-080027bd9056"
            }
        ]
    },
    "spec": {
        "containers": [
            {
                "image": "nginx:1.15-alpine",
                "imagePullPolicy": "Always",
                "name": "nginx",
                "resources": {},
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/etc/config",
                        "name": "config-volume"
                    },
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-rkzb2",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "nodeName": "minikube",
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "volumes": [
            {
                "configMap": {
                    "defaultMode": 420,
                    "name": "configmap-rollout-mfonkaf3"
                },
                "name": "config-volume"
            },
            {
                "name": "default-token-rkzb2",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-rkzb2"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T05:04:10Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T05:04:13Z",
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T05:04:10Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://a91bc460f583402484ceeef5801a0f6221bb71f184359e79a8e795e7f463ba02",
                "image": "nginx:1.15-alpine",
                "imageID": "docker-pullable://nginx@sha256:23e4dacbc60479fa7f23b3b8e18aad41bd8445706d0538b25ba1d575a6e2410b",
                "lastState": {},
                "name": "nginx",
                "ready": true,
                "restartCount": 0,
                "state": {
                    "running": {
                        "startedAt": "2018-08-03T05:04:13Z"
                    }
                }
            }
        ],
        "hostIP": "192.168.99.100",
        "phase": "Running",
        "podIP": "172.17.0.5",
        "qosClass": "BestEffort",
        "startTime": "2018-08-03T05:04:10Z"
    }
}

`, replicaSetName, replicaSetName, namespace, name, replicaSetName))
	if err != nil {
		panic(err)
	}
	return obj
}

func deployedFailedPod(namespace, name, replicaSetName string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "annotations": {
            "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"default\",\"name\":\"%s\",\"uid\":\"c80dda50-96e4-11e8-9050-080027bd9056\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"819008\"}}\n"
        },
        "generateName": "%s-",
        "labels": {
            "app": "foo",
            "pod-template-hash": "3789350985"
        },
        "namespace": "%s",
        "name": "%s",
        "ownerReferences": [
            {
                "apiVersion": "extensions/v1beta1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "ReplicaSet",
                "name": "%s",
                "uid": "c80dda50-96e4-11e8-9050-080027bd9056"
            }
        ]
    },
    "spec": {
        "containers": [
            {
                "image": "sdkjlsdlkj",
                "imagePullPolicy": "Always",
                "name": "nginx",
                "resources": {},
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/etc/config",
                        "name": "config-volume"
                    },
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-rkzb2",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "nodeName": "minikube",
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "volumes": [
            {
                "configMap": {
                    "defaultMode": 420,
                    "name": "configmap-rollout-mfonkaf3"
                },
                "name": "config-volume"
            },
            {
                "name": "default-token-rkzb2",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-rkzb2"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T06:16:38Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T06:16:38Z",
                "message": "containers with unready status: [nginx]",
                "reason": "ContainersNotReady",
                "status": "False",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T06:16:38Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "image": "sdkjlsdlkj",
                "imageID": "",
                "lastState": {},
                "name": "nginx",
                "ready": false,
                "restartCount": 0,
                "state": {
                    "waiting": {
                        "message": "Back-off pulling image \"sdkjlsdlkj\"",
                        "reason": "ImagePullBackOff"
                    }
                }
            }
        ],
        "hostIP": "192.168.99.100",
        "phase": "Pending",
        "podIP": "172.17.0.7",
        "qosClass": "BestEffort",
        "startTime": "2018-08-03T06:16:38Z"
    }
}`, replicaSetName, replicaSetName, namespace, name, replicaSetName))
	if err != nil {
		panic(err)
	}
	return obj
}
