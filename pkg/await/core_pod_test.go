package await

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

func Test_Core_Pod(t *testing.T) {
	tests := []struct {
		description   string
		do            func(deployments chan watch.Event, timeout chan time.Time)
		expectedError error
	}{
		{
			description: "Should succeed after booting containers",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				// API server successfully initializes Pod.
				pods <- watchAddedEvent(podRunning("default", "foo-4setj4y6"))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "Should transition to success state after initial image pull error",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				// API server successfully adds Pod; pod transitions through various error states
				// until it successfully reaches running state.
				pods <- watchAddedEvent(podAdded("default", "foo-4setj4y6"))
				pods <- watchAddedEvent(podScheduled("default", "foo-4setj4y6"))
				pods <- watchAddedEvent(podContainerCreating("default", "foo-4setj4y6"))
				pods <- watchAddedEvent(podErrImagePull("default", "foo-4setj4y6"))
				pods <- watchAddedEvent(podImagePullBackoff("default", "foo-4setj4y6"))
				pods <- watchAddedEvent(podRunning("default", "foo-4setj4y6"))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "Should fail if Kubelet doesn't boot containers",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				// API server begins to initialize Pod, but does not succeed in time.
				pods <- watchAddedEvent(podAdded("default", "foo-4setj4y6"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: "foo-4setj4y6", subErrors: []string{}},
		},
		{
			description: "Should fail if Pod is scheduled but containers aren't created",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				// API server passes initialized service and endpoint objects back.
				pods <- watchAddedEvent(podScheduled("default", "foo-4setj4y6"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: "foo-4setj4y6", subErrors: []string{}},
		},
		{
			description: "Should fail if Pod is unschedulable",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				// API server fails to schedule Pod, marks status as such.
				pods <- watchAddedEvent(podUnschedulable("default", "foo-4setj4y6"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: "foo-4setj4y6",
				subErrors: []string{
					"Pod unscheduled: [Unschedulable] No nodes are available that match all " +
						"of the predicates: Insufficient cpu (3).",
				},
			},
		},
		{
			description: "Should fail if timeout when creating containers",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				// API server passes initialized service and endpoint objects back.
				pods <- watchAddedEvent(podContainerCreating("default", "foo-4setj4y6"))

				// Mark endpoint objects as having settled. Success.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: "foo-4setj4y6", subErrors: []string{
					"Pod not ready: [ContainersNotReady] containers with unready status: [nginx]",
				}},
		},
		{
			description: "Should fail if timeout after image pull failure",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				// API server tries to initialize Pod, but can't pull the image from $WHATEVER registry.
				pods <- watchAddedEvent(podErrImagePull("default", "foo-4setj4y6"))

				// Mark endpoint objects as having settled. Success.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: "foo-4setj4y6",
				subErrors: []string{
					"Pod not ready: [ContainersNotReady] containers with unready status: [nginx]",
					"[ErrImagePull] repository dsjkdsjkljks not found: does not exist or no " +
						"pull access",
				},
			},
		},
		{
			description: "Should fail if container terminated with error",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				pods <- watchAddedEvent(podTerminatedError("default", "foo-4setj4y6"))

				// Mark endpoint objects as having settled. Success.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: "foo-4setj4y6",
				subErrors: []string{
					"Pod not ready: [ContainersNotReady] containers with unready status: [completer]",
					"[RunContainerError] failed to start container " +
						"\"ce5652ee18060c0f58144968587f5c333cf97dad907c6743e1e95a63541aab72\": Error " +
						"response from daemon: oci runtime error: container_linux.go:262: starting " +
						"container process caused \"exec: \\\"echo foo\\\": executable file not " +
						"found in $PATH\"",
				},
			},
		},
		{
			description: "Should fail if container terminated successfully",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				pods <- watchAddedEvent(podTerminatedSuccess("default", "foo-4setj4y6"))

				// Mark endpoint objects as having settled. Success.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				objectName: "foo-4setj4y6",
				subErrors: []string{
					"Pod not ready: [ContainersNotReady] containers with unready status: [completer]",
					"[Completed] Container completed with exit code 0",
				},
			},
		},
		{
			description: "Should succeed terminated successfully and restart policy is 'Never'",
			do: func(pods chan watch.Event, timeout chan time.Time) {
				pods <- watchAddedEvent(podSucceeded("default", "foo-4setj4y6"))

				// Mark endpoint objects as having settled. Success.
				timeout <- time.Now()
			},
		},
	}

	for _, test := range tests {
		awaiter := makePodInitAwaiter(mockAwaitConfig(podInput("default", "foo-4setj4y6")))
		pods := make(chan watch.Event)

		timeout := make(chan time.Time)
		go test.do(pods, timeout)

		err := awaiter.await(&chanWatcher{results: pods}, timeout)
		assert.Equal(t, test.expectedError, err, test.description)
	}
}

