// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	k8s "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/kustomize"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		provider, err := k8s.NewProvider(ctx, "k8s", &k8s.ProviderArgs{})
		if err != nil {
			return err
		}
		ns, err := corev1.NewNamespace(ctx, "test", &corev1.NamespaceArgs{}, pulumi.Provider(provider))
		if err != nil {
			return err
		}
		nsProvider, err := k8s.NewProvider(ctx, "k8s-ns", &k8s.ProviderArgs{
			Namespace: ns.Metadata.Name(),
		})
		if err != nil {
			return err
		}

		_, err = kustomize.NewDirectory(ctx, "helloWorld",
			kustomize.DirectoryArgs{Directory: pulumi.String("helloWorld")},
			pulumi.Provider(nsProvider),
		)
		if err != nil {
			return err
		}

		return nil
	})
}
