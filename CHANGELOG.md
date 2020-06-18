## HEAD (Unreleased)

### Improvements

-   Fix prometheus-operator test to wait for the CRD to be ready before use (https://github.com/pulumi/pulumi-kubernetes/pull/1172)

## 2.3.1 (June 17, 2020)

### Improvements

-   Update resource deprecation/removal warnings. (https://github.com/pulumi/pulumi-kubernetes/pull/1162)

### Bug Fixes

-   Fix regression in TypeScript YAML SDK (https://github.com/pulumi/pulumi-kubernetes/pull/1157)

## 2.3.0 (June 9, 2020)

### Improvements

- NodeJS SDK updated to align with other Pulumi NodeJS SDKs. (https://github.com/pulumi/pulumi-kubernetes/pull/1151)
- .NET SDK updated to align with other Pulumi .NET SDKs. (https://github.com/pulumi/pulumi-kubernetes/pull/1132)
    - Deprecated resources are now marked as `Obsolete`.
    - Many classes are moved to new locations on disk while preserving the public namespaces and API.
    - Several unused argument/output classes were removed without any impact on resources (e.g. `DeploymentRollbackArgs`).
    - Fixed the type of some properties in `JSONSchemaPropsArgs` (there's no need to have 2nd-level inputs there):
        - `InputList<InputJson>` -> `InputList<JsonElement>`
        - `InputMap<Union<TArgs, InputList<string>>>` -> `InputMap<Union<TArgs, ImmutableArray<string>>>`
        
### Bug Fixes

-   Fix incorrect schema consts for apiVersion and kind (https://github.com/pulumi/pulumi-kubernetes/pull/1153)

## 2.2.2 (May 27, 2020)

-   2.2.1 SDK release process failed, so pushing a new tag.

## 2.2.1 (May 27, 2020)

### Improvements

-   Update deprecated/removed resource warnings. (https://github.com/pulumi/pulumi-kubernetes/pull/1135)
-   Update to client-go 1.18. (https://github.com/pulumi/pulumi-kubernetes/pull/1136)
-   Don't replace Service on .spec.type change. (https://github.com/pulumi/pulumi-kubernetes/pull/1139)

### Bug Fixes

-   Fix regex in python `include-crds` logic (https://github.com/pulumi/pulumi-kubernetes/pull/1145)

## 2.2.0 (May 15, 2020)

### Improvements

-   Support helm v3 `include-crds` argument. (https://github.com/pulumi/pulumi-kubernetes/pull/1102)
-   Bump python requests version dependency. (https://github.com/pulumi/pulumi-kubernetes/pull/1121)
-   Add apiextensions.CustomResource to Go SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/1125)

## 2.1.1 (May 8, 2020)

-   Python and .NET packages failed to publish for 2.1.0, so bumping release version.

## 2.1.0 (May 8, 2020)

### Improvements

-   Add YAML support to Go SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/1093).
-   Add Helm support to Go SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/1105).

### Bug Fixes

-   fix(customresources): use a 3-way merge patch instead of strategic merge. (https://github.com/pulumi/pulumi-kubernetes/pull/1095)
-   Fix required input props in Go SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/1090)
-   Update Go SDK using latest codegen packages. (https://github.com/pulumi/pulumi-kubernetes/pull/1089)
-   Fix schema type for Fields and RawExtension. (https://github.com/pulumi/pulumi-kubernetes/pull/1086)
-   Fix error parsing YAML in python 3.8 (https://github.com/pulumi/pulumi-kubernetes/pull/1079)
-   Fix HELM_HOME handling for Helm v3. (https://github.com/pulumi/pulumi-kubernetes/pull/1076)

## 2.0.0 (April 16, 2020)

### Improvements

-   Add consts to Go SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/1062).
-   Add `CustomResource` to .NET SDK (https://github.com/pulumi/pulumi-kubernetes/pull/1067).
-   Upgrade to Pulumi v2.0.0

### Bug fixes

-   Sort fetched helm charts into alphabetical order. (https://github.com/pulumi/pulumi-kubernetes/pull/1064)

## 1.6.0 (March 25, 2020)

### Improvements

-   Add a Go SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/1029) (https://github.com/pulumi/pulumi-kubernetes/pull/1042).
-   Add support for Kubernetes 1.18. (https://github.com/pulumi/pulumi-kubernetes/pull/872) (https://github.com/pulumi/pulumi-kubernetes/pull/1042).

### Bug fixes

-   Update the Python `Provider` class to use parameter naming consistent with other resources. (https://github.com/pulumi/pulumi-kubernetes/pull/1039).
-   Change URN for apiregistration resources. (https://github.com/pulumi/pulumi-kubernetes/pull/1021).

## 1.5.8 (March 16, 2020)

### Improvements

-   Automatically populate type aliases and additional secret outputs in the .NET SDK.  (https://github.com/pulumi/pulumi-kubernetes/pull/1026).
-   Update to Pulumi NuGet 1.12.1 and .NET Core 3.1.  (https://github.com/pulumi/pulumi-kubernetes/pull/1030).

## 1.5.7 (March 10, 2020)

### Bug fixes

-   Change URN for apiregistration resources. (https://github.com/pulumi/pulumi-kubernetes/pull/1021).
-   Replace PersistentVolume if volume source changes. (https://github.com/pulumi/pulumi-kubernetes/pull/1015).
-   Fix bool Python provider opts. (https://github.com/pulumi/pulumi-kubernetes/pull/1027).

## 1.5.6 (February 28, 2020)

### Bug fixes

-   Replace Daemonset if .spec.selector changes. (https://github.com/pulumi/pulumi-kubernetes/pull/1008).
-   Display error when pulumi plugin install fails. (https://github.com/pulumi/pulumi-kubernetes/pull/1010).

## 1.5.5 (February 25, 2020)

### Bug fixes

-   Upgrade pulumi/pulumi dep to 1.11.0 (fixes #984). (https://github.com/pulumi/pulumi-kubernetes/pull/1005).

## 1.5.4 (February 19, 2020)

### Improvements

-   Auto-generate aliases for all resource kinds. (https://github.com/pulumi/pulumi-kubernetes/pull/991).

### Bug fixes

-   Fix aliases for several resource kinds. (https://github.com/pulumi/pulumi-kubernetes/pull/990).
-   Don't require valid cluster for YAML render mode. (https://github.com/pulumi/pulumi-kubernetes/pull/997).
-   Fix .NET resources with empty arguments. (https://github.com/pulumi/pulumi-kubernetes/pull/983).
-   Fix panic condition in Pod await logic. (https://github.com/pulumi/pulumi-kubernetes/pull/998).

### Improvements

-   .NET SDK supports resources to work with YAML Kubernetes files and Helm charts.
(https://github.com/pulumi/pulumi-kubernetes/pull/980).

## 1.5.3 (February 11, 2020)

### Bug fixes

-   Change invoke call to always use latest version. (https://github.com/pulumi/pulumi-kubernetes/pull/987).

## 1.5.2 (February 10, 2020)

### Improvements

-   Optionally render YAML for k8s resources. (https://github.com/pulumi/pulumi-kubernetes/pull/936).

## 1.5.1 (February 7, 2020)

### Bug fixes

-   Specify provider version for invokes. (https://github.com/pulumi/pulumi-kubernetes/pull/982).

## 1.5.0 (February 4, 2020)

### Improvements

-   Update nodejs SDK to use optional chaining in constructor. (https://github.com/pulumi/pulumi-kubernetes/pull/959).
-   Automatically set Secret inputs as pulumi.secret. (https://github.com/pulumi/pulumi-kubernetes/pull/961).
-   Create helm.v3 alias. (https://github.com/pulumi/pulumi-kubernetes/pull/970).

### Bug fixes

-   Fix hang on large YAML files. (https://github.com/pulumi/pulumi-kubernetes/pull/974).
-   Use resourcePrefix all code paths. (https://github.com/pulumi/pulumi-kubernetes/pull/977).

## 1.4.5 (January 22, 2020)

### Bug fixes

-   Handle invalid kubeconfig context. (https://github.com/pulumi/pulumi-kubernetes/pull/960).

## 1.4.4 (January 21, 2020)

### Improvements

-   Improve namespaced Kind check. (https://github.com/pulumi/pulumi-kubernetes/pull/947).
-   Add helm template `apiVersions` flag. (https://github.com/pulumi/pulumi-kubernetes/pull/894)
-   Move YAML decode logic into provider and improve handling of default namespaces for Helm charts. (https://github.com
/pulumi/pulumi-kubernetes/pull/952).

### Bug fixes

-   Gracefully handle unreachable k8s cluster. (https://github.com/pulumi/pulumi-kubernetes/pull/946).
-   Fix deprecation notice for CSINode. (https://github.com/pulumi/pulumi-kubernetes/pull/944).

## 1.4.3 (January 8, 2020)

### Bug fixes

-   Revert invoke changes. (https://github.com/pulumi/pulumi-kubernetes/pull/941).

## 1.4.2 (January 7, 2020)

### Improvements

-   Move YAML decode logic into provider. (https://github.com/pulumi/pulumi-kubernetes/pull/925).
-   Improve handling of default namespaces for Helm charts. (https://github.com/pulumi/pulumi-kubernetes/pull/934).

### Bug fixes

-   Fix panic condition in Ingress await logic. (https://github.com/pulumi/pulumi-kubernetes/pull/928).
-   Fix deprecation warnings and docs. (https://github.com/pulumi/pulumi-kubernetes/pull/929).
-   Fix projection of array-valued output properties in .NET. (https://github.com/pulumi/pulumi-kubernetes/pull/931)

## 1.4.1 (December 17, 2019)

### Bug fixes

-   Fix deprecation warnings and docs. (https://github.com/pulumi/pulumi-kubernetes/pull/918 and https://github.com /pulumi/pulumi-kubernetes/pull/921).

## 1.4.0 (December 9, 2019)

### Important

The discovery.v1alpha1.EndpointSlice and discovery.v1alpha1.EndpointSliceList APIs were removed in k8s 1.17,
and no longer appear in the Pulumi Kubernetes SDKs. These resources can now be found at
discovery.v1beta1.EndpointSlice and discovery.v1beta1.EndpointSliceList.

### Major changes

-   Add support for Kubernetes v1.17.0 (https://github.com/pulumi/pulumi-kubernetes/pull/706)

## 1.3.4 (December 5, 2019)

### Improvements

-   Use HELM_HOME as default if set. (https://github.com/pulumi/pulumi-kubernetes/pull/855).
-   Use `namespace` provided by `KUBECONFIG`, if it is not explicitly set in the provider (https://github.com/pulumi/pulumi-kubernetes/pull/903).

## 1.3.3 (November 29, 2019)

### Improvements

-   Add `Provider` for .NET. (https://github.com/pulumi/pulumi-kubernetes/pull/897)

## 1.3.2 (November 26, 2019)

### Improvements

-   Add support for .NET. (https://github.com/pulumi/pulumi-kubernetes/pull/885)

## 1.3.1 (November 18, 2019)

### Improvements

-   Add support for helm 3 CLI tool. (https://github.com/pulumi/pulumi-kubernetes/pull/882).

## 1.3.0 (November 13, 2019)

### Improvements

-   Increase maxBuffer for helm template exec. (https://github.com/pulumi/pulumi-kubernetes/pull/864).
-   Add StreamInvoke RPC call, along with stream invoke implementations for
    kubernetes:kubernetes:watch, kubernetes:kubernetes:list, and kubernetes:kubernetes:logs. (#858, #873, #876).

## 1.2.3 (October 17, 2019)

### Bug fixes

-   Correctly merge provided opts for k8s resources. (https://github.com/pulumi/pulumi-kubernetes/pull/850).
-   Fix a bug that causes helm crash when referencing 'scoped packages' that start with '@'. (https://github.com/pulumi/pulumi-kubernetes/pull/846)

## 1.2.2 (October 10, 2019)

### Improvements

-   Stop using initialApiVersion annotation. (https://github.com/pulumi/pulumi-kubernetes/pull/837).
-   Cache the parsed OpenAPI schema to improve performance. (https://github.com/pulumi/pulumi-kubernetes/pull/836).

## 1.2.1 (October 8, 2019)

### Improvements

-   Cache the OpenAPI schema to improve performance. (https://github.com/pulumi/pulumi-kubernetes/pull/833).
-   Aggregate error messages from Pods on Job Read. (https://github.com/pulumi/pulumi-kubernetes/pull/831).
-   Improve interactive status for Jobs. (https://github.com/pulumi/pulumi-kubernetes/pull/832).

## 1.2.0 (October 4, 2019)

### Improvements

-   Add logic to check for Job readiness. (https://github.com/pulumi/pulumi-kubernetes/pull/633).
-   Automatically mark Secret data and stringData as secret. (https://github.com/pulumi/pulumi-kubernetes/pull/803).
-   Auto-alias resource apiVersions. (https://github.com/pulumi/pulumi-kubernetes/pull/798).
-   Provide detailed error for removed apiVersions. (https://github.com/pulumi/pulumi-kubernetes/pull/809).

## 1.1.0 (September 18, 2019)

### Major changes

-   Add support for Kubernetes v1.16.0 (https://github.com/pulumi/pulumi-kubernetes/pull/669)

### Improvements

-   Implement customTimeout for resource deletion. (https://github.com/pulumi/pulumi-kubernetes/pull/802).
-   Increase default readiness timeouts to 10 mins. (https://github.com/pulumi/pulumi-kubernetes/pull/721).
-   Add suppressDeprecationWarnings flag. (https://github.com/pulumi/pulumi-kubernetes/pull/808).
-   Warn for invalid usage of Helm repo parameter. (https://github.com/pulumi/pulumi-kubernetes/pull/805).
-   Add PodAggregator for use by resource awaiters. (https://github.com/pulumi/pulumi-kubernetes/pull/785).

## 1.0.1 (September 11, 2019)

### Improvements

-   Warn for deprecated apiVersions.
    (https://github.com/pulumi/pulumi-kubernetes/pull/779).

### Bug fixes

-   Fix await logic for extensions/v1beta1/Deployment
    (https://github.com/pulumi/pulumi-kubernetes/pull/794).
-   Fix error reporting
    (https://github.com/pulumi/pulumi-kubernetes/pull/782).

## 1.0.0 (September 3, 2019)

### Bug fixes

-   Fix name collisions in the Charts/YAML Python packages
    (https://github.com/pulumi/pulumi-kubernetes/pull/771).
-   Implement `{ConfigFile, ConfigGroup, Chart}#get_resource`
    (https://github.com/pulumi/pulumi-kubernetes/pull/771).
-   Upgrade Pulumi dependency to 1.0.0.

## 1.0.0-rc.1 (August 28, 2019)

### Improvements

### Bug fixes

-   Do not leak unencrypted secret values into the state file (fixes https://github.com/pulumi/pulumi-kubernetes/issues/734).

## 1.0.0-beta.2 (August 26, 2019)

### Improvements

-   Refactor and update the docs of the repo for 1.0. (https://github.com/pulumi/pulumi-kubernetes/pull/736).
-   Document await logic in the SDKs. (https://github.com/pulumi/pulumi-kubernetes/pull/711).
-   Document await timeouts and how to override. (https://github.com/pulumi/pulumi-kubernetes/pull/718).
-   Improve CustomResource for Python SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/700).
-   Clean up Python SDK get methods. (https://github.com/pulumi/pulumi-kubernetes/pull/740).
-   Remove undocumented kubectl replace invoke method. (https://github.com/pulumi/pulumi-kubernetes/pull/738).
-   Don't populate `.status` in input types (https://github.com/pulumi/pulumi-kubernetes/pull/635).
-   Allow a user to pass CustomTimeouts as part of ResourceOptions (fixes https://github.com/pulumi/pulumi-kubernetes/issues/672)
-   Don't panic when an Asset or an Archive are passed into a resource definition (https://github.com/pulumi/pulumi-kubernetes/pull/751).

### Bug fixes

-   Fix error messages for resources with default namespace. (https://github.com/pulumi/pulumi-kubernetes/pull/749).
-   Correctly compute version number for plugin to send with registration requests (fixes https://github.com/pulumi/pulumi-kubernetes/issues/732).

## 1.0.0-beta.1 (August 13, 2019)

### Improvements

-   Add .get() to Python SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/435).

## 0.25.6 (August 7, 2019)

### Bug fixes

-   Align YAML parsing with core Kubernetes supported YAML subset. (https://github.com/pulumi/pulumi-kubernetes/pull/690).
-   Handle string values in the equalNumbers function. (https://github.com/pulumi/pulumi-kubernetes/pull/691).
-   Properly detect readiness for Deployment scaled to 0. (https://github.com/pulumi/pulumi-kubernetes/pull/688).
-   Fix a bug that caused crashes when empty array values were added to resource inputs. (https://github.com/pulumi/pulumi-kubernetes/pull/696)

## 0.25.5 (August 2, 2019)

### Bug fixes

-   Fall back to client-side diff if server-side diff fails. (https://github.com/pulumi/pulumi-kubernetes/pull/685).
-   Fix namespace arg for Python Helm SDK (https://github.com/pulumi/pulumi-kubernetes/pull/670).
-   Detect namespace diff for first-class providers. (https://github.com/pulumi/pulumi-kubernetes/pull/674).
-   Fix values arg for Python Helm SDK (https://github.com/pulumi/pulumi-kubernetes/pull/678).
-   Fix Python Helm LocalChartOpts to inherit from BaseChartOpts (https://github.com/pulumi/pulumi-kubernetes/pull/681).

## 0.25.4 (August 1, 2019)

### Important

This release reverts the default diff behavior back to the pre-`0.25.3` behavior. A new flag has
been added to the provider options called `enableDryRun`, that can be used to opt in to the new
diff behavior. This will eventually become the default behavior after further testing to ensure
that this change is not disruptive.

### Major changes

-   Disable dryRun diff behavior by default. (https://github.com/pulumi/pulumi-kubernetes/pull/686)

### Improvements

-   Improve error messages for StatefulSet. (https://github.com/pulumi/pulumi-kubernetes/pull/673)

### Bug fixes

-   Properly reference override values in Python Helm SDK (https://github.com/pulumi/pulumi-kubernetes/pull/676).
-   Handle Output values in diffs. (https://github.com/pulumi/pulumi-kubernetes/pull/682).

## 0.25.3 (July 29, 2019)

### Bug fixes

-   Allow `yaml.ConfigGroup` to take URLs as argument
    (https://github.com/pulumi/pulumi-kubernetes/pull/638).
-   Return useful errors when we fail to fetch URL YAML
    (https://github.com/pulumi/pulumi-kubernetes/pull/638).
-   Use JSON_SCHEMA when parsing Kubernetes YAML, to conform with the expectations of the Kubernetes
    core resource types. (https://github.com/pulumi/pulumi-kubernetes/pull/638).
-   Don't render emoji on Windows. (https://github.com/pulumi/pulumi-kubernetes/pull/634)
-   Emit a useful error message (rather than a useless one) if we fail to parse the YAML data in
    `kubernetes:config:kubeconfig` (https://github.com/pulumi/pulumi-kubernetes/pull/636).
-   Provide useful contexts in provider errors, particularly those that originate from the API
    server (https://github.com/pulumi/pulumi-kubernetes/pull/636).
-   Expose all Kubernetes types through the SDK
    (https://github.com/pulumi/pulumi-kubernetes/pull/637).
-   Use `opts` instead of `__opts__` and `resource_name` instead of `__name__` in Python SDK
    (https://github.com/pulumi/pulumi-kubernetes/pull/639).
-   Properly detect failed Deployment on rollout. (https://github.com/pulumi/pulumi-kubernetes/pull/646
    and https://github.com/pulumi/pulumi-kubernetes/pull/657).
-   Use dry-run support if available when diffing the actual and desired state of a resource
    (https://github.com/pulumi/pulumi-kubernetes/pull/649)
-   Fix panic when `.metadata.label` is mistyped
    (https://github.com/pulumi/pulumi-kubernetes/pull/655).
-   Fix unexpected diffs when running against an API server that does not support dry-run.
    (https://github.com/pulumi/pulumi-kubernetes/pull/658)

## 0.25.2 (July 11, 2019)

### Improvements

-   The Kubernetes provider can now communicate detailed information about the difference between a resource's
desired and actual state during a Pulumi update. (https://github.com/pulumi/pulumi-kubernetes/pull/618).
-   Refactor Pod await logic for easier testing and maintenance (https://github.com/pulumi/pulumi-kubernetes/pull/590).
-   Update to client-go v12.0.0 (https://github.com/pulumi/pulumi-kubernetes/pull/621).
-   Fallback to JSON merge if strategic merge fails (https://github.com/pulumi/pulumi-kubernetes/pull/622).

### Bug fixes

-   Fix Helm Chart resource by passing `resourcePrefix` to the yaml template resources (https://github.com/pulumi/pulumi-kubernetes/pull/625).

## 0.25.1 (July 2, 2019)

### Improvements

-   Unify diff behavior between `Diff` and `Update`. This should result in better detection of state drift as
    well as behavior that is more consistent with respect to `kubectl`. (https://github.com/pulumi/pulumi-kubernetes/pull/604)
-   The Kubernetes provider now supports the internal features necessary for the Pulumi engine to detect diffs between the actual and desired state of a resource after a `pulumi refresh` (https://github.com/pulumi/pulumi-kubernetes/pull/477).
-   The Kubernetes provider now sets the `"kubectl.kubernetes.io/last-applied-configuration"` annotation to the last deployed configuration for a resource. This enables better interoperability with `kubectl`.

### Bug fixes

-   Add more props that force replacement of Pods (https://github.com/pulumi/pulumi-kubernetes/pull/613)

## 0.25.0 (June 19, 2019)

### Major changes

-   Add support for Kubernetes v1.15.0 (https://github.com/pulumi/pulumi-kubernetes/pull/557)

### Improvements

-   Enable multiple instances of Helm charts per stack (https://github.com/pulumi/pulumi-kubernetes/pull/599).
-   Enable multiple instances of YAML manifests per stack (https://github.com/pulumi/pulumi-kubernetes/pull/594).

### Bug fixes

-   None

## 0.24.0 (June 5, 2019)

### Important

BREAKING: This release changes the behavior of the provider `namespace` flag introduced
in `0.23.0`. Previously, this flag was treated as an override, which ignored namespace
values set directly on resources. Now, the flag is a default, and will only set the
namespace if one is not already set. If you have created resources using a provider
with the `namespace` flag set, this change may cause these resources to be recreated
on the next update.

### Major changes

-   BREAKING: Change the recently added `transformations` callback in Python to match JavaScript API (https://github.com/pulumi/pulumi-kubernetes/pull/575)
-   BREAKING: Remove `getInputs` from Kubernetes resource implementations. (https://github.com/pulumi/pulumi-kubernetes/pull/580)
-   BREAKING: Change provider namespace from override to default. (https://github.com/pulumi/pulumi-kubernetes/pull/585)

### Improvements

-   Enable configuring `ResourceOptions` via `transformations` (https://github.com/pulumi/pulumi-kubernetes/pull/575).
-   Changing k8s cluster config now correctly causes dependent resources to be replaced (https://github.com/pulumi/pulumi-kubernetes/pull/577).
-   Add user-defined type guard `isInstance` to all Kubernetes `CustomResource` implementations (https://github.com/pulumi/pulumi-kubernetes/pull/582).

### Bug fixes

-   Fix panics during preview when `metadata` is a computed value (https://github.com/pulumi/pulumi-kubernetes/pull/572)

## 0.23.1 (May 10, 2019)

### Major changes

-   None

### Improvements

-   Update to use client-go v11.0.0 (https://github.com/pulumi/pulumi-kubernetes/pull/549)
-   Deduplicate provider logs (https://github.com/pulumi/pulumi-kubernetes/pull/558)

### Bug fixes

-   Fix namespaceable check for diff (https://github.com/pulumi/pulumi-kubernetes/pull/554)

## 0.23.0 (April 30, 2019)

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

### Major changes

-   Add support for Kubernetes v1.14.0 (https://github.com/pulumi/pulumi-kubernetes/pull/371)

### Improvements

-   Add CustomResource to Python SDK (https://github.com/pulumi/pulumi-kubernetes/pull/543)

### Bug fixes

-   None

## 0.21.1 (March 18, 2019)

### Major changes

-   None

### Improvements

-   Split up nodejs SDK into multiple files (https://github.com/pulumi/pulumi-kubernetes/pull/480)

### Bug fixes

-   Check for unexpected RPC ID and return an error (https://github.com/pulumi/pulumi-kubernetes/pull/475)
-   Fix an issue where the Python `pulumi_kubernetes` package was depending on an older `pulumi` package.
-   Fix YAML parsing for computed namespaces (https://github.com/pulumi/pulumi-kubernetes/pull/483)

## 0.21.0 (Released March 6, 2019)

### Important

Updating to v0.17.0 version of `@pulumi/pulumi`.  This is an update that will not play nicely
in side-by-side applications that pull in prior versions of this package.

See https://github.com/pulumi/pulumi/commit/7f5e089f043a70c02f7e03600d6404ff0e27cc9d for more details.

As such, we are rev'ing the minor version of the package from 0.16 to 0.17.  Recent version of `pulumi` will now detect, and warn, if different versions of `@pulumi/pulumi` are loaded into the same application.  If you encounter this warning, it is recommended you move to versions of the `@pulumi/...` packages that are compatible.  i.e. keep everything on 0.16.x until you are ready to move everything to 0.17.x.

## 0.20.4 (March 1, 2019)

### Major changes

-   None

### Improvements

-   Allow the default timeout for awaiters to be overridden (https://github.com/pulumi/pulumi-kubernetes/pull/457)

### Bug fixes

-   Properly handle computed values in labels and annotations (https://github.com/pulumi/pulumi-kubernetes/pull/461)

## 0.20.3 (February 20, 2019)

### Major changes

-   None

### Improvements

-   None

### Bug fixes

-   Move mocha dependencies to devDependencies (https://github.com/pulumi/pulumi-kubernetes/pull/441)
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
