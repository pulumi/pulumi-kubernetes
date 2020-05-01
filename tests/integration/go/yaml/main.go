package main

import (
	"path/filepath"

	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yaml.NewConfigFile(ctx, "guestbook",
			&yaml.ConfigFileArgs{File: "guestbook-all-in-one.yaml"},
		)
		if err != nil {
			return err
		}

		_, err = yaml.NewConfigGroup(ctx, "guestbook",
			&yaml.ConfigGroupArgs{Files: []string{filepath.Join("yaml", "*.yaml")}},
		)
		return err
	})
}
