// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v2

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gk "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
	gs "github.com/onsi/gomega/gstruct"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/internals"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	pgm "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gomega"
)

const (
	manifest = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: crontabs.stable.example.com
  # Purposefully set annotations here to be null: pulumi/pulumi-kubernetes#3585
  annotations: 
spec:
  group: stable.example.com
  versions:
  - name: v1
  served: true
  storage: true
  schema:
    openAPIV3Schema:
      type: object
      properties:
        spec:
          type: object
          properties:
            cronSpec:
              type: string
            image:
              type: string
            replicas:
              type: integer
  scope: Namespaced
  names:
    plural: crontabs
    singular: crontab
    kind: CronTab
---
apiVersion: v1
kind: Namespace
metadata:
  name: my-namespace
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-map
  namespace: my-namespace
data:
  altGreeting: "Good Morning!"
---
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: my-new-cron-object
  namespace: my-namespace
spec:
  cronSpec: "* * * * */5"
  image: my-awesome-cron-image
`

	list = `
apiVersion: v1
kind: List
items:
- apiVersion: v1
  kind: ConfigMap
  metadata:
    name: map-1
- apiVersion: v1
  kind: ConfigMap
  metadata:
    name: map-2
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: map-3
`
)

var _ = gk.Describe("Register", func() {
	var tc *componentProviderTestContext
	var registerOpts *RegisterOptions

	gk.BeforeEach(func() {
		tc = newTestContext(gk.GinkgoTB())
		registerOpts = &RegisterOptions{}
	})

	register := func(ctx context.Context) (resources pulumi.ArrayOutput, err error) {
		// use RunWithContext to reliably wait for outstanding RPCs to complete
		err = pulumi.RunWithContext(tc.NewContext(ctx), func(ctx *pulumi.Context) error {
			resources, err = Register(ctx, *registerOpts)
			return err
		})
		return resources, err
	}

	gk.Context("given the objects in the manifest", func() {
		gk.BeforeEach(func() {
			resources, err := yamlDecode(manifest)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			registerOpts.Objects = resources
		})

		gk.It("should register some resources", func(ctx context.Context) {
			_, err := register(ctx)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())

			gm.Expect(tc.monitor.Resources()).To(gs.MatchAllKeys(gs.Keys{
				"urn:pulumi:stack::project::kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::" +
					"crontabs.stable.example.com": pgm.MatchProps(
					gs.IgnoreExtras,
					pgm.Props{
						"state": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
							"metadata": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
								"name": pgm.MatchValue("crontabs.stable.example.com"),
							}),
						}),
					},
				),
				"urn:pulumi:stack::project::kubernetes:core/v1:Namespace::my-namespace": pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
					"state": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
						"metadata": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
							"name": pgm.MatchValue("my-namespace"),
						}),
					}),
				}),
				"urn:pulumi:stack::project::kubernetes:core/v1:ConfigMap::my-namespace/my-map": pgm.MatchProps(
					gs.IgnoreExtras,
					pgm.Props{
						"state": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
							"metadata": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
								"name": pgm.MatchValue("my-map"),
							}),
						}),
					},
				),
				"urn:pulumi:stack::project::kubernetes:stable.example.com/v1:CronTab::" +
					"my-namespace/my-new-cron-object": pgm.MatchProps(
					gs.IgnoreExtras,
					pgm.Props{
						"state": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
							"metadata": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
								"name": pgm.MatchValue("my-new-cron-object"),
							}),
						}),
					},
				),
			}))
		})

		gk.It("should return an array of resource outputs", func(ctx context.Context) {
			resources, err := register(ctx)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())

			resourceArray, err := internals.UnsafeAwaitOutput(ctx, resources)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(resourceArray.Value).To(gm.HaveLen(4))
		})

		gk.Context("when a prefix is configured", func() {
			gk.BeforeEach(func() {
				registerOpts.ResourcePrefix = "prefixed"
			})
			gk.It("should apply the prefix to the resource (but not to the Kubernetes object)", func(ctx context.Context) {
				_, err := register(ctx)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())

				gm.Expect(tc.monitor.Resources()).To(gs.MatchAllKeys(gs.Keys{
					"urn:pulumi:stack::project::kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::" +
						"prefixed:crontabs.stable.example.com": pgm.MatchProps(
						gs.IgnoreExtras,
						pgm.Props{
							"state": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
								"metadata": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
									"name": pgm.MatchValue("crontabs.stable.example.com"),
								}),
							}),
						},
					),
					"urn:pulumi:stack::project::kubernetes:core/v1:Namespace::prefixed:my-namespace": pgm.MatchProps(
						gs.IgnoreExtras,
						pgm.Props{
							"state": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
								"metadata": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
									"name": pgm.MatchValue("my-namespace"),
								}),
							}),
						},
					),
					"urn:pulumi:stack::project::kubernetes:core/v1:ConfigMap::prefixed:my-namespace/my-map": pgm.MatchProps(
						gs.IgnoreExtras,
						pgm.Props{
							"state": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
								"metadata": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
									"name": pgm.MatchValue("my-map"),
								}),
							}),
						},
					),
					"urn:pulumi:stack::project::kubernetes:stable.example.com/v1:CronTab::" +
						"prefixed:my-namespace/my-new-cron-object": pgm.MatchProps(
						gs.IgnoreExtras,
						pgm.Props{
							"state": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
								"metadata": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
									"name": pgm.MatchValue("my-new-cron-object"),
								}),
							}),
						},
					),
				}))
			})
		})

		gk.Context("when skipAwait is enabled", func() {
			gk.BeforeEach(func() {
				registerOpts.SkipAwait = true
			})
			gk.It("should apply the skipAwait annotation", func(ctx context.Context) {
				_, err := register(ctx)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())

				gm.Expect(tc.monitor.Resources()).To(gs.MatchAllKeys(gs.Keys{
					"urn:pulumi:stack::project::kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::" +
						"crontabs.stable.example.com": pgm.MatchProps(
						gs.IgnoreExtras,
						pgm.Props{
							"state": pgm.BeObject(pgm.HaveSkipAwaitAnnotation()),
						},
					),
					"urn:pulumi:stack::project::kubernetes:core/v1:Namespace::my-namespace": pgm.MatchProps(
						gs.IgnoreExtras,
						pgm.Props{
							"state": pgm.BeObject(pgm.HaveSkipAwaitAnnotation()),
						},
					),
					"urn:pulumi:stack::project::kubernetes:core/v1:ConfigMap::my-namespace/my-map": pgm.MatchProps(
						gs.IgnoreExtras,
						pgm.Props{
							"state": pgm.BeObject(pgm.HaveSkipAwaitAnnotation()),
						},
					),
					"urn:pulumi:stack::project::kubernetes:stable.example.com/v1:CronTab::" +
						"my-namespace/my-new-cron-object": pgm.MatchProps(
						gs.IgnoreExtras,
						pgm.Props{
							"state": pgm.BeObject(pgm.HaveSkipAwaitAnnotation()),
						},
					),
				}))
			})
		})

		gk.Describe("Ordering", func() {
			gk.Context("implicit dependencies", func() {
				gk.It("should apply a DependsOn option on the dependents", func(ctx context.Context) {
					_, err := register(ctx)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())

					gm.Expect(tc.monitor.Registrations()).To(gs.MatchAllKeys(gs.Keys{
						"urn:pulumi:stack::project::kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::" +
							"crontabs.stable.example.com": gs.MatchFields(
							gs.IgnoreExtras,
							gs.Fields{
								"Request": gs.MatchFields(gs.IgnoreExtras, gs.Fields{
									"Dependencies": gm.BeEmpty(),
								}),
							},
						),
						"urn:pulumi:stack::project::kubernetes:core/v1:Namespace::my-namespace": gs.MatchFields(
							gs.IgnoreExtras,
							gs.Fields{
								"Request": gs.MatchFields(gs.IgnoreExtras, gs.Fields{
									"Dependencies": gm.BeEmpty(),
								}),
							},
						),
						"urn:pulumi:stack::project::kubernetes:core/v1:ConfigMap::my-namespace/my-map": gs.MatchFields(
							gs.IgnoreExtras,
							gs.Fields{
								"Request": gs.MatchFields(gs.IgnoreExtras, gs.Fields{
									"Dependencies": gm.ConsistOf(
										"urn:pulumi:stack::project::kubernetes:core/v1:Namespace::my-namespace",
									),
								}),
							},
						),
						"urn:pulumi:stack::project::kubernetes:stable.example.com/v1:CronTab::" +
							"my-namespace/my-new-cron-object": gs.MatchFields(
							gs.IgnoreExtras,
							gs.Fields{
								"Request": gs.MatchFields(gs.IgnoreExtras, gs.Fields{
									"Dependencies": gm.ConsistOf(
										"urn:pulumi:stack::project::kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::"+
											"crontabs.stable.example.com",
										"urn:pulumi:stack::project::kubernetes:core/v1:Namespace::my-namespace",
									),
								}),
							},
						),
					}))
				})
			})

			gk.Context("explicit dependencies (config.kubernetes.io/depends-on annotation)", func() {
				gk.BeforeEach(func() {
					registerOpts.Objects = append(registerOpts.Objects, unstructured.Unstructured{
						Object: map[string]any{
							"apiVersion": "v1",
							"kind":       "Pod",
							"metadata": map[string]any{
								"name": "my-pod",
								"annotations": map[string]any{
									"config.kubernetes.io/depends-on": "/Namespace/my-namespace",
								},
							},
						},
					})
				})
				gk.It("should apply a DependsOn option on the dependents", func(ctx context.Context) {
					_, err := register(ctx)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())

					gm.Expect(tc.monitor.Registrations()).To(gs.MatchKeys(gs.IgnoreExtras, gs.Keys{
						"urn:pulumi:stack::project::kubernetes:core/v1:Pod::my-pod": gs.MatchFields(gs.IgnoreExtras, gs.Fields{
							"Request": gs.MatchFields(gs.IgnoreExtras, gs.Fields{
								"Dependencies": gm.ConsistOf(
									"urn:pulumi:stack::project::kubernetes:core/v1:Namespace::my-namespace",
								),
							}),
						}),
					}))
				})
			})
		})
	})

	gk.Describe("Kubernetes object specifics", func() {
		gk.Context("when the object is a Secret", func() {
			gk.BeforeEach(func() {
				registerOpts.Objects = []unstructured.Unstructured{{
					Object: map[string]any{
						"apiVersion": "v1",
						"kind":       "Secret",
						"metadata": map[string]any{
							"name": "my-secret",
						},
						"stringData": map[string]interface{}{
							"foo": "bar",
						},
					},
				}}
			})
			gk.It("should mark the contents as a secret", func(ctx context.Context) {
				_, err := register(ctx)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())

				gm.Expect(tc.monitor.Resources()).To(gs.MatchKeys(gs.IgnoreExtras, gs.Keys{
					"urn:pulumi:stack::project::kubernetes:core/v1:Secret::my-secret": pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"state": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
							"stringData": pgm.BeSecret(),
						}),
					}),
				}))
			})
		})
	})
})

var _ = gk.Describe("Parse", func() {
	var args ParseOptions

	gk.BeforeEach(func() {
		args = ParseOptions{}
	})

	tempFiles := func(manifests ...string) string {
		tempDir := gk.GinkgoTB().TempDir()
		for i, m := range manifests {
			name := filepath.Join(tempDir, fmt.Sprintf("manifest-%02d.yaml", i+1))
			err := os.WriteFile(name, []byte(m), 0o600)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
		}
		return tempDir
	}

	manifestAssertions := func() {
		gk.GinkgoHelper()

		gk.It("should produce the objects in the manifest", func(ctx context.Context) {
			objs, err := Parse(ctx, args)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(objs).To(gm.ConsistOf(
				matchUnstructured(
					gs.Keys{"metadata": gs.MatchKeys(gs.IgnoreExtras, gs.Keys{"name": gm.Equal("my-namespace")})}),
				matchUnstructured(
					gs.Keys{"metadata": gs.MatchKeys(gs.IgnoreExtras, gs.Keys{"name": gm.Equal("crontabs.stable.example.com")})},
				),
				matchUnstructured(gs.Keys{"metadata": gs.MatchKeys(gs.IgnoreExtras, gs.Keys{"name": gm.Equal("my-map")})}),
				matchUnstructured(
					gs.Keys{"metadata": gs.MatchKeys(gs.IgnoreExtras, gs.Keys{"name": gm.Equal("my-new-cron-object")})}),
			))
		})
	}

	gk.Describe("yamls", func() {
		gk.Context("when the input is a valid YAML string", func() {
			gk.BeforeEach(func() {
				args.YAML = manifest
			})
			manifestAssertions()
		})

		gk.Context("when the manifest has empty blocks", func() {
			gk.BeforeEach(func() {
				args.YAML = "---"
			})
			gk.It("should do nothing", func(ctx context.Context) {
				_, err := Parse(ctx, args)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
			})
		})
	})

	gk.Describe("files", func() {
		gk.Describe("globbing", func() {
			gk.BeforeEach(func() {
				args.Glob = true
			})

			gk.Context("when the pattern matches no files", func() {
				gk.BeforeEach(func() {
					args.Files = []string{"nosuchfile-*.yaml"}
				})
				gk.It("should do nothing", func(ctx context.Context) {
					_, err := Parse(ctx, args)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
				})
			})

			gk.Context("when the pattern matches some files", func() {
				gk.BeforeEach(func() {
					tempDir := tempFiles(manifest)
					args.Files = []string{filepath.Join(tempDir, "*.yaml")}
				})
				manifestAssertions()
			})
		})

		gk.Context("when the file doesn't exist", func() {
			gk.BeforeEach(func() {
				args.Files = []string{"nosuchfile.yaml"}
			})
			gk.It("should fail", func(ctx context.Context) {
				_, err := Parse(ctx, args)
				gm.Expect(err).Should(gm.HaveOccurred())
			})
		})

		gk.Context("when the input is a valid YAML file", func() {
			gk.BeforeEach(func() {
				tempDir := tempFiles(manifest)
				args.Files = []string{filepath.Join(tempDir, "manifest-01.yaml")}
			})
			manifestAssertions()
		})

		gk.Context("when the input is a valid YAML URL", func() {
			gk.BeforeEach(func() {
				args.Files = []string{
					`https://raw.githubusercontent.com/pulumi/pulumi-kubernetes/master/tests/sdk/nodejs/examples/` +
						`yaml-guestbook/yaml/guestbook.yaml`,
				}
			})
			gk.It("should download and use the document", func(ctx context.Context) {
				objs, err := Parse(ctx, args)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(objs).To(gm.HaveLen(6))
			})
		})
	})
})

