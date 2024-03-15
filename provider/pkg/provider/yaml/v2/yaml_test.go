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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	. "github.com/onsi/gomega/gstruct"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	. "github.com/pulumi/pulumi-kubernetes/tests/v4/gomega"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/internals"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	manifest = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-map
data:
  altGreeting: "Good Morning!"
---
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: my-new-cron-object
spec:
  cronSpec: "* * * * */5"
  image: my-awesome-cron-image
`
)

var _ = Describe("Register", func() {
	var tc *componentProviderTestContext
	var registerOpts *RegisterOptions

	BeforeEach(func() {
		tc = newTestContext(GinkgoTB())
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

	Context("given the objects in the manifest", func() {
		BeforeEach(func() {
			resources, err := yamlDecode(manifest, nil)
			Expect(err).ShouldNot(HaveOccurred())
			registerOpts.Objects = resources
		})

		It("should register some resources", func(ctx context.Context) {
			_, err := register(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(tc.monitor.Resources()).To(MatchKeys(IgnoreExtras, Keys{
				"urn:pulumi:stack::project::kubernetes:core/v1:ConfigMap::my-map": MatchProps(IgnoreExtras, Props{
					"state": MatchObject(IgnoreExtras, Props{
						"metadata": MatchObject(IgnoreExtras, Props{
							"name": MatchValue("my-map"),
						}),
					}),
				}),
				"urn:pulumi:stack::project::kubernetes:stable.example.com/v1:CronTab::my-new-cron-object": MatchProps(IgnoreExtras, Props{
					"state": MatchObject(IgnoreExtras, Props{
						"metadata": MatchObject(IgnoreExtras, Props{
							"name": MatchValue("my-new-cron-object"),
						}),
					}),
				}),
			}))
		})

		It("should return an array of resource outputs", func(ctx context.Context) {
			resources, err := register(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			resourceArray, err := internals.UnsafeAwaitOutput(ctx, resources)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resourceArray.Value).To(HaveLen(2))
		})

		Context("when a prefix is configured", func() {
			BeforeEach(func() {
				registerOpts.ResourcePrefix = "prefixed"
			})
			It("should apply the prefix to the resource (but not to the Kubernetes object)", func(ctx context.Context) {
				_, err := register(ctx)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(tc.monitor.Resources()).To(MatchKeys(IgnoreExtras, Keys{
					"urn:pulumi:stack::project::kubernetes:core/v1:ConfigMap::prefixed-my-map": MatchProps(IgnoreExtras, Props{
						"state": MatchObject(IgnoreExtras, Props{
							"metadata": MatchObject(IgnoreExtras, Props{
								"name": MatchValue("my-map"),
							}),
						}),
					}),
					"urn:pulumi:stack::project::kubernetes:stable.example.com/v1:CronTab::prefixed-my-new-cron-object": MatchProps(IgnoreExtras, Props{
						"state": MatchObject(IgnoreExtras, Props{
							"metadata": MatchObject(IgnoreExtras, Props{
								"name": MatchValue("my-new-cron-object"),
							}),
						}),
					}),
				}))
			})
		})

		Context("when skipAwait is enabled", func() {
			BeforeEach(func() {
				registerOpts.SkipAwait = true
			})
			It("should apply the skipAwait annotation", func(ctx context.Context) {
				_, err := register(ctx)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(tc.monitor.Resources()).To(MatchKeys(IgnoreExtras, Keys{
					"urn:pulumi:stack::project::kubernetes:core/v1:ConfigMap::my-map": MatchProps(IgnoreExtras, Props{
						"state": MatchObject(IgnoreExtras, Props{
							"metadata": MatchObject(IgnoreExtras, Props{
								"annotations": MatchObject(IgnoreExtras, Props{
									"pulumi.com/skipAwait": MatchValue("true"),
								}),
							}),
						}),
					}),
					"urn:pulumi:stack::project::kubernetes:stable.example.com/v1:CronTab::my-new-cron-object": MatchProps(IgnoreExtras, Props{
						"state": MatchObject(IgnoreExtras, Props{
							"metadata": MatchObject(IgnoreExtras, Props{
								"annotations": MatchObject(IgnoreExtras, Props{
									"pulumi.com/skipAwait": MatchValue("true"),
								}),
							}),
						}),
					}),
				}))
			})
		})
	})

	Describe("Kubernetes object specifics", func() {
		Context("when the object has no kind", func() {
			BeforeEach(func() {
				registerOpts.Objects = []unstructured.Unstructured{{
					Object: map[string]any{
						"apiVersion": "v1",
						"metadata":   map[string]any{},
					},
				}}
			})
			It("should fail", func(ctx context.Context) {
				_, err := register(ctx)
				Expect(err).Should(MatchError(ContainSubstring("Kubernetes resources require a kind and apiVersion")))
			})
		})

		Context("when the object has no metadata.name", func() {
			BeforeEach(func() {
				registerOpts.Objects = []unstructured.Unstructured{{
					Object: map[string]any{
						"apiVersion": "v1",
						"kind":       "Secret",
						"metadata":   map[string]any{},
					},
				}}
			})
			It("should fail", func(ctx context.Context) {
				_, err := register(ctx)
				Expect(err).Should(MatchError(ContainSubstring("YAML object does not have a .metadata.name")))
			})
		})

		Context("when the object is a list", func() {
			BeforeEach(func() {
				registerOpts.Objects = []unstructured.Unstructured{{
					Object: map[string]any{
						"apiVersion": "v1",
						"kind":       "List",
						"items": []any{
							map[string]any{
								"apiVersion": "v1",
								"kind":       "Secret",
								"metadata": map[string]any{
									"name": "my-secret",
								},
							},
						},
					},
				}}
			})
			It("should flatten the list", func(ctx context.Context) {
				_, err := register(ctx)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(tc.monitor.Resources()).To(HaveKey("urn:pulumi:stack::project::kubernetes:core/v1:Secret::my-secret"))
			})
		})

		Context("when the object is a Secret", func() {
			BeforeEach(func() {
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
			It("should mark the contents as a secret", func(ctx context.Context) {
				_, err := register(ctx)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(tc.monitor.Resources()).To(MatchKeys(IgnoreExtras, Keys{
					"urn:pulumi:stack::project::kubernetes:core/v1:Secret::my-secret": MatchProps(IgnoreExtras, Props{
						"state": MatchObject(IgnoreExtras, Props{
							"stringData": BeSecret(),
						}),
					}),
				}))
			})
		})
	})
})

var _ = Describe("Parse", func() {
	var clientSet *clients.DynamicClientSet
	var args ParseOptions

	BeforeEach(func() {
		args = ParseOptions{}
	})

	JustBeforeEach(func() {
	})

	parse := func(ctx context.Context) (objs []unstructured.Unstructured, err error) {
		return Parse(context.Background(), clientSet, args)
	}

	tempFiles := func(manifests ...string) string {
		tempDir := GinkgoTB().TempDir()
		for i, m := range manifests {
			name := filepath.Join(tempDir, fmt.Sprintf("manifest-%02d.yaml", i+1))
			err := os.WriteFile(name, []byte(m), 0o600)
			Expect(err).ShouldNot(HaveOccurred())
		}
		return tempDir
	}

	manifestAssertions := func() {
		GinkgoHelper()

		It("should produce the objects in the manifest", func(ctx context.Context) {
			objs, err := parse(ctx)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(objs).To(HaveExactElements(
				matchUnstructured(Keys{
					"metadata": MatchKeys(IgnoreExtras, Keys{
						"name": Equal("my-map"),
					}),
				}),
				matchUnstructured(Keys{
					"metadata": MatchKeys(IgnoreExtras, Keys{
						"name": Equal("my-new-cron-object"),
					}),
				}),
			))
		})
	}

	Describe("yamls", func() {
		Context("when the input is a valid YAML string", func() {
			BeforeEach(func() {
				args.YAML = manifest
			})
			manifestAssertions()
		})

		Context("when the manifest has empty blocks", func() {
			BeforeEach(func() {
				args.YAML = "---"
			})
			It("should do nothing", func(ctx context.Context) {
				_, err := parse(ctx)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("files", func() {
		Describe("globbing", func() {
			BeforeEach(func() {
				args.Glob = true
			})

			Context("when the pattern matches no files", func() {
				BeforeEach(func() {
					args.Files = []string{"nosuchfile-*.yaml"}
				})
				It("should do nothing", func(ctx context.Context) {
					_, err := parse(ctx)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			Context("when the pattern matches some files", func() {
				BeforeEach(func() {
					tempDir := tempFiles(manifest)
					args.Files = []string{filepath.Join(tempDir, "*.yaml")}
				})
				manifestAssertions()
			})
		})

		Context("when the file doesn't exist", func() {
			BeforeEach(func() {
				args.Files = []string{"nosuchfile.yaml"}
			})
			It("should fail", func(ctx context.Context) {
				_, err := parse(ctx)
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("when the input is a valid YAML file", func() {
			BeforeEach(func() {
				tempDir := tempFiles(manifest)
				args.Files = []string{filepath.Join(tempDir, "manifest-01.yaml")}
			})
			manifestAssertions()
		})

		Context("when the input is a valid YAML URL", func() {
			BeforeEach(func() {
				args.Files = []string{`https://raw.githubusercontent.com/pulumi/pulumi-kubernetes/master/tests/sdk/nodejs/examples/yaml-guestbook/yaml/guestbook.yaml`}
			})
			It("should download and use the document", func(ctx context.Context) {
				objs, err := parse(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(objs).To(HaveLen(6))
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

func matchUnstructured(keys gstruct.Keys) gomegatypes.GomegaMatcher {
	return WithTransform(func(obj unstructured.Unstructured) map[string]interface{} {
		return obj.Object
	}, MatchKeys(IgnoreExtras, keys))
}
