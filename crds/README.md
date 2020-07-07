# Generating strongly-typed Kubernetes CRDs (Custom Resource Definitions)

## Usage

If you have some Kubernetes CRD properly defined in some file `resourcedefinition.yaml` and you want to generate a code snippet called `output.ts` that has the
schema strongly typed, you can run:

```bash
$ go run main.go gen.go parse.go nestedMap.go <resourcedefinition.yaml> > <output.ts>
```

## Example

If we wanted to manually create a CronTab CRD from the `resourcedefinition.yaml`
specified in the [Kubernetes Docs](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/), we would write something like this:

```ts
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const cronTabDefinition = new k8s.apiextensions.v1.CustomResourceDefinition("cronTabDefinition",
    {
        "apiVersion": "apiextensions.k8s.io/v1",
        "kind": "CustomResourceDefinition",
        "metadata": {
            "name": "crontabs.stable.example.com"
        },
        "spec": {
            "group": "stable.example.com",
            "names": {
                "kind": "CronTab",
                "plural": "crontabs",
                "shortNames": [
                    "ct"
                ],
                "singular": "crontab"
            },
            "scope": "Namespaced",
            "versions": [
                {
                    "name": "v1",
                    "schema": {
                        "openAPIV3Schema": {
                            "properties": {
                                "spec": {
                                    "properties": {
                                        "cronSpec": {
                                            "type": "string"
                                        },
                                        "image": {
                                            "type": "string"
                                        },
                                        "replicas": {
                                            "type": "integer"
                                        }
                                    },
                                    "required": [
                                        "cronSpec"
                                    ],
                                    "type": "object"
                                }
                            },
                            "type": "object"
                        }
                    },
                    "served": true,
                    "storage": true
                }
            ]
        },
    }
)

const newCronObject = new k8s.apiextensions.CustomResource("newCronObject", 
    {
        apiVersion: "stable.example.com/v1",
        kind: "CronTab",
        metadata: {
            name: "my-new-cron-object"
        },
        spec: {
            cronSpec: "* * * * */5",
            image: "my-awesome-cron-image"
        }
    }
)

export const urn = newCronObject.urn
```

However with this tool, we can simply can simply run:
```bash
$ go run main.go gen.go parse.go nestedMap.go resourcedefinition.yaml > crontab.ts
```
Which will generate `crontab.ts` in our current directoy. This gives us the
`CronTabDefinition` and `CronTab` classes, which makes it much faster to
provision a CRD and create instances of it. Notice how we no longer need to
pass in `apiVersion` and `kind` values for the `CronTab` instance, since we
were able to infer those values from `resourcedefinition.yaml`.
This tool also generates `CronTabArgs` and `CronTabSpecArgs` interfaces, which gives us
the benefit of strongly typed inputs when creating instances of `CronTab`. 

```ts
import * as crontab from "./crontab"

const cronTabDefinition = new crontab.CronTabDefinition("cronTabDefinition")
const newCronObject = new crontab.CronTab("newCronObject", 
    {
        metadata: {
            name: "my-new-cron-object"
        },
        spec: {
            cronSpec: "* * * * */5",
            image: "my-awesome-cron-image"
        }
    }
)
export const urn = newCronObject.urn
```