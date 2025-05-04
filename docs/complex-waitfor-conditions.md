# Complex Conditions for `pulumi.com/waitFor`

The `pulumi.com/waitFor` annotation supports complex logical expressions to define custom readiness criteria for Kubernetes resources. This document explains the syntax and usage.

## Basic Syntax

The basic forms of `pulumi.com/waitFor` are:

1. Single condition (string): 
   ```
   "pulumi.com/waitFor": "condition=Ready"
   ```

2. JSONPath expression (string):
   ```
   "pulumi.com/waitFor": "jsonpath={.status.phase}=Running"
   ```

3. AND conditions (array of strings):
   ```
   "pulumi.com/waitFor": ["condition=Ready", "jsonpath={.status.phase}=Running"]
   ```

## New Complex Condition Syntax

The enhanced syntax allows for:

1. Explicit AND operator:
   ```json
   {
     "operator": "and",
     "conditions": [
       "condition=Ready",
       "jsonpath={.status.phase}=Running"
     ]
   }
   ```

2. OR operator:
   ```json
   {
     "operator": "or",
     "conditions": [
       "jsonpath={.status.phase}=Running",
       "jsonpath={.status.phase}=Succeeded"
     ]
   }
   ```

3. Nested expressions:
   ```json
   {
     "operator": "and",
     "conditions": [
       "condition=Ready",
       {
         "operator": "or",
         "conditions": [
           "jsonpath={.status.phase}=Running",
           "jsonpath={.status.phase}=Succeeded"
         ]
       }
     ]
   }
   ```

## Usage Example in Pulumi

Here's an example of how to use complex conditions in a Pulumi program:

```typescript
import * as k8s from "@pulumi/kubernetes";

const resource = new k8s.apiextensions.CustomResource(
  "my-resource",
  {
    apiVersion: "example.com/v1",
    kind: "MyResource",
    metadata: {
      annotations: {
        // This resource will be considered ready when:
        // (it has a Ready=True condition) AND 
        // (its status.phase is either Running OR Succeeded)
        "pulumi.com/waitFor": JSON.stringify({
          operator: "and",
          conditions: [
            "condition=Ready",
            {
              operator: "or",
              conditions: [
                "jsonpath={.status.phase}=Running",
                "jsonpath={.status.phase}=Succeeded"
              ]
            }
          ]
        })
      }
    },
    spec: {
      // Resource spec...
    }
  }
);
```

## Logical Evaluation

The conditions are evaluated as follows:

- For `"operator": "and"`, all conditions must be true
- For `"operator": "or"`, at least one condition must be true
- Conditions are evaluated recursively for nested expressions

This allows for arbitrarily complex logical expressions to define exactly when a resource should be considered ready.
