package main

import (
	k8s "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/kustomize"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		ns, err := corev1.NewNamespace(ctx, "test", &corev1.NamespaceArgs{})
		if err != nil {
			return err
		}
		provider, err := k8s.NewProvider(ctx, "k8s", &k8s.ProviderArgs{
			Kubeconfig: pulumi.String("~/.kube/config"),
			Namespace:  ns.Metadata.Name(),
		})
		if err != nil {
			return err
		}
		_, err = kustomize.NewDirectory(ctx, "helloWorld",
			kustomize.DirectoryArgs{Directory: pulumi.String("helloWorld")},
			pulumi.Provider(provider),
		)
		if err != nil {
			return err
		}

		return nil
	})
}
