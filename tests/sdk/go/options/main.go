package main

import (
	"strings"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/kustomize"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-time/sdk/go/time"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		bootstrapProvider, err := kubernetes.NewProvider(ctx, "bootstrap", &kubernetes.ProviderArgs{})
		if err != nil {
			return err
		}
		nulloptsNs, err := corev1.NewNamespace(ctx, "nullopts", &corev1.NamespaceArgs{}, pulumi.Provider(bootstrapProvider))
		if err != nil {
			return err
		}
		aNs, err := corev1.NewNamespace(ctx, "a", &corev1.NamespaceArgs{}, pulumi.Provider(bootstrapProvider))
		if err != nil {
			return err
		}
		bNs, err := corev1.NewNamespace(ctx, "b", &corev1.NamespaceArgs{}, pulumi.Provider(bootstrapProvider))
		if err != nil {
			return err
		}
		nulloptsProvider, err := kubernetes.NewProvider(ctx, "nullopts", &kubernetes.ProviderArgs{Namespace: nulloptsNs.Metadata.Name()})
		if err != nil {
			return err
		}
		aProvider, err := kubernetes.NewProvider(ctx, "a", &kubernetes.ProviderArgs{Namespace: aNs.Metadata.Name()})
		if err != nil {
			return err
		}
		bProvider, err := kubernetes.NewProvider(ctx, "b", &kubernetes.ProviderArgs{Namespace: bNs.Metadata.Name()})
		if err != nil {
			return err
		}

		// a sleep resource to exercise the "depends_on" component-level option
		sleep, err := time.NewSleep(ctx, "sleep", &time.SleepArgs{CreateDuration: pulumi.String("1s")})
		if err != nil {
			return err
		}

		// apply_default_opts is a stack transformation that applies default opts to any resource whose name ends with "-nullopts".
		// this is intended to be applied to component resources only.
		applyDefaultOpts := func(args *pulumi.ResourceTransformationArgs) *pulumi.ResourceTransformationResult {
			if strings.HasSuffix(args.Name, "-nullopts") {
				return &pulumi.ResourceTransformationResult{
					Props: args.Props,
					Opts:  append(args.Opts, pulumi.Provider(nulloptsProvider)),
				}
			}
			return nil
		}
		ctx.RegisterStackTransformation(applyDefaultOpts)

		// applyAlias is a Pulumi transformation that applies a unique alias to each resource.
		applyAlias := func(args *pulumi.ResourceTransformationArgs) *pulumi.ResourceTransformationResult {
			return &pulumi.ResourceTransformationResult{
				Props: args.Props,
				Opts:  append(args.Opts, pulumi.Aliases([]pulumi.Alias{{Name: pulumi.String(args.Name + "-aliased")}})),
			}
		}

		// transform_k8s is a Kubernetes transformation that applies a unique alias and annotation to each resource.
		transformK8s := func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
			metadata := state["metadata"].(map[string]interface{})
			metadata["annotations"] = map[string]interface{}{"transformed": "true"}

			// note: pulumi-kubernetes Go SDK doesn't provide a way to mutate the options (e.g. to add a "-k8s-aliased" alias)
			// https://github.com/pulumi/pulumi-kubernetes/issues/2666
		}

		// --- ConfigGroup ---

		_, err = yaml.NewConfigGroup(ctx, "cg-options", &yaml.ConfigGroupArgs{
			ResourcePrefix:  "cg-options",
			SkipAwait:       true,
			Transformations: []yaml.Transformation{transformK8s},
			Files:           []string{"./testdata/options/configgroup/*.yaml"},
			YAML: []string{`
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-options-cm-1
`},
		},
			pulumi.Providers(aProvider),
			pulumi.Aliases([]pulumi.Alias{{Name: pulumi.String("cg-options-old")}}),
			pulumi.IgnoreChanges([]string{"ignored"}),
			pulumi.Protect(true),
			pulumi.DependsOn([]pulumi.Resource{sleep}),
			pulumi.Transformations([]pulumi.ResourceTransformation{applyAlias}),
			pulumi.Version("1.2.3"),
			pulumi.PluginDownloadURL("https://a.pulumi.test"),
		)
		if err != nil {
			return err
		}

		// "provider" option
		_, err = yaml.NewConfigGroup(ctx, "cg-provider", &yaml.ConfigGroupArgs{
			ResourcePrefix: "cg-provider",
			YAML: []string{`
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-provider-cm-1
`},
		}, pulumi.Provider(bProvider))
		if err != nil {
			return err
		}

		// null opts
		_, err = yaml.NewConfigGroup(ctx, "cg-nullopts", &yaml.ConfigGroupArgs{
			ResourcePrefix: "cg-nullopts",
			YAML: []string{`
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-nullopts-cm-1
`},
		})
		if err != nil {
			return err
		}

		//  --- ConfigFile ---
		_, err = yaml.NewConfigFile(ctx, "cf-options", &yaml.ConfigFileArgs{
			ResourcePrefix:  "cf-options",
			SkipAwait:       true,
			Transformations: []yaml.Transformation{transformK8s},
			File:            "./testdata/options/configfile/manifest.yaml",
		},
			pulumi.Providers(aProvider),
			pulumi.Aliases([]pulumi.Alias{{Name: pulumi.String("cf-options-old")}}),
			pulumi.IgnoreChanges([]string{"ignored"}),
			pulumi.Protect(true),
			pulumi.DependsOn([]pulumi.Resource{sleep}),
			pulumi.Transformations([]pulumi.ResourceTransformation{applyAlias}),
			pulumi.Version("1.2.3"),
			pulumi.PluginDownloadURL("https://a.pulumi.test"),
		)
		if err != nil {
			return err
		}

		// "provider" option
		_, err = yaml.NewConfigFile(ctx, "cf-provider", &yaml.ConfigFileArgs{
			ResourcePrefix: "cf-provider",
			File:           "./testdata/options/configfile/manifest.yaml",
		}, pulumi.Providers(bProvider))
		if err != nil {
			return err
		}

		// null opts
		_, err = yaml.NewConfigFile(ctx, "cf-nullopts", &yaml.ConfigFileArgs{
			ResourcePrefix: "cf-nullopts",
			File:           "./testdata/options/configfile/manifest.yaml",
		})
		if err != nil {
			return err
		}

		// empty manifest
		_, err = yaml.NewConfigFile(ctx, "cf-empty", &yaml.ConfigFileArgs{
			ResourcePrefix: "cf-empty",
			File:           "./testdata/options/configfile/empty.yaml",
		}, pulumi.Providers(bProvider))
		if err != nil {
			return err
		}

		// --- Directory ---
		_, err = kustomize.NewDirectory(ctx, "kustomize-options", kustomize.DirectoryArgs{
			Directory:       pulumi.String("./testdata/options/kustomize"),
			ResourcePrefix:  "kustomize-options",
			Transformations: []yaml.Transformation{transformK8s},
		},
			pulumi.Providers(aProvider),
			pulumi.Aliases([]pulumi.Alias{{Name: pulumi.String("kustomize-options-old")}}),
			pulumi.IgnoreChanges([]string{"ignored"}),
			pulumi.Protect(true),
			pulumi.DependsOn([]pulumi.Resource{sleep}),
			pulumi.Transformations([]pulumi.ResourceTransformation{applyAlias}),
			pulumi.Version("1.2.3"),
			pulumi.PluginDownloadURL("https://a.pulumi.test"),
		)
		if err != nil {
			return err
		}

		// "provider" option
		_, err = kustomize.NewDirectory(ctx, "kustomize-provider", kustomize.DirectoryArgs{
			Directory:      pulumi.String("./testdata/options/kustomize"),
			ResourcePrefix: "kustomize-provider",
		}, pulumi.Provider(bProvider))
		if err != nil {
			return err
		}

		// null opts
		_, err = kustomize.NewDirectory(ctx, "kustomize-nullopts", kustomize.DirectoryArgs{
			Directory:      pulumi.String("./testdata/options/kustomize"),
			ResourcePrefix: "kustomize-nullopts",
		})
		if err != nil {
			return err
		}

		// --- Chart ---
		_, err = helm.NewChart(ctx, "chart-options", helm.ChartArgs{
			Path:            pulumi.String("./testdata/options/chart"),
			ResourcePrefix:  "chart-options",
			Transformations: []yaml.Transformation{transformK8s},
			SkipAwait:       pulumi.Bool(true),
		},
			pulumi.Providers(aProvider),
			pulumi.Aliases([]pulumi.Alias{{Name: pulumi.String("chart-options-old")}}),
			pulumi.IgnoreChanges([]string{"ignored"}),
			pulumi.Protect(true),
			pulumi.DependsOn([]pulumi.Resource{sleep}),
			pulumi.Transformations([]pulumi.ResourceTransformation{applyAlias}),
			pulumi.Version("1.2.3"),
			pulumi.PluginDownloadURL("https://a.pulumi.test"),
		)
		if err != nil {
			return err
		}

		// "provider" option
		_, err = helm.NewChart(ctx, "chart-provider", helm.ChartArgs{
			Path:           pulumi.String("./testdata/options/chart"),
			ResourcePrefix: "chart-provider",
		}, pulumi.Provider(bProvider))
		if err != nil {
			return err
		}

		// null opts
		_, err = helm.NewChart(ctx, "chart-nullopts", helm.ChartArgs{
			Path:           pulumi.String("./testdata/options/chart"),
			ResourcePrefix: "chart-nullopts",
		})
		if err != nil {
			return err
		}

		return nil
	})
}
