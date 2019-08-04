package await

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

func Test_Apps_StatefulSet(t *testing.T) {
	const (
		inputNamespace = "default"
		inputName      = "foo"
		targetService  = "foo-service"
	)
	tests := []struct {
		description   string
		do            func(statefulsets, pods chan watch.Event, timeout chan time.Time)
		expectedError error
	}{
		{
			description: "[Revision 1] Should succeed after creating StatefulSet",
			do: func(statefulsets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully creates and initializes StatefulSet object.
				statefulsets <- watchAddedEvent(
					statefulsetAdded(inputNamespace, inputName, targetService))
				statefulsets <- watchAddedEvent(
					statefulsetCreating(inputNamespace, inputName, targetService))
				statefulsets <- watchAddedEvent(
					statefulsetProgressing(inputNamespace, inputName, targetService))
				statefulsets <- watchAddedEvent(
					statefulsetReady(inputNamespace, inputName, targetService))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 2] Should succeed after updating StatefulSet",
			do: func(statefulsets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully updates StatefulSet object.
				statefulsets <- watchAddedEvent(
					statefulsetUpdate(inputNamespace, inputName, targetService))
				statefulsets <- watchAddedEvent(
					statefulsetUpdating(inputNamespace, inputName, targetService))
				statefulsets <- watchAddedEvent(
					statefulsetUpdatingWithActiveReplica(inputNamespace, inputName, targetService))
				statefulsets <- watchAddedEvent(
					statefulsetUpdateSuccess(inputNamespace, inputName, targetService))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 1] Should fail if timeout occurs before successful creation",
			do: func(statefulsets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully creates StatefulSet object.
				statefulsets <- watchAddedEvent(
					statefulsetAdded(inputNamespace, inputName, targetService))
				statefulsets <- watchAddedEvent(
					statefulsetCreating(inputNamespace, inputName, targetService))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: statefulsetCreating(inputNamespace, inputName, targetService),
				subErrors: []string{
					"0 out of 2 replicas succeeded readiness checks",
				}},
		},
		{
			description: "[Revision 2] Should fail if timeout occurs before successful update rollout",
			do: func(statefulsets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully updates StatefulSet object.
				statefulsets <- watchAddedEvent(
					statefulsetUpdate(inputNamespace, inputName, targetService))
				statefulsets <- watchAddedEvent(
					statefulsetUpdating(inputNamespace, inputName, targetService))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: statefulsetUpdating(inputNamespace, inputName, targetService),
				subErrors: []string{
					"0 out of 2 replicas succeeded readiness checks",
					"StatefulSet controller failed to advance from revision \"foo-7b5cf87b78\" to revision \"foo-789c4b994f\"",
				}},
		},
		{
			description: "[Revision 1] Failure should only report Pods from active StatefulSet, part 1",
			do: func(statefulsets, pods chan watch.Event, timeout chan time.Time) {
				podName := "foo-0"

				statefulsets <- watchAddedEvent(
					statefulsetProgressing(inputNamespace, inputName, targetService))

				// Ready Pod should generate no errors.
				pods <- watchAddedEvent(statefulsetReadyPod(inputNamespace, podName, inputName))

				// Pod belonging to some other StatefulSet should not show up in the errors.
				pods <- watchAddedEvent(statefulsetFailedPod(inputNamespace, podName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: statefulsetProgressing(inputNamespace, inputName, targetService),
				subErrors: []string{
					"1 out of 2 replicas succeeded readiness checks",
				}},
		},
		{
			description: "[Revision 2] Failure should only report Pods from active StatefulSet, part 1",
			do: func(statefulsets, pods chan watch.Event, timeout chan time.Time) {
				podName := "foo-0"

				statefulsets <- watchAddedEvent(
					statefulsetUpdating(inputNamespace, inputName, targetService))

				// Ready Pod should generate no errors.
				pods <- watchAddedEvent(statefulsetReadyPod(inputNamespace, podName, inputName))

				// Pod belonging to some other StatefulSet should not show up in the errors.
				pods <- watchAddedEvent(statefulsetFailedPod(inputNamespace, podName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: statefulsetUpdating(inputNamespace, inputName, targetService),
				subErrors: []string{
					"0 out of 2 replicas succeeded readiness checks",
					"StatefulSet controller failed to advance from revision \"foo-7b5cf87b78\" to revision \"foo-789c4b994f\"",
				}},
		},
		{
			description: "[Revision 1] Failure should only report Pods from active StatefulSet, part 2",
			do: func(statefulsets, pods chan watch.Event, timeout chan time.Time) {
				podName := "foo-0"

				statefulsets <- watchAddedEvent(
					statefulsetProgressing(inputNamespace, inputName, targetService))

				// Failed Pod should show up in the errors.
				pods <- watchAddedEvent(statefulsetFailedPod(inputNamespace, podName, inputName))

				// Unrelated successful Pod should generate no errors.
				pods <- watchAddedEvent(statefulsetReadyPod(inputNamespace, podName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: statefulsetProgressing(inputNamespace, inputName, targetService),
				subErrors: []string{
					"1 out of 2 replicas succeeded readiness checks",
					"containers with unready status: [nginx] -- [ErrImagePull] manifest for nginx:busted not found",
				}},
		},
		{
			description: "[Revision 2] Failure should only report Pods from active StatefulSet, part 2",
			do: func(statefulsets, pods chan watch.Event, timeout chan time.Time) {
				podName := "foo-0"

				statefulsets <- watchAddedEvent(
					statefulsetUpdating(inputNamespace, inputName, targetService))

				// Failed Pod should show up in the errors.
				pods <- watchAddedEvent(statefulsetFailedPod(inputNamespace, podName, inputName))

				// Unrelated successful Pod should generate no errors.
				pods <- watchAddedEvent(statefulsetReadyPod(inputNamespace, podName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: statefulsetUpdating(inputNamespace, inputName, targetService),
				subErrors: []string{
					"0 out of 2 replicas succeeded readiness checks",
					"StatefulSet controller failed to advance from revision \"foo-7b5cf87b78\" to revision \"foo-789c4b994f\"",
					"containers with unready status: [nginx] -- [ErrImagePull] manifest for nginx:busted not found",
				}},
		},
	}

	for _, test := range tests {
		awaiter := makeStatefulSetInitAwaiter(
			updateAwaitConfig{
				createAwaitConfig: mockAwaitConfig(statefulsetInput(inputNamespace, inputName, targetService)),
			})
		statefulsets := make(chan watch.Event)
		pods := make(chan watch.Event)

		timeout := make(chan time.Time)
		period := make(chan time.Time)
		go test.do(statefulsets, pods, timeout)

		err := awaiter.await(
			&chanWatcher{results: statefulsets}, &chanWatcher{results: pods}, timeout, period)
		assert.Equal(t, test.expectedError, err, test.description)
	}
}

func Test_Apps_StatefulSet_MultipleUpdates(t *testing.T) {
	tests := []struct {
		description string
		inputs      func() *unstructured.Unstructured
		firstUpdate func(statefulsets, pods chan watch.Event, timeout chan time.Time,
			setLast setLastInputs)
		secondUpdate        func(statefulsets, pods chan watch.Event, timeout chan time.Time)
		firstExpectedError  error
		secondExpectedError error
	}{
		{
			description: "StatefulSet fails, is updated with working config, and then succeeds",
			inputs:      statefulsetFailed,
			firstUpdate: func(
				statefulsets, pods chan watch.Event, timeout chan time.Time,
				setLast setLastInputs,
			) {
				statefulsets <- watchAddedEvent(statefulsetFailed())

				setLast(statefulsetFailed())
				// Timeout. Failed.
				timeout <- time.Now()
			},
			firstExpectedError: &timeoutError{
				object: statefulsetFailed(),
				subErrors: []string{
					"0 out of 2 replicas succeeded readiness checks",
				}},
			secondUpdate: func(statefulset, pods chan watch.Event, timeout chan time.Time) {
				statefulset <- watchAddedEvent(statefulsetUpdatedAfterFailed())
				statefulset <- watchAddedEvent(statefulsetSucceedAfterFailed())

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
	}

	for _, test := range tests {
		awaiter := makeStatefulSetInitAwaiter(
			updateAwaitConfig{
				createAwaitConfig: mockAwaitConfig(test.inputs()),
			})
		statefulsets := make(chan watch.Event)
		pods := make(chan watch.Event)

		timeout := make(chan time.Time)
		period := make(chan time.Time)
		go test.firstUpdate(statefulsets, pods, timeout,
			func(obj *unstructured.Unstructured) {
				awaiter.config.lastInputs = obj
			})

		err := awaiter.await(
			&chanWatcher{results: statefulsets}, &chanWatcher{results: pods}, timeout, period)
		assert.Equal(t, test.firstExpectedError, err, test.description)

		statefulsets = make(chan watch.Event)
		pods = make(chan watch.Event)

		timeout = make(chan time.Time)
		period = make(chan time.Time)
		go test.secondUpdate(statefulsets, pods, timeout)

		err = awaiter.await(
			&chanWatcher{results: statefulsets}, &chanWatcher{results: pods}, timeout, period)
		assert.Equal(t, test.secondExpectedError, err, test.description)
	}
}

func Test_Apps_StatefulSetRead(t *testing.T) {
	const (
		inputNamespace = "default"
		inputName      = "foo"
		targetService  = "foo-service"
	)
	tests := []struct {
		description       string
		statefulset       func(namespace, name, targetService string) *unstructured.Unstructured
		expectedSubErrors []string
	}{
		{
			description: "Read should fail if StatefulSet status empty",
			statefulset: statefulsetAdded,
			expectedSubErrors: []string{
				"0 out of 2 replicas succeeded readiness checks",
				"StatefulSet controller failed to advance from revision \"\" to revision \"\"",
			},
		},
		{
			description: "Read should fail if StatefulSet is progressing",
			statefulset: statefulsetProgressing,
			expectedSubErrors: []string{
				"1 out of 2 replicas succeeded readiness checks",
			},
		},
		{
			description: "Read should succeed if StatefulSet is ready",
			statefulset: statefulsetReady,
		},
		{
			description: "Read should succeed if StatefulSet update is ready",
			statefulset: statefulsetUpdateSuccess,
		},
	}

	for _, test := range tests {
		awaiter := makeStatefulSetInitAwaiter(
			updateAwaitConfig{
				createAwaitConfig: mockAwaitConfig(statefulsetInput(inputNamespace, inputName, targetService)),
			})
		statefulset := test.statefulset(inputNamespace, inputName, targetService)
		err := awaiter.read(statefulset, unstructuredList())
		if test.expectedSubErrors != nil {
			assert.Equal(t, test.expectedSubErrors, err.(*initializationError).SubErrors(), test.description)
		} else {
			assert.Nil(t, err, test.description)
		}
	}
}

// --------------------------------------------------------------------------

// StatefulSet objects.

// --------------------------------------------------------------------------

// statefulsetInput is the user-provided declaration of a StatefulSet
func statefulsetInput(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "%s",
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
                        "image": "nginx",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    }
}`, namespace, name, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetAdded is the initial state of the StatefulSet object after being added through the API
func statefulsetAdded(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
		"generation": 1,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "%s",
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
                        "image": "nginx",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
        "replicas": 0
    }
}`, namespace, name, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetCreating is the state of the StatefulSet object while initial Pods are created but before any are ready
func statefulsetCreating(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
		"generation": 1,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "%s",
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
                        "image": "nginx",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
		"replicas": 1,
		"collisionCount": 0,
		"currentReplicas": 1,
		"currentRevision": "foo-7b5cf87b78",
		"observedGeneration": 1,
		"updateRevision": "foo-7b5cf87b78"
    }
}`, namespace, name, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetProgressing is the state of the StatefulSet object after a Pod is ready
func statefulsetProgressing(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
		"generation": 1,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "%s",
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
                        "image": "nginx",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
		"replicas": 2,
		"collisionCount": 0,
		"currentReplicas": 2,
		"currentRevision": "foo-7b5cf87b78",
		"observedGeneration": 1,
		"updateRevision": "foo-7b5cf87b78",
		"readyReplicas": 1
    }
}`, namespace, name, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetReady is the state of the StatefulSet object after all Pods are ready
func statefulsetReady(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
		"generation": 1,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "%s",
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
                        "image": "nginx",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
		"replicas": 2,
		"collisionCount": 0,
		"currentReplicas": 2,
		"currentRevision": "foo-7b5cf87b78",
		"observedGeneration": 1,
		"updateRevision": "foo-7b5cf87b78",
		"readyReplicas": 2
    }
}`, namespace, name, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetReady is the state of the StatefulSet object after an update is issued through the API
func statefulsetUpdate(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
		"generation": 2,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "%s",
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
                        "image": "nginx:stable",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    }
}`, namespace, name, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetUpdating is the state of the StatefulSet object while an update is rolling out
func statefulsetUpdating(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
		"generation": 2,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "%s",
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
                        "image": "nginx:stable",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
		"replicas": 2,
		"collisionCount": 0,
		"currentReplicas": 1,
		"currentRevision": "foo-7b5cf87b78",
		"observedGeneration": 2,
		"updateRevision": "foo-789c4b994f",
		"readyReplicas": 2
    }
}`, namespace, name, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetUpdating is the state of the StatefulSet object while an update is rolling out and a new Pod is active
func statefulsetUpdatingWithActiveReplica(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
		"generation": 2,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "%s",
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
                        "image": "nginx:stable",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
		"replicas": 2,
		"collisionCount": 0,
		"currentReplicas": 1,
		"currentRevision": "foo-7b5cf87b78",
		"observedGeneration": 2,
		"updateRevision": "foo-789c4b994f",
		"readyReplicas": 1,
		"updatedReplicas": 1
    }
}`, namespace, name, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetUpdateSuccess is the state of the StatefulSet object after an update is rolled out successfully
func statefulsetUpdateSuccess(namespace, name, targetService string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
		"generation": 2,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "%s",
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
                        "image": "nginx:stable",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
		"replicas": 2,
		"collisionCount": 0,
		"currentReplicas": 2,
		"currentRevision": "foo-789c4b994f",
		"observedGeneration": 2,
		"updateRevision": "foo-789c4b994f",
		"readyReplicas": 2
    }
}`, namespace, name, targetService))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetReadyPod is an example Pod created by a StatefulSet that is ready
func statefulsetReadyPod(namespace, name, statefulsetName string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "annotations": {
            "kubernetes.io/limit-ranger": "LimitRanger plugin set: cpu request for container nginx"
        },
        "creationTimestamp": "2018-11-30T21:59:10Z",
        "generateName": "%s-",
        "labels": {
            "app": "foo",
            "controller-revision-hash": "foo-78fd4cddbd",
            "statefulset.kubernetes.io/pod-name": "%s-0"
        },
        "name": "%s",
        "namespace": "%s",
        "ownerReferences": [
            {
                "apiVersion": "apps/v1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "StatefulSet",
                "name": "%s",
                "uid": "984ac0f5-f4ea-11e8-bebe-42010a8a0080"
            }
        ],
        "resourceVersion": "459191",
        "selfLink": "/api/v1/namespaces/default/pods/%s-0",
        "uid": "2a73a5d2-f4eb-11e8-bebe-42010a8a0080"
    },
    "spec": {
        "containers": [
            {
                "image": "nginx",
                "imagePullPolicy": "IfNotPresent",
                "name": "nginx",
                "ports": [
                    {
                        "containerPort": 80,
                        "name": "web",
                        "protocol": "TCP"
                    }
                ],
                "resources": {
                    "requests": {
                        "cpu": "100m"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/usr/share/nginx/html",
                        "name": "www"
                    },
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-p74mp",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "hostname": "foo-0",
        "nodeName": "gke-gke-cluster-8d214cd-default-pool-df2b3fc2-zlkv",
        "priority": 0,
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "subdomain": "ss-service",
        "terminationGracePeriodSeconds": 10,
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
                "name": "www",
                "persistentVolumeClaim": {
                    "claimName": "www-%s-0"
                }
            },
            {
                "name": "default-token-p74mp",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-p74mp"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-11-30T21:59:10Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-11-30T21:59:21Z",
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": null,
                "status": "True",
                "type": "ContainersReady"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-11-30T21:59:10Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://4a89a0e2ab5ad945aad2af0fb12d9660a2715ea77a2e9f214a732a5446088c55",
                "image": "nginx:latest",
                "imageID": "docker-pullable://nginx@sha256:87e9b6904b4286b8d41bba4461c0b736835fcc218f7ecbe5544b53fdd467189f",
                "lastState": {},
                "name": "nginx",
                "ready": true,
                "restartCount": 0,
                "state": {
                    "running": {
                        "startedAt": "2018-11-30T21:59:21Z"
                    }
                }
            }
        ],
        "hostIP": "10.138.0.2",
        "phase": "Running",
        "podIP": "10.32.1.26",
        "qosClass": "Burstable",
        "startTime": "2018-11-30T21:59:10Z"
    }
}`, statefulsetName, statefulsetName, name, namespace, statefulsetName, statefulsetName, statefulsetName))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetFailedPod is an example Pod created by a StatefulSet that is failed
func statefulsetFailedPod(namespace, name, statefulsetName string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "annotations": {
            "kubernetes.io/limit-ranger": "LimitRanger plugin set: cpu request for container nginx"
        },
        "creationTimestamp": "2018-11-30T21:59:10Z",
        "generateName": "%s-",
        "labels": {
            "app": "foo",
            "controller-revision-hash": "foo-78fd4cddbd",
            "statefulset.kubernetes.io/pod-name": "%s-0"
        },
        "name": "%s",
        "namespace": "%s",
        "ownerReferences": [
            {
                "apiVersion": "apps/v1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "StatefulSet",
                "name": "%s",
                "uid": "984ac0f5-f4ea-11e8-bebe-42010a8a0080"
            }
        ],
        "resourceVersion": "459191",
        "selfLink": "/api/v1/namespaces/default/pods/%s-0",
        "uid": "2a73a5d2-f4eb-11e8-bebe-42010a8a0080"
    },
    "spec": {
        "containers": [
            {
                "image": "nginx:busted",
                "imagePullPolicy": "IfNotPresent",
                "name": "nginx",
                "ports": [
                    {
                        "containerPort": 80,
                        "name": "web",
                        "protocol": "TCP"
                    }
                ],
                "resources": {
                    "requests": {
                        "cpu": "100m"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/usr/share/nginx/html",
                        "name": "www"
                    },
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-p74mp",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "hostname": "foo-0",
        "nodeName": "gke-gke-cluster-8d214cd-default-pool-df2b3fc2-zlkv",
        "priority": 0,
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "subdomain": "ss-service",
        "terminationGracePeriodSeconds": 10,
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
                "name": "www",
                "persistentVolumeClaim": {
                    "claimName": "www-%s-0"
                }
            },
            {
                "name": "default-token-p74mp",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-p74mp"
                }
            }
        ]
    },
	"status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-11-30T23:10:58Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-11-30T23:10:58Z",
                "message": "containers with unready status: [nginx]",
                "reason": "ContainersNotReady",
                "status": "False",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": null,
                "message": "containers with unready status: [nginx]",
                "reason": "ContainersNotReady",
                "status": "False",
                "type": "ContainersReady"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-11-30T23:10:58Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "image": "nginx:busted",
                "imageID": "",
                "lastState": {},
                "name": "nginx",
                "ready": false,
                "restartCount": 0,
                "state": {
                    "waiting": {
                        "message": "rpc error: code = Unknown desc = Error response from daemon: manifest for nginx:busted not found",
                        "reason": "ErrImagePull"
                    }
                }
            }
        ],
        "hostIP": "10.138.0.2",
        "phase": "Running",
        "podIP": "10.32.1.26",
        "qosClass": "Burstable",
        "startTime": "2018-11-30T21:59:10Z"
    }
}`, statefulsetName, statefulsetName, name, namespace, statefulsetName, statefulsetName, statefulsetName))
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetFailed is the state of the StatefulSet object that is failing to be ready (invalid image)
func statefulsetFailed() *unstructured.Unstructured {
	obj, err := decodeUnstructured(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "default",
        "name": "foo",
		"generation": 1,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "ss-service",
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
                        "image": "nginx:busted",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
		"replicas": 1,
		"collisionCount": 0,
		"currentReplicas": 1,
		"currentRevision": "foo-7b5cf87b78",
		"observedGeneration": 1,
		"updateRevision": "foo-7b5cf87b78"
    }
}`)
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetFailed is the state of the StatefulSet object that is updating after failing to be ready (invalid image)
func statefulsetUpdatedAfterFailed() *unstructured.Unstructured {
	obj, err := decodeUnstructured(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "default",
        "name": "foo",
		"generation": 2,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "ss-service",
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
                        "image": "nginx:stable",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
		"collisionCount": 0,
		"currentRevision": "foo-7b5cf87b78",
		"observedGeneration": 2,
		"readyReplicas": 1,
		"replicas": 2,
		"updateRevision": "foo-789c4b994f",
		"updatedReplicas": 1
    }
}`)
	if err != nil {
		panic(err)
	}
	return obj
}

// statefulsetSucceedAfterFailed is the state of the StatefulSet object that succeeded
// after failing to be ready (invalid image)
func statefulsetSucceedAfterFailed() *unstructured.Unstructured {
	obj, err := decodeUnstructured(`{
    "kind": "StatefulSet",
    "apiVersion": "apps/v1beta1",
    "metadata": {
        "namespace": "default",
        "name": "foo",
		"generation": 2,
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
		"serviceName": "ss-service",
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
                        "image": "nginx:stable",
						"volumeMounts": [
							{
								"mountPath": "/usr/share/nginx/html",
								"name": "www"
							}
						]
                    }
                ],
				"terminationGracePeriodSeconds": 10
            }
        },
		"volumeClaimTemplates": [
			{
				"metadata": {
					"name": "www"
				},
				"spec": {
					"accessModes": [
						"ReadWriteOnce"
					],
					"resources": {
						"requests": {
							"storage": "1Gi"
						}
					}
				}
			}
		]
    },
	"status": {
		"collisionCount": 0,
		"currentReplicas": 2,
		"currentRevision": "foo-789c4b994f",
		"observedGeneration": 2,
		"readyReplicas": 2,
		"replicas": 2,
		"updateRevision": "foo-789c4b994f"
    }
}`)
	if err != nil {
		panic(err)
	}
	return obj
}