var _ = gk.Describe("Normalize", func() {
	var objs []unstructured.Unstructured
	var defaultNamespace string
	var clientSet *clients.DynamicClientSet
	var disco *fake.SimpleDiscovery

	gk.BeforeEach(func() {
		objs = []unstructured.Unstructured{}
		defaultNamespace = "default"
		clientSet, disco, _, _ = fake.NewSimpleDynamicClient()

		// populate the discovery client with some custom resources
		var fakeResources = []*metav1.APIResourceList{
			{
				GroupVersion: "stable.example.com/v1",
				APIResources: []metav1.APIResource{
					{Name: "issuers", Namespaced: true, Kind: "Issuer"},
					{Name: "clusterissuers", Namespaced: false, Kind: "ClusterIssuer"},
				},
			},
		}
		disco.Resources = append(disco.Resources, fakeResources...)
	})

	gk.Describe("validation", func() {
		gk.Context("when the object has no kind", func() {
			gk.BeforeEach(func() {
				objs = []unstructured.Unstructured{{
					Object: map[string]any{
						"apiVersion": "v1",
						"metadata":   map[string]any{},
					},
				}}
			})
			gk.It("should fail", func(_ /* ctx */ context.Context) {
				_, err := Normalize(objs, defaultNamespace, clientSet)
				gm.Expect(err).Should(gm.MatchError(gm.ContainSubstring("Kubernetes resources require a kind and apiVersion")))
			})
		})

		gk.Context("when the object has no metadata.name", func() {
			gk.BeforeEach(func() {
				objs = []unstructured.Unstructured{{
					Object: map[string]any{
						"apiVersion": "v1",
						"kind":       "Secret",
						"metadata":   map[string]any{},
					},
				}}
			})
			gk.It("should fail", func(_ /* ctx */ context.Context) {
				_, err := Normalize(objs, defaultNamespace, clientSet)
				gm.Expect(err).Should(gm.MatchError(gm.ContainSubstring("Kubernetes resources require a .metadata.name")))
			})
		})
	})

	gk.Describe("namespacing", func() {
		gk.Context("when the object has a namespace-scoped kind", func() {
			gk.BeforeEach(func() {
				objs = []unstructured.Unstructured{{
					Object: map[string]any{
						"apiVersion": "v1",
						"kind":       "Secret",
						"metadata": map[string]any{
							"name": "my-secret",
						},
					},
				}}
			})

			gk.It("should apply the default namespace", func(_ /* ctx */ context.Context) {
				objs, err := Normalize(objs, defaultNamespace, clientSet)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(objs).To(gm.HaveExactElements(
					matchUnstructured(gs.Keys{"metadata": gs.MatchKeys(gs.IgnoreExtras, gs.Keys{"namespace": gm.Equal("default")})}),
				))
			})
		})

		gk.Context("when the object has a cluster-scoped kind", func() {
			gk.BeforeEach(func() {
				objs = []unstructured.Unstructured{{
					Object: map[string]any{
						"apiVersion": "rbac.authorization.k8s.io/v1",
						"kind":       "ClusterRole",
						"metadata": map[string]any{
							"name": "my-secret",
						},
					},
				}}
			})

			gk.It("should not apply the default namespace", func(_ /* ctx */ context.Context) {
				objs, err := Normalize(objs, defaultNamespace, clientSet)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(objs).To(gm.HaveExactElements(
					matchUnstructured(gs.Keys{"metadata": gm.Not(gm.HaveKey("namespace"))}),
				))
			})
		})
	})

	gk.Describe("special-case kinds", func() {
		gk.Context("when the object is a list", func() {
			gk.BeforeEach(func() {
				resources, err := yamlDecode(list)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				objs = resources
			})
			gk.It("should flatten the list", func(_ /* ctx */ context.Context) {
				objs, err := Normalize(objs, defaultNamespace, clientSet)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(objs).To(gm.HaveExactElements(
					matchUnstructured(gs.Keys{"metadata": gs.MatchKeys(gs.IgnoreExtras, gs.Keys{"name": gm.Equal("map-1")})}),
					matchUnstructured(gs.Keys{"metadata": gs.MatchKeys(gs.IgnoreExtras, gs.Keys{"name": gm.Equal("map-2")})}),
					matchUnstructured(gs.Keys{"metadata": gs.MatchKeys(gs.IgnoreExtras, gs.Keys{"name": gm.Equal("map-3")})}),
				))
			})
		})

		gk.Context("when the object has kind 'core/v1'", func() {
			gk.BeforeEach(func() {
				objs = []unstructured.Unstructured{{
					Object: map[string]any{
						"apiVersion": "core/v1",
						"kind":       "Secret",
						"metadata": map[string]any{
							"name": "my-secret",
						},
					},
				}}
			})

			gk.It("should replace with 'v1", func(_ /* ctx */ context.Context) {
				objs, err := Normalize(objs, defaultNamespace, clientSet)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(objs).To(gm.HaveExactElements(
					matchUnstructured(gs.Keys{"apiVersion": gm.Equal("v1"), "kind": gm.Equal("Secret")}),
				))
			})
		})
	})
})

func TestIsGlobPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pattern  string
		expected bool
	}{
		{pattern: `manifest.yaml`, expected: false},
		{pattern: `*.yaml`, expected: true},
		{pattern: `*`, expected: true},
		{pattern: `test-?.yaml`, expected: true},
		{pattern: `ba[rz].yaml`, expected: true},
		{pattern: `escaped-\*.yaml`, expected: false},
		{pattern: `\*.yaml`, expected: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()

			isPattern := isGlobPattern(tt.pattern)
			if tt.expected {
				assert.Truef(t, isPattern, "expected %q to be a pattern", tt.pattern)
			} else {
				assert.Falsef(t, isPattern, "expected %q to not be a pattern", tt.pattern)
			}
		})
	}
}

func matchUnstructured(keys gs.Keys) gomegatypes.GomegaMatcher {
	return gm.WithTransform(func(obj unstructured.Unstructured) map[string]interface{} {
		return obj.Object
	}, gs.MatchKeys(gs.IgnoreExtras, keys))
}