// --------------------------------------------------------------------------

// Utility constructs.

// --------------------------------------------------------------------------

func podInput(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "containers": [
            {"name": "nginx", "image": "nginx:1.15-alpine"}
        ]
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podAdded(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "containers": [
            {"name": "nginx", "image": "nginx:1.15-alpine"}
        ]
    },
    "status": {
        "phase": "Pending",
        "qosClass": "Burstable"
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podUnschedulable(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "containers": [
            {"name": "nginx", "image": "nginx:1.15-alpine"}
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T20:37:58Z",
                "message": "No nodes are available that match all of the predicates: Insufficient cpu (3).",
                "reason": "Unschedulable",
                "status": "False",
                "type": "PodScheduled"
            }
        ],
        "phase": "Pending",
        "qosClass": "Burstable"
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podScheduled(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "containers": [
            {"name": "nginx", "image": "nginx:1.15-alpine"}
        ]
    },
    "status": {
        "phase": "Pending",
        "conditions": [
            {
                "type": "PodScheduled",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z"
            }
        ],
        "qosClass": "Burstable"
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podContainerCreating(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "containers": [
            {"name": "nginx", "image": "nginx:1.15-alpine"}
        ]
    },
    "status": {
        "phase": "Pending",
        "conditions": [
            {
                "type": "Initialized",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z"
            },
            {
                "type": "Ready",
                "status": "False",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z",
                "reason": "ContainersNotReady",
                "message": "containers with unready status: [nginx]"
            },
            {
                "type": "PodScheduled",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z"
            }
        ],
        "hostIP": "10.0.0.7",
        "startTime": "2018-07-30T03:19:45Z",
        "containerStatuses": [
            {
                "name": "nginx",
                "state": {
                    "waiting": {
                        "reason": "ContainerCreating"
                    }
                },
                "lastState": {},
                "ready": false,
                "restartCount": 0,
                "image": "dsjkdsjkljks",
                "imageID": ""
            }
        ],
        "qosClass": "Burstable"
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podErrImagePull(namespace, name string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "containers": [
            {"name": "nginx", "image": "dsjkdsjkljks"}
        ]
    },
    "status": {
        "phase": "Pending",
        "conditions": [
            {
                "type": "Initialized",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z"
            },
            {
                "type": "Ready",
                "status": "False",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z",
                "reason": "ContainersNotReady",
                "message": "containers with unready status: [nginx]"
            },
            {
                "type": "PodScheduled",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z"
            }
        ],
        "hostIP": "10.0.0.7",
        "startTime": "2018-07-30T03:19:45Z",
        "containerStatuses": [
            {
                "name": "nginx",
                "state": {
                    "waiting": {
                        "reason": "ErrImagePull",
                        "message": "rpc error: code = Unknown desc = Error response from daemon: repository dsjkdsjkljks not found: does not exist or no pull access"
                    }
                },
                "lastState": {},
                "ready": false,
                "restartCount": 0,
                "image": "dsjkdsjkljks",
                "imageID": ""
            }
        ],
        "qosClass": "Burstable"
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podImagePullBackoff(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "containers": [
            {"name": "nginx", "image": "dsjkdsjkljks"}
        ]
    },
    "status": {
        "phase": "Pending",
        "conditions": [
            {
                "type": "Initialized",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z"
            },
            {
                "type": "Ready",
                "status": "False",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z",
                "reason": "ContainersNotReady",
                "message": "containers with unready status: [nginx]"
            },
            {
                "type": "PodScheduled",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z"
            }
        ],
        "hostIP": "10.0.0.7",
        "podIP": "10.32.1.13",
        "startTime": "2018-07-30T03:19:45Z",
        "containerStatuses": [
            {
                "name": "nginx",
                "state": {
                    "waiting": {
                        "reason": "ImagePullBackOff",
                        "message": "Back-off pulling image \"dsjkdsjkljks\""
                    }
                },
                "lastState": {},
                "ready": false,
                "restartCount": 0,
                "image": "dsjkdsjkljks",
                "imageID": ""
            }
        ],
        "qosClass": "Burstable"
      }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podRunning(namespace, name string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "containers": [
            {"name": "nginx", "image": "nginx:1.15-alpine"}
        ]
    },
    "status": {
        "phase": "Running",
        "conditions": [
            {
                "type": "Initialized",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z"
            },
            {
                "type": "Ready",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:42:08Z"
            },
            {
                "type": "PodScheduled",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-07-30T03:19:45Z"
            }
        ],
        "hostIP": "10.0.0.7",
        "podIP": "10.32.1.13",
        "startTime": "2018-07-30T03:19:45Z",
        "containerStatuses": [
            {
                "name": "nginx",
                "state": {
                    "running": {
                        "startedAt": "2018-07-30T03:42:07Z"
                    }
                },
                "lastState": {},
                "ready": true,
                "restartCount": 0,
                "image": "nginx:latest",
                "imageID": "docker-pullable://nginx@sha256:4ffd9758ea9ea360fd87d0cee7a2d1cf9dba630bb57ca36b3108dcd3708dc189",
                "containerID": "docker://045ad24362a3c834af31b691e099d30bf7e803c6e79d461d6ae540ac0b21b30f"
            }
        ],
        "qosClass": "Burstable"
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podTerminatedError(namespace, name string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Pod",
    "apiVersion": "v1",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "volumes": [
            {
                "name": "default-token-rkzb2",
                "secret": {
                    "secretName": "default-token-rkzb2",
                    "defaultMode": 420
                }
            }
        ],
        "containers": [
            {
                "name": "completer",
                "image": "alpine",
                "command": [
                    "echo foo"
                ],
                "resources": {},
                "volumeMounts": [
                    {
                        "name": "default-token-rkzb2",
                        "readOnly": true,
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
                    }
                ],
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "imagePullPolicy": "Always"
            }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "minikube",
        "securityContext": {},
        "schedulerName": "default-scheduler"
    },
    "status": {
        "phase": "Running",
        "conditions": [
            {
                "type": "Initialized",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T00:27:03Z"
            },
            {
                "type": "Ready",
                "status": "False",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T00:27:03Z",
                "reason": "ContainersNotReady",
                "message": "containers with unready status: [completer]"
            },
            {
                "type": "PodScheduled",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T00:27:03Z"
            }
        ],
        "hostIP": "192.168.99.100",
        "podIP": "172.17.0.4",
        "startTime": "2018-08-03T00:27:03Z",
        "containerStatuses": [
            {
                "name": "completer",
                "state": {
                    "waiting": {
                        "reason": "RunContainerError",
                        "message": "failed to start container \"ce5652ee18060c0f58144968587f5c333cf97dad907c6743e1e95a63541aab72\": Error response from daemon: oci runtime error: container_linux.go:262: starting container process caused \"exec: \\\"echo foo\\\": executable file not found in $PATH\""
                    }
                },
                "lastState": {
                    "terminated": {
                        "exitCode": 127,
                        "reason": "ContainerCannotRun",
                        "message": "oci runtime error: container_linux.go:262: starting container process caused \"exec: \\\"echo foo\\\": executable file not found in $PATH\"\n",
                        "startedAt": "2018-08-03T00:27:09Z",
                        "finishedAt": "2018-08-03T00:27:09Z",
                        "containerID": "docker://ce5652ee18060c0f58144968587f5c333cf97dad907c6743e1e95a63541aab72"
                    }
                },
                "ready": false,
                "restartCount": 1,
                "image": "alpine:latest",
                "imageID": "docker-pullable://alpine@sha256:7043076348bf5040220df6ad703798fd8593a0918d06d3ce30c6c93be117e430",
                "containerID": "docker://ce5652ee18060c0f58144968587f5c333cf97dad907c6743e1e95a63541aab72"
            }
        ],
        "qosClass": "BestEffort"
    }
}
`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podTerminatedSuccess(namespace, name string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Pod",
    "apiVersion": "v1",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "volumes": [
            {
                "name": "default-token-rkzb2",
                "secret": {
                    "secretName": "default-token-rkzb2",
                    "defaultMode": 420
                }
            }
        ],
        "containers": [
            {
                "name": "completer",
                "image": "alpine",
                "command": [
                    "/bin/sh"
                ],
                "resources": {},
                "volumeMounts": [
                    {
                        "name": "default-token-rkzb2",
                        "readOnly": true,
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
                    }
                ],
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "imagePullPolicy": "Always"
            }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "minikube",
        "securityContext": {},
        "schedulerName": "default-scheduler"
    },
    "status": {
        "phase": "Running",
        "conditions": [
            {
                "type": "Initialized",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T00:53:09Z"
            },
            {
                "type": "Ready",
                "status": "False",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T00:53:09Z",
                "reason": "ContainersNotReady",
                "message": "containers with unready status: [completer]"
            },
            {
                "type": "PodScheduled",
                "status": "True",
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T00:53:09Z"
            }
        ],
        "hostIP": "192.168.99.100",
        "podIP": "172.17.0.4",
        "startTime": "2018-08-03T00:53:09Z",
        "containerStatuses": [
            {
                "name": "completer",
                "state": {
                    "terminated": {
                        "exitCode": 0,
                        "reason": "Completed",
                        "startedAt": "2018-08-03T00:53:16Z",
                        "finishedAt": "2018-08-03T00:53:16Z",
                        "containerID": "docker://7dd9e12bb149cd60eada7d5008ca1c694dd4dcf6b642aa4cdc9c1b2ca607f27b"
                    }
                },
                "lastState": {
                    "terminated": {
                        "exitCode": 0,
                        "reason": "Completed",
                        "startedAt": "2018-08-03T00:53:13Z",
                        "finishedAt": "2018-08-03T00:53:13Z",
                        "containerID": "docker://f56da68168029d2b8d0e7f7915d6aa0a5446879cc0c01924c9464bdca0167d79"
                    }
                },
                "ready": false,
                "restartCount": 1,
                "image": "alpine:latest",
                "imageID": "docker-pullable://alpine@sha256:7043076348bf5040220df6ad703798fd8593a0918d06d3ce30c6c93be117e430",
                "containerID": "docker://7dd9e12bb149cd60eada7d5008ca1c694dd4dcf6b642aa4cdc9c1b2ca607f27b"
            }
        ],
        "qosClass": "BestEffort"
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func podSucceeded(namespace, name string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "%s",
        "name": "%s"
    },
    "spec": {
        "containers": [
            {
                "command": [
                    "/tools/bats/bats",
                    "-t",
                    "/tests/run.sh"
                ],
                "image": "mysql:5.7.14",
                "imagePullPolicy": "IfNotPresent",
                "name": "mysql-test",
                "resources": {
                    "requests": {
                        "cpu": "100m"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/tests",
                        "name": "tests",
                        "readOnly": true
                    },
                    {
                        "mountPath": "/tools",
                        "name": "tools"
                    },
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-bx8t4",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "initContainers": [
            {
                "command": [
                    "bash",
                    "-c",
                    "set -ex\n# copy bats to tools dir\ncp -R /usr/local/libexec/ /tools/bats/\n"
                ],
                "image": "dduportal/bats:0.4.0",
                "imagePullPolicy": "IfNotPresent",
                "name": "test-framework",
                "resources": {
                    "requests": {
                        "cpu": "100m"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/tools",
                        "name": "tools"
                    },
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-bx8t4",
                        "readOnly": true
                    }
                ]
            }
        ],
        "nodeName": "gke-test-ci-default-pool-098d687f-jv41",
        "restartPolicy": "Never",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
            {
                "effect": "NoExecute",
                "key": "node.kubernetes.io/not-ready",
                "operator": "Exists",
                "tolerationSeconds": 300
            },
            {
                "effect": "NoExecute",
                "key": "node.kubernetes.io/unreachable",
                "operator": "Exists",
                "tolerationSeconds": 300
            }
        ],
        "volumes": [
            {
                "configMap": {
                    "defaultMode": 420,
                    "name": "mysql-test"
                },
                "name": "tests"
            },
            {
                "emptyDir": {},
                "name": "tools"
            },
            {
                "name": "default-token-bx8t4",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-bx8t4"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-05T19:27:09Z",
                "reason": "PodCompleted",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-05T19:27:07Z",
                "reason": "PodCompleted",
                "status": "False",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-05T19:27:07Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://3d85d5c4fa4f71505ea6d593e0753f6bf71762e7d317605bad599deb765e2f70",
                "image": "mysql:5.7.14",
                "imageID": "docker-pullable://mysql@sha256:c8f03238ca1783d25af320877f063a36dbfce0daa56a7b4955e6c6e05ab5c70b",
                "lastState": {},
                "name": "mysql-test",
                "ready": false,
                "restartCount": 0,
                "state": {
                    "terminated": {
                        "containerID": "docker://3d85d5c4fa4f71505ea6d593e0753f6bf71762e7d317605bad599deb765e2f70",
                        "exitCode": 0,
                        "finishedAt": "2018-08-05T19:27:09Z",
                        "reason": "Completed",
                        "startedAt": "2018-08-05T19:27:09Z"
                    }
                }
            }
        ],
        "hostIP": "10.0.0.3",
        "initContainerStatuses": [
            {
                "containerID": "docker://ebf0e41f162380240a5f8885a4052e253649f9ecd6a8e2601867c675e0774c76",
                "image": "dduportal/bats:0.4.0",
                "imageID": "docker-pullable://dduportal/bats@sha256:b2d533b27109f7c9ea1e270e23f212c47906346f9cffaa4da6da48ed9d8031da",
                "lastState": {},
                "name": "test-framework",
                "ready": true,
                "restartCount": 0,
                "state": {
                    "terminated": {
                        "containerID": "docker://ebf0e41f162380240a5f8885a4052e253649f9ecd6a8e2601867c675e0774c76",
                        "exitCode": 0,
                        "finishedAt": "2018-08-05T19:27:09Z",
                        "reason": "Completed",
                        "startedAt": "2018-08-05T19:27:09Z"
                    }
                }
            }
        ],
        "phase": "Succeeded",
        "podIP": "10.44.2.67",
        "qosClass": "Burstable",
        "startTime": "2018-08-05T19:27:07Z"
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}
