## 0.20.3 (Unreleased)

### Major changes

-   None

### Improvements

-   None

### Bug fixes

-   Include managed-by label in diff preview (https://github.com/pulumi/pulumi-kubernetes/pull/431)

## 0.20.2 (Released February 13, 2019)

### Major changes

-   None

### Improvements

-   Allow awaiters to be skipped by setting an annotation (https://github.com/pulumi/pulumi-kubernetes/pull/417)
-   Set managed-by: pulumi label on all created resources (https://github.com/pulumi/pulumi-kubernetes/pull/418)
-   Clean up docstrings for Helm package (https://github.com/pulumi/pulumi-kubernetes/pull/396)
-   Support explicit `deleteBeforeReplace` (https://github.com/pulumi/pulumi/pull/2415)

### Bug fixes

-   Fix an issue with variable casing (https://github.com/pulumi/pulumi-kubernetes/pull/412)
-   Use modified copy of memcache client (https://github.com/pulumi/pulumi-kubernetes/pull/414)

## 0.20.1 (Released February 6, 2019)

### Bug fixes

-   Fix namespace handling regression (https://github.com/pulumi/pulumi-kubernetes/pull/403)
-   Nest Input<T> inside arrays (https://github.com/pulumi/pulumi-kubernetes/pull/395)

## 0.20.0 (Released February 1, 2019)

### Major changes

-   Add support for first-class Python providers (https://github.com/pulumi/pulumi-kubernetes/pull/350)
-   Upgrade to client-go 0.10.0 (https://github.com/pulumi/pulumi-kubernetes/pull/348)

### Improvements

-   Consider PVC events in Deployment await logic (https://github.com/pulumi/pulumi-kubernetes/pull/355)
-   Improve info message for Ingress with default path (https://github.com/pulumi/pulumi-kubernetes/pull/388)
-   Autogenerate Python casing table from OpenAPI spec (https://github.com/pulumi/pulumi-kubernetes/pull/387)

### Bug fixes

-   Use `node-fetch` rather than `got` to support Node 6 (https://github.com/pulumi/pulumi-kubernetes/pull/390)
-   Prevent orphaned resources on cancellation during delete (https://github.com/pulumi/pulumi-kubernetes/pull/368)
-   Handle buggy case for headless Service with no port (https://github.com/pulumi/pulumi-kubernetes/pull/366)


## 0.19.0 (Released January 15, 2019)

### Major changes

-   Implement incremental status updates for `StatefulSet`
    (https://github.com/pulumi/pulumi-kubernetes/pull/307)
-   Allow the `@pulumi/kubernetes` YAML API to understand arbitrary URLs
    (https://github.com/pulumi/pulumi-kubernetes/pull/328)
-   Add support for `.get` on CustomResources
    (https://github.com/pulumi/pulumi-kubernetes/pull/329)
-   Add support for `.get` for first-class providers
    (https://github.com/pulumi/pulumi-kubernetes/pull/340)

### Improvements

-   Fix Ingress await logic for ExternalName Services
    (https://github.com/pulumi/pulumi-kubernetes/pull/320)
-   Fix replacement logic for Job
    (https://github.com/pulumi/pulumi-kubernetes/pull/324 and https://github.com/pulumi/pulumi-kubernetes/pull/324)
-   Fix Cluster/RoleBinding replace semantics
    (https://github.com/pulumi/pulumi-kubernetes/pull/337)
-   Improve typing for `apiVersion` and `kind`
    (https://github.com/pulumi/pulumi-kubernetes/pull/341)

## 0.18.0 (Released December 4, 2018)

### Major changes

-   Allow Helm Charts to have `pulumi.Input` in their `values`
    (https://github.com/pulumi/pulumi-kubernetes/pull/241)

### Improvements

-   Retry REST calls to Kubernetes if they fail, greatly improving resiliance against resorce
    operation ordering problems.
-   Add support for creating CRDs and CRs in the same app
    (https://github.com/pulumi/pulumi-kubernetes/pull/271,
    https://github.com/pulumi/pulumi-kubernetes/pull/280)
-   Add incremental await for logic for `Ingress`
    (https://github.com/pulumi/pulumi-kubernetes/pull/283)
-   Allow users to specify a Chart's source any way they can do it from the CLI
    (https://github.com/pulumi/pulumi-kubernetes/pull/284)
-   "Fix" "bug" that cases Pulumi to crash if there is a duplicate key in a YAML template, to conform
    with Helm's behavior (https://github.com/pulumi/pulumi-kubernetes/pull/289)
-   Emit better error when the API server is unreachable
    (https://github.com/pulumi/pulumi-kubernetes/pull/291)
-   Add support for Kubernetes v0.12.\* (https://github.com/pulumi/pulumi-kubernetes/pull/293)
-   Fix bug that spuriously requires `.metadata.name` to be specified in Kubernetes list types
    (_e.g._, `v1/List`) (https://github.com/pulumi/pulumi-kubernetes/pull/294,
    https://github.com/pulumi/pulumi-kubernetes/pull/296)
-   Add Kubernetes v0.13.\* support (https://github.com/pulumi/pulumi-kubernetes/pull/306)
-   Improve error message when `Service` fails to initialized
    (https://github.com/pulumi/pulumi-kubernetes/pull/309)
-   Fix bug that causes us to erroneously report `Pod`'s owner
    (https://github.com/pulumi/pulumi-kubernetes/pull/311)
