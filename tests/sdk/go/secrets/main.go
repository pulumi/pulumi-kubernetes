package main

import (
	b64 "encoding/base64"
	"fmt"
	"math/rand"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		pw := conf.RequireSecret("message")
		rawPW := conf.Require("message")

		cmData, err := corev1.NewConfigMap(ctx, "cmdata", &corev1.ConfigMapArgs{
			Data: pulumi.StringMap{
				"password": pw,
			},
		})
		if err != nil {
			return err
		}

		cmBinaryData, err := corev1.NewConfigMap(ctx, "cmbinarydata", &corev1.ConfigMapArgs{
			Data: pulumi.StringMap{
				"password": pw.ApplyT(func(s string) string {
					return b64.StdEncoding.EncodeToString([]byte(s))
				}).(pulumi.StringOutput),
			},
		})
		if err != nil {
			return err
		}

		sStringData, err := corev1.NewSecret(ctx, "sstringdata", &corev1.SecretArgs{
			StringData: pulumi.StringMap{
				"password": pulumi.String(rawPW),
			},
		})
		if err != nil {
			return err
		}

		sData, err := corev1.NewSecret(ctx, "sdata", &corev1.SecretArgs{
			Data: pulumi.StringMap{
				"password": pulumi.String(b64.StdEncoding.EncodeToString([]byte(rawPW))),
			},
		})
		if err != nil {
			return err
		}

		randSuffix := func() string {
			b := make([]rune, 5)
			letters := []rune("abcdefghijklmnopqrstuvwxyz")
			for i := range b {
				//nolint:gosec
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		}
		name := fmt.Sprintf("test-%s", randSuffix())

		secretYAML := fmt.Sprintf(`
apiVersion: v1
kind: Secret
metadata:
  name: %s
stringData:
  password: %s
`, name, rawPW)
		cg, err := yaml.NewConfigGroup(ctx, "example",
			&yaml.ConfigGroupArgs{
				YAML: []string{secretYAML},
			})
		if err != nil {
			return err
		}
		cgSecret := cg.GetResource("v1/Secret", name, "")

		ctx.Export("cmData", cmData.Data)
		ctx.Export("cmBinaryData", cmBinaryData.BinaryData)
		ctx.Export("sStringData", sStringData.StringData)
		ctx.Export("sData", sData.Data)
		ctx.Export("cgData", cgSecret.(*corev1.Secret).Data)

		return nil
	})
}
