# crd2pulumi

Generate strongly-typed CustomResources based on a Kubernetes CustomResourceDefinition (CRD).

## Goals

`crd2pulumi` is a CLI tool that generates strongly-typed CustomResources based on a Kubernetes CRD. CustomResourceDefinitions allow you to extend the Kubernetes API by defining your schemas for custom objects. Pulumi lets you create [CustomResources](https://www.pulumi.com/docs/reference/pkg/kubernetes/apiextensions/customresource/), but previously there was no strong-typing for these objects since every schema was, well, custom. This can be a massive headache for popular CRDs such as [cert-manager](https://github.com/jetstack/cert-manager/tree/master/deploy/crds) or [istio](https://github.com/istio/istio/tree/0321da58ca86fc786fb03a68afd29d082477e4f2/manifests/charts/base/crds), which contain thousands of lines of complex YAML schemas. By generating strongly-typed versions of CustomResources, crd2pulumi makes filling out their arguments more convenient because it lets you leverage existing IDE type checking and autocomplete features.

Currently, TypeScript and Go are supported, with Python and .NET coming soon.

## Building and Installation

If you wish to use `crd2pulumi` without developing the tool itself, you can use one of the [binary releases](https://github.com/pulumi/pulumi-kubernetes/releases/tag/crd2pulumi/v1.0.0) hosted on GitHub.

`crd2pulumi` uses Go modules to manage dependencies. If you want to develop `crd2pulumi` itself, you'll need to have Go installed in order to build. Once this prerequisite is installed, run the following to build the `crd2pulumi` binary and install it into `$GOPATH/bin`:

```bash
$ go build -o $GOPATH/bin/crd2pulumi main.go
```

Go should automatically handle pulling the dependencies for you.

If `$GOPATH/bin` is not on your path, you may want to move the `crd2pulumi` binary from `$GOPATH/bin` into a directory that is on your path.

## Usage

```bash
$ crd2pulumi <language> <crd file path> [output directory] [--force]
```

`<language>` is the target language to generate code for, so either `nodejs` or `go`.

`<crd file path>` is the path to the k8s CRD YAML file.

`[output directory]` is an optional path to the directory in which to output the
code to. If this field is not specified, then crd2pulumi will automatically output
to the same directory as the CRD YAML file.

`[--force]` is an optional flag to overwrite existing files if crd2pulumi would
write to directories or files that already exist. By default this is set to false.

## Examples

Let's use the example CronTab CRD specified in `resourcedefinition.yaml` from the [Kubernetes Documentation](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/). 

### Output to TypeScript

To generate a strongly-typed CronTab CustomResource in TypeScript, we can run this command:

```bash
$ crd2pulumi nodejs resourcedefinition.yaml
```

> By default, this will create the folder `./crontabs` in the same directory as `resourcedefinition.yaml`. Each versioned CustomResource will live under this folder, so since we had just specified `v1`, you'll see that a `./crontabs/v1` directory has been created. If we had defined a `v2` version in the CRD YAML file, you'd also see `./crontabs/v2`.

`./crontabs` contains two useful classes: `v1.CronTab` and `CronTabDefinition`. `v1.Crontab` is the typed CustomResource for CronTab, and `CronTabDefinition` is a helper class that provisions the CRD YAML schema in a single line. Now let's import the generated code into a Pulumi program that provisions the CRD and creates an instance of it.

```typescript
import * as crontabs from "./crontabs"
import * as pulumi from "@pulumi/pulumi"

const cronTabDefinition = new crontabs.CronTabDefinition("my-crontab-definition")

const myCronTab = new crontabs.v1.CronTab("my-new-cron-object",
{
    metadata: {
        name: "my-new-cron-object",
    },
    spec: {
        cronSpec: "* * * * */5",
        image: "my-awesome-cron-image",
    }
})

export const urn = myCronTab.urn;
```

As you can see, the `v1.CronTab` object is typed! For example, if you try to set
`cronSpec` to a non-string or add an extra field, your IDE should immediately warn you.

### Output to Go

Here's an example of the same program, but written in Go. To generate a strongly-typed CronTab CustomResource, we can run this command:

```bash
$ crd2pulumi go resourcedefinition.yaml
```

> Like with TypeScript, this will create a `./crontabs` folder. If we had multiple versions, they would be generated in `./crontabs/v2/crontabs.go`, `./crontabs/v3/crontabs.go`, etc...

Now we can access the `NewCronTab()` constructor. Create a `main.go` file with the following, swapping out the `go_pulumi` in `import crontabs "go_pulumi/v1"` with your own module's name.

```go
package main

import (
  crontabs "go_pulumi/v1"

  v1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apiextensions/v1"
  metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
  "github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Register the CronTab CRD.
		_, err := yaml.NewConfigFile(ctx, "my-crontab-definition",
			&yaml.ConfigFileArgs{
				File: "resourcedefinition.yaml",
			},
		)
		if err != nil {
      return err
		}
		// Instantiate a CronTab resource.
		cronTabInstance, err := crontabs.NewCronTab(ctx, "my-new-cron-object",
			&crontabs.CronTabArgs{
				Metadata: metav1.ObjectMetaArgs{
					Name: pulumi.StringPtr("my-new-cron-object"),
				},
				Spec: crontabs.CronTabSpecArgs{
					CronSpec: pulumi.StringPtr("* * * * */5"),
					Image:    pulumi.StringPtr("my-awesome-cron-image"),
				},
			},
		)
		if err != nil {
			return err
		}
		ctx.Export("urn", cronTabInstance.URN())
		return nil
	})
}
```

You can use the generated `crontabs.CronTabSpecArgs` struct to ensure your arguments are valid! This is a great improvement compared to filling out a `map[string]interface{}`, which was the previous procedure for Kubernetes `CustomResources`.

Now let's run the program and perform the update.

```bash
$ pulumi up
Previewing update (dev):
  Type                                                      Name                Plan
  pulumi:pulumi:Stack                                       examples-dev
 +   ├─ kubernetes:stable.example.com:CronTab                   my-new-cron-object  create
 +   └─ kubernetes:apiextensions.k8s.io:CustomResourceDefinition  my-crontab-definition  create
Resources:
  + 2 to create
  1 unchanged
Do you want to perform this update? yes
Updating (dev):
  Type                                                      Name                Status
  pulumi:pulumi:Stack                                       examples-dev
 +   ├─ kubernetes:stable.example.com:CronTab                   my-new-cron-object  created
 +   └─ kubernetes:apiextensions.k8s.io:CustomResourceDefinition  my-crontab-definition  created
Outputs:
  urn: "urn:pulumi:dev::examples::kubernetes:stable.example.com/v1:CronTab::my-new-cron-object"
Resources:
  + 2 created
  1 unchanged
Duration: 17s
Permalink: https://app.pulumi.com/albert-zhong/examples/dev/updates/4
```

It looks like both the CronTab definition and instance were both created! Finally, let's verify that they were created
by manually viewing the raw YAML data:

```bash
$ kubectl get ct -o yaml
```

```yaml
- apiVersion: stable.example.com/v1
  kind: CronTab
  metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"stable.example.com/v1","kind":"CronTab","metadata":{"labels":{"app.kubernetes.io/managed-by":"pulumi"},"name":"my-new-cron-object"},"spec":{"cronSpec":"* * * * */5","image":"my-awesome-cron-image"}}
  creationTimestamp: "2020-08-10T09:50:38Z"
  generation: 1
  labels:
    app.kubernetes.io/managed-by: pulumi
  name: my-new-cron-object
  namespace: default
  resourceVersion: "1658962"
  selfLink: /apis/stable.example.com/v1/namespaces/default/crontabs/my-new-cron-object
  uid: 5e2c56a2-7332-49cf-b0fc-211a0892c3d5
  spec:
  cronSpec: '* * * * */5'
  image: my-awesome-cron-image
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""
```
