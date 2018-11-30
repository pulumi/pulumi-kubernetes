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
		inputName = "foo"
		targetService = "foo-service"
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

		err := awaiter.await(&chanWatcher{results: statefulsets}, &chanWatcher{results: pods}, timeout, period)
		assert.Equal(t, test.expectedError, err, test.description)
	}
}

// --------------------------------------------------------------------------

// StatefulSet objects.

// --------------------------------------------------------------------------

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
