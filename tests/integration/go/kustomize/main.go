package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/kustomize"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := kustomize.NewDirectory(ctx, "helloWorld",
			kustomize.DirectoryArgs{Directory: pulumi.String("helloWorld")},
		)
		if err != nil {
			return err
		}

		return nil
	})
}
