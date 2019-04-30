## 0.23.1 (Unreleased)

## 0.23.0 (April 30, 2019)

### Supported Kubernetes versions

- v1.14.x
- v1.13.x
- v1.12.x

### Important

This release fixes a longstanding issue with the provider namespace flag. Previously, this
flag was erroneously ignored, but will now cause any resources using this provider to be
created in the specified namespace. **This may cause resources to be recreated!** Unset the
namespace parameter to avoid this behavior. Also note that this parameter takes precedence
over any namespace defined on the underlying resource.

The Python SDK now supports YAML manifests and Helm charts, including `CustomResourceDefinitions`
and `CustomResources`!

### Major changes

-   Put all resources in specified provider namespace (https://github.com/pulumi/pulumi-kubernetes/pull/538)
-   Add Helm support to Python SDK (https://github.com/pulumi/pulumi-kubernetes/pull/544)

### Bug fixes

-   Fix Helm repo quoting for Windows (https://github.com/pulumi/pulumi-kubernetes/pull/540)
-   Fix Python YAML SDK (https://github.com/pulumi/pulumi-kubernetes/pull/545)

## 0.22.2 (April 11, 2019)

### Supported Kubernetes versions

- v1.14.x
- v1.13.x
- v1.12.x

### Important

This release improves handling for CustomResources (CRs) and CustomResourceDefinitions (CRDs).
CRs without a matching CRD will now be considered deleted during `pulumi refresh`, and `pulumi destroy`
will not fail to delete a CR if the related CRD is missing.
See https://github.com/pulumi/pulumi-kubernetes/pull/530 for details.

### Major changes

-   None

### Improvements

-   Improve error handling for "no match found" errors (https://github.com/pulumi/pulumi-kubernetes/pull/530)

### Bug fixes

-   None

## 0.22.1 (April 9, 2019)

### Supported Kubernetes versions

- v1.14.x
- v1.13.x
- v1.12.x

### Major changes

-   Add basic YAML support to Python SDK (https://github.com/pulumi/pulumi-kubernetes/pull/499)
-   Add transforms to YAML support for Python SDK (https://github.com/pulumi/pulumi-kubernetes/pull/500)

### Improvements

-   Move helm module into a directory (https://github.com/pulumi/pulumi-kubernetes/pull/512)
-   Move yaml module into a directory (https://github.com/pulumi/pulumi-kubernetes/pull/513)

### Bug fixes

-   Fix Deployment await logic for old API schema (https://github.com/pulumi/pulumi-kubernetes/pull/523)
-   Replace PodDisruptionBudget if spec changes (https://github.com/pulumi/pulumi-kubernetes/pull/527)

## 0.22.0 (March 25, 2019)

### Supported Kubernetes versions

- v1.14.x
- v1.13.x
- v1.12.x

### Major changes

-   Add support for Kubernetes v1.14.0 (https://github.com/pulumi/pulumi-kubernetes/pull/371)

### Improvements

-   Add CustomResource to Python SDK (https://github.com/pulumi/pulumi-kubernetes/pull/543)

### Bug fixes

-   None

## 0.21.1 (March 18, 2019)

### Supported Kubernetes versions

- v1.13.x
- v1.12.x
- v1.11.x

### Major changes

-   None

### Improvements

-   Split up nodejs SDK into multiple files (https://github.com/pulumi/pulumi-kubernetes/pull/480)

### Bug fixes

-   Check for unexpected RPC ID and return an error (https://github.com/pulumi/pulumi-kubernetes/pull/475)
-   Fix an issue where the Python `pulumi_kubernetes` package was depending on an older `pulumi` package.
-   Fix YAML parsing for computed namespaces (https://github.com/pulumi/pulumi-kubernetes/pull/483)

## 0.21.0 (Released March 6, 2019)

### Supported Kubernetes versions

- v1.13.x
- v1.12.x
- v1.11.x

### Important

Updating to v0.17.0 version of `@pulumi/pulumi`.  This is an update that will not play nicely
in side-by-side applications that pull in prior versions of this package.

See https://github.com/pulumi/pulumi/commit/7f5e089f043a70c02f7e03600d6404ff0e27cc9d for more details.

As such, we are rev'ing the minor version of the package from 0.16 to 0.17.  Recent version of `pulumi` will now detect, and warn, if different versions of `@pulumi/pulumi` are loaded into the same application.  If you encounter this warning, it is recommended you move to versions of the `@pulumi/...` packages that are compatible.  i.e. keep everything on 0.16.x until you are ready to move everything to 0.17.x.

## 0.20.4 (March 1, 2019)

### Supported Kubernetes versions

- v1.13.x
- v1.12.x
- v1.11.x

### Major changes

-   None

### Improvements

-   Allow the default timeout for awaiters to be overridden (https://github.com/pulumi/pulumi-kubernetes/pull/457)

### Bug fixes

-   Properly handle computed values in labels and annotations (https://github.com/pulumi/pulumi-kubernetes/pull/461)

## 0.20.3 (February 20, 2019)

### Supported Kubernetes versions

- v1.13.x
- v1.12.x
- v1.11.x

### Major changes

-   None

### Improvements

-   None

### Bug fixes

-   Move mocha dependencies to devDependencies (https://github.com/pulumi/pulumi-kubernetes/pull/441)
-   Include managed-by label in diff preview (https://github.com/pulumi/pulumi-kubernetes/pull/431)

## 0.20.2 (Released February 13, 2019)

### Supported Kubernetes versions

- v1.13.x
- v1.12.x
- v1.11.x

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

### Supported Kubernetes versions

- v1.13.x
- v1.12.x
- v1.11.x

### Bug fixes

-   Fix namespace handling regression (https://github.com/pulumi/pulumi-kubernetes/pull/403)
-   Nest Input<T> inside arrays (https://github.com/pulumi/pulumi-kubernetes/pull/395)

## 0.20.0 (Released February 1, 2019)

### Supported Kubernetes versions

- v1.13.x
- v1.12.x
- v1.11.x

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

### Supported Kubernetes versions

- v1.13.x
- v1.12.x
- v1.11.x

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

### Supported Kubernetes versions

- v1.13.x
- v1.12.x
- v1.11.x

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
