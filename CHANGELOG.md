## Unreleased

### Added

- `clusterIdentifier` configuration can now be used to manually control the
  replacement behavior of a provider resource.
  (https://github.com/pulumi/pulumi-kubernetes/pull/3068)
- Documentation is now generated for all languages supported by overlay types.
  (https://github.com/pulumi/pulumi-kubernetes/pull/3107)

### Fixed

- Updated logic to accurately detect if a resource is a Patch variant. (https://github.com/pulumi/pulumi-kubernetes/pull/3102)
- Added java as supported language for CustomResource overlays. (https://github.com/pulumi/pulumi-kubernetes/pull/3120)

## 4.15.0 (July 9, 2024)

### Changed

- `CustomResource` should have plain `apiVersion` and `kind` properties (https://github.com/pulumi/pulumi-kubernetes/pull/3079)

### Fixed

- Prevent CustomResourceDefinitions from always being applied to the cluster during preview operations (https://github.com/pulumi/pulumi-kubernetes/pull/3096)

## 4.14.0 (June 28, 2024)

### Added

- `TypedDict` input types for the Python SDK (https://github.com/pulumi/pulumi-kubernetes/pull/3070)

### Changed

- The `Release` resource no longer ignores empty lists when merging values. (https://github.com/pulumi/pulumi-kubernetes/pull/2995)

### Fixed

- `Chart` v4 now handles an array of assets. (https://github.com/pulumi/pulumi-kubernetes/pull/3061)
- Fix previews always failing when a resource is to be replaced (https://github.com/pulumi/pulumi-kubernetes/pull/3053)

## 4.13.1 (June 4, 2024)

### Added

- Kustomize Directory v2 resource (https://github.com/pulumi/pulumi-kubernetes/pull/3036) 
- CustomResource for Java SDK (https://github.com/pulumi/pulumi-kubernetes/pull/3020)

### Changed

- Update to pulumi-java v0.12.0 (https://github.com/pulumi/pulumi-kubernetes/pull/3025)

### Fixed

- Fixed Chart v4 fails on update (https://github.com/pulumi/pulumi-kubernetes/pull/3046)
- Fixed a panic that occurs when diffing Job resources containing `replaceUnready` annotations and an unreachable cluster connection. (https://github.com/pulumi/pulumi-kubernetes/pull/3024)
- Fixed spurious diffing for updates when in renderYaml mode (https://github.com/pulumi/pulumi-kubernetes/pull/3030)

## 4.12.0 (May 21, 2024)

### Added

- Added a new Helm Chart v4 resource. (https://github.com/pulumi/pulumi-kubernetes/pull/2947)
- Added support for deletion propagation policies (e.g. Orphan). (https://github.com/pulumi/pulumi-kubernetes/pull/3011)
- Server-side apply conflict errors now include the original field manager's name. (https://github.com/pulumi/pulumi-kubernetes/pull/2983)

### Changed 

- Pulumi will now wait for DaemonSets to become ready. (https://github.com/pulumi/pulumi-kubernetes/pull/2953)
- The Release resource's merge behavior for `valueYamlFiles` now more closely matches Helm's behavior. (https://github.com/pulumi/pulumi-kubernetes/pull/2963)

### Fixed

- Helm Chart V3 previews no longer fail when the cluster is unreachable. (https://github.com/pulumi/pulumi-kubernetes/pull/2992)
- Fixed a panic that could occur when a missing field became `null`. (https://github.com/pulumi/pulumi-kubernetes/issues/1970)

## 4.11.0 (April 17, 2024)

- [dotnet] Unknowns for previews involving an uninitialized provider (https://github.com/pulumi/pulumi-kubernetes/pull/2957)
- Update Kubernetes schemas and libraries to v1.30.0 (https://github.com/pulumi/pulumi-kubernetes/pull/2932)

## 4.10.0 (April 11, 2024)

- ConfigGroup V2 (https://github.com/pulumi/pulumi-kubernetes/pull/2844)
- ConfigFile V2 (https://github.com/pulumi/pulumi-kubernetes/pull/2862)
- Bugfix for ambiguous kinds (https://github.com/pulumi/pulumi-kubernetes/pull/2889)
- [yaml/v2] Support for resource ordering (https://github.com/pulumi/pulumi-kubernetes/pull/2894)
- Bugfix for deployment await logic not referencing the correct deployment status (https://github.com/pulumi/pulumi-kubernetes/pull/2943)

### New Features

A new MLC-based implementation of `ConfigGroup` and of `ConfigFile` is now available in the "yaml/v2" package. These resources are
usable in all Pulumi languages, including Pulumi YAML and in the Java Pulumi SDK.

Note that transformations aren't supported in this release (see https://github.com/pulumi/pulumi/issues/12996).

## 4.9.1 (March 13, 2024)

- Use async invokes to avoid hangs/stalls in Python `helm`, `kustomize`, and `yaml` components (https://github.com/pulumi/pulumi-kubernetes/pull/2863)

## 4.9.0 (March 4, 2024)

- Fix SSA ignoreChanges by enhancing field manager path comparisons (https://github.com/pulumi/pulumi-kubernetes/pull/2828)
- Update nodejs SDK dependencies (https://github.com/pulumi/pulumi-kubernetes/pull/2858, https://github.com/pulumi/pulumi-kubernetes/pull/2861)

## 4.8.1 (February 22, 2024)

- skip normalization in preview w/ computed fields (https://github.com/pulumi/pulumi-kubernetes/pull/2846)

## 4.8.0 (February 22, 2024)

- Fix DiffConfig issue when when provider's kubeconfig is set to file path (https://github.com/pulumi/pulumi-kubernetes/pull/2771)
- Fix for replacement having incorrect status messages (https://github.com/pulumi/pulumi-kubernetes/pull/2810)
- Use output properties for await logic (https://github.com/pulumi/pulumi-kubernetes/pull/2790)
- Support for metadata.generateName (CSA) (https://github.com/pulumi/pulumi-kubernetes/pull/2808)
- Fix unmarshalling of Helm values yaml file (https://github.com/pulumi/pulumi-kubernetes/issues/2815)
- Handle unknowns in Helm Release resource (https://github.com/pulumi/pulumi-kubernetes/pull/2822)

## 4.7.1 (January 17, 2024)

- Fix deployment await logic for accurate rollout detection

## 4.7.0 (January 17, 2024)
- Fix JSON encoding of KubeVersion and Version on Chart resource (.NET SDK) (https://github.com/pulumi/pulumi-kubernetes/pull/2740)
- Fix option propagation in component resources (Python SDK) (https://github.com/pulumi/pulumi-kubernetes/pull/2717)
- Fix option propagation in component resources (.NET SDK) (https://github.com/pulumi/pulumi-kubernetes/pull/2720)
- Fix option propagation in component resources (NodeJS SDK) (https://github.com/pulumi/pulumi-kubernetes/pull/2713)
- Fix option propagation in component resources (Go SDK) (https://github.com/pulumi/pulumi-kubernetes/pull/2709)

### Breaking Changes
In previous versions of the pulumi-kubernetes .NET SDK, the `ConfigFile` and `ConfigGroup` component resources inadvertently assigned the wrong parent to the child resource(s).
This would happen when the component resource itself had a parent; the child would be assigned that same parent. This also had the effect of disregarding the component resource's provider in favor of the parent's provider.

For example, here's a before/after look at the component hierarchy:

Before:

```
├─ pkg:index:MyComponent            parent
│  ├─ kubernetes:core/v1:ConfigMap  cg-options-cg-options-cm-1
│  ├─ kubernetes:yaml:ConfigFile    cg-options-testdata/options/configgroup/manifest.yaml
│  ├─ kubernetes:core/v1:ConfigMap  cg-options-configgroup-cm-1
│  ├─ kubernetes:yaml:ConfigFile    cg-options-testdata/options/configgroup/empty.yaml
│  └─ kubernetes:yaml:ConfigGroup   cg-options
```

After:

```
└─ pkg:index:MyComponent                  parent
   └─ kubernetes:yaml:ConfigGroup         cg-options
      ├─ kubernetes:yaml:ConfigFile       cg-options-testdata/options/configgroup/manifest.yaml
      │  └─ kubernetes:core/v1:ConfigMap  cg-options-configgroup-cm-1
      └─ kubernetes:core/v1:ConfigMap     cg-options-cg-options-cm-1
```

This release addresses this issue and attempts to heal existing stacks using aliases. This is effective at avoiding a replacement except in the case where the child was created with the wrong provider. In this case, __Pulumi will suggest a replacement of the child resource(s), such that they use the correct provider__.

## 4.6.1 (December 14, 2023)
- Fix: Refine URN lookup by using its core type for more accurate resource identification (https://github.com/pulumi/pulumi-kubernetes/issues/2719)

## 4.6.0 (December 13, 2023)
- Fix: Helm OCI chart deployment fails in Windows (https://github.com/pulumi/pulumi-kubernetes/pull/2648)
- Fix: compute version field in Check for content detection (https://github.com/pulumi/pulumi-kubernetes/pull/2672)
- Fix: Fix: Helm Release fails with "the server could not find the requested resource" (https://github.com/pulumi/pulumi-kubernetes/pull/2677)
- Fix Helm Chart resource lookup key handling for objects in default namespace (https://github.com/pulumi/pulumi-kubernetes/pull/2655)
- Update Kubernetes schemas and libraries to v1.29.0 (https://github.com/pulumi/pulumi-kubernetes/pull/2690)
- Fix panic when using `PULUMI_KUBERNETES_MANAGED_BY_LABEL` env var with SSA created objects (https://github.com/pulumi/pulumi-kubernetes/pull/2711)
- Fix normalization of base64 encoded secrets.data values to strip whitespace (https://github.com/pulumi/pulumi-kubernetes/issues/2715)

### Resources Renamed:
- `#/types/kubernetes:core/v1:ResourceRequirements`
  - renamed to: `#/types/kubernetes:core/v1:VolumeResourceRequirements`
- `#/types/kubernetes:core/v1:ResourceRequirementsPatch`
  - renamed to: `#/types/kubernetes:core/v1:VolumeResourceRequirementsPatch`

### New Resources:
- `flowcontrol.apiserver.k8s.io/v1.FlowSchema`
- `flowcontrol.apiserver.k8s.io/v1.FlowSchemaList`
- `flowcontrol.apiserver.k8s.io/v1.FlowSchemaPatch`
- `flowcontrol.apiserver.k8s.io/v1.PriorityLevelConfiguration`
- `flowcontrol.apiserver.k8s.io/v1.PriorityLevelConfigurationList`
- `flowcontrol.apiserver.k8s.io/v1.PriorityLevelConfigurationPatch`
- `networking.k8s.io/v1alpha1.ServiceCIDR`
- `networking.k8s.io/v1alpha1.ServiceCIDRList`
- `networking.k8s.io/v1alpha1.ServiceCIDRPatch`
- `storage.k8s.io/v1alpha1.VolumeAttributesClass`
- `storage.k8s.io/v1alpha1.VolumeAttributesClassList`
- `storage.k8s.io/v1alpha1.VolumeAttributesClassPatch`

## 4.5.5 (November 28, 2023)
- Fix: Make the invoke calls for Helm charts and YAML config resilient to the value being None or an empty dict (https://github.com/pulumi/pulumi-kubernetes/pull/2665)

## 4.5.4 (November 8, 2023)
- Fix: Helm Release: chart requires kubeVersion (https://github.com/pulumi/pulumi-kubernetes/pull/2653)

## 4.5.3 (October 31, 2023)
- Fix: Update pulumi version to 3.91.1 to pick up fixes in python codegen (https://github.com/pulumi/pulumi-kubernetes/pull/2647)

## 4.5.2 (October 26, 2023)
- Fix: Do not patch field managers for Patch resources (https://github.com/pulumi/pulumi-kubernetes/pull/2640)

## 4.5.1 (October 24, 2023)
- Revert: Normalize provider inputs and make available as outputs (https://github.com/pulumi/pulumi-kubernetes/pull/2627)

## 4.5.0 (October 23, 2023)

- helm.v3.ChartOpts: Add KubeVersion field that can be passed to avoid asking the kubernetes API server for the version (https://github.com/pulumi/pulumi-kubernetes/pull/2593)
- Fix for Helm Import regression (https://github.com/pulumi/pulumi-kubernetes/pull/2605)
- Improved search functionality for Helm Import (https://github.com/pulumi/pulumi-kubernetes/pull/2610)
- Fix SSA dry-run previews when a Pulumi program uses Apply on the status subresource (https://github.com/pulumi/pulumi-kubernetes/pull/2615)
- Normalize provider inputs and make available as outputs (https://github.com/pulumi/pulumi-kubernetes/pull/2598)

## 4.4.0 (October 12, 2023)

- Fix normalizing fields with empty objects/slices (https://github.com/pulumi/pulumi-kubernetes/pull/2576)
- helm.v3.Release: Improved cancellation support (https://github.com/pulumi/pulumi-kubernetes/pull/2579)
- Update Kubernetes client library to v0.28.2 (https://github.com/pulumi/pulumi-kubernetes/pull/2585)

## 4.3.0 (September 25, 2023)

- helm.v3.Release: Detect changes to local charts (https://github.com/pulumi/pulumi-kubernetes/pull/2568)
- Ignore read-only inputs in Kubernetes object metadata (https://github.com/pulumi/pulumi-kubernetes/pull/2571)
- Handle fields specified in ignoreChanges gracefully without needing a refresh when drift has occurred (https://github.com/pulumi/pulumi-kubernetes/pull/2566)

## 4.2.0 (September 14, 2023)

- Reintroduce switching builds to pyproject.toml; when publishing the package to PyPI both
  source-based and wheel distributions are now published. For most users the installs will now favor
  the wheel distribution, but users invoking pip with `--no-binary :all:` will continue having
  installs based on the source distribution.
- Return mapping information for terraform conversions (https://github.com/pulumi/pulumi-kubernetes/pull/2457)
- feature: added skipUpdateUnreachable flag to proceed with the updates without failing (https://github.com/pulumi/pulumi-kubernetes/pull/2528)

## 4.1.1 (August 23, 2023)

- Revert the switch to pyproject.toml and wheel-based PyPI publishing as it impacts users that run pip with --no-binary
  (see https://github.com/pulumi/pulumi-kubernetes/issues/2540)

## 4.1.0 (August 15, 2023)

- fix: ensure CSA does not hit API Server for preview (https://github.com/pulumi/pulumi-kubernetes/pull/2522)
- Fix helm.v3.Release replace behavior (https://github.com/pulumi/pulumi-kubernetes/pull/2532)
- [sdk/python] Switch to pyproject.toml and wheel-based PyPI publishing (https://github.com/pulumi/pulumi-kubernetes/pull/2493)
- Update Kubernetes to v1.28.0 (https://github.com/pulumi/pulumi-kubernetes/pull/2526)

## 4.0.3 (July 21, 2023)

- fix: ensure data is not dropped when normalizing Secrets (https://github.com/pulumi/pulumi-kubernetes/pull/2514)

## 4.0.2 (July 20, 2023)

- [sdk/python] Drop unused pyyaml dependency (https://github.com/pulumi/pulumi-kubernetes/pull/2502)
- Fix continuous diff for Secret stringData field (https://github.com/pulumi/pulumi-kubernetes/pull/2511)
- Fix diff for CRD with status set on input (https://github.com/pulumi/pulumi-kubernetes/pull/2512)

## 4.0.1 (July 19, 2023)

- Gracefully handle undefined resource schemes (https://github.com/pulumi/pulumi-kubernetes/pull/2504)
- Fix diff for CRD .spec.preserveUnknownFields (https://github.com/pulumi/pulumi-kubernetes/pull/2506)


## 4.0.0 (July 19, 2023)

Breaking changes:

- Remove deprecated enableDryRun provider flag (https://github.com/pulumi/pulumi-kubernetes/pull/2400)
- Remove deprecated helm/v2 SDK (https://github.com/pulumi/pulumi-kubernetes/pull/2396)
- Remove deprecated enableReplaceCRD provider flag (https://github.com/pulumi/pulumi-kubernetes/pull/2402)
- Drop support for Kubernetes clusters older than v1.13 (https://github.com/pulumi/pulumi-kubernetes/pull/2414)
- Make all resource output properties required (https://github.com/pulumi/pulumi-kubernetes/pull/2422)

Other changes:

- Enable Server-side Apply by default (https://github.com/pulumi/pulumi-kubernetes/pull/2398)
- Automatically fall back to client-side preview if server-side preview fails (https://github.com/pulumi/pulumi-kubernetes/pull/2419)
- Drop support for legacy pulumi.com/initialApiVersion annotation (https://github.com/pulumi/pulumi-kubernetes/pull/2443)
- Overhaul logic for resource diffing (https://github.com/pulumi/pulumi-kubernetes/pull/2445)
    - Drop usage of the "kubectl.kubernetes.io/last-applied-configuration" annotation.
    - Compute preview diffs using resource inputs rather than making a dry-run API call.
    - Automatically update .metadata.managedFields to work with resources that were managed with client-side apply, and later upgraded to use server-side apply.
    - Fix a bug with the diff calculation so that resource drift is detected accurately after a refresh.
- Update go module version to v4 (https://github.com/pulumi/pulumi-kubernetes/pull/2466)
- Upgrade to latest helm dependency (https://github.com/pulumi/pulumi-kubernetes/pull/2474)
- Improve error handling for List resources (https://github.com/pulumi/pulumi-kubernetes/pull/2493)

## 3.30.2 (July 11, 2023)

- Improve deleteUnreachable workflow for unreachable clusters (https://github.com/pulumi/pulumi-kubernetes/pull/2489)

## 3.30.1 (June 29, 2023)

- Add experimental helmChart support to kustomize.Directory (https://github.com/pulumi/pulumi-kubernetes/pull/2471)

## 3.30.0 (June 28, 2023)

- [sdk/python] Fix bug with class methods for YAML transformations (https://github.com/pulumi/pulumi-kubernetes/pull/2469)
- Fix StatefulSet await logic for OnDelete update (https://github.com/pulumi/pulumi-kubernetes/pull/2473)
- Skip wait for Pods on headless Service (https://github.com/pulumi/pulumi-kubernetes/pull/2475)

## 3.29.1 (June 14, 2023)

- Fix provider handling of CustomResources with Patch suffix (https://github.com/pulumi/pulumi-kubernetes/pull/2438)
- Improve status message for Deployment awaiter (https://github.com/pulumi/pulumi-kubernetes/pull/2456)

## 3.29.0 (June 2, 2023)

- Fix regression in file/folder checking logic that caused incorrect parsing of compressed chart files (https://github.com/pulumi/pulumi-kubernetes/pull/2428)
- Update Patch resources rather than replacing (https://github.com/pulumi/pulumi-kubernetes/pull/2429)

## 3.28.1 (May 24, 2023)

- Add a "strict mode" configuration option (https://github.com/pulumi/pulumi-kubernetes/pull/2425)

## 3.28.0 (May 19, 2023)

- Handle resource change from static name to autoname under SSA (https://github.com/pulumi/pulumi-kubernetes/pull/2392)
- Fix Helm release creation when the name of the chart conflicts with the name of a folder in the current working directory (https://github.com/pulumi/pulumi-kubernetes/pull/2410)
- Remove imperative authentication and authorization resources: TokenRequest, TokenReview, LocalSubjectAccessReview,
    SelfSubjectReview, SelfSubjectAccessReview, SelfSubjectRulesReview, and SubjectAccessReview (https://github.com/pulumi/pulumi-kubernetes/pull/2413)
- Improve check for existing resource GVK (https://github.com/pulumi/pulumi-kubernetes/pull/2418)

## 3.27.1 (May 11, 2023)

- Update Kubernetes client library to v0.27.1 (https://github.com/pulumi/pulumi-kubernetes/pull/2380)
- Increase default client burst and QPS to avoid throttling (https://github.com/pulumi/pulumi-kubernetes/pull/2381)
- Add HTTP request timeout option to KubeClientSettings (https://github.com/pulumi/pulumi-kubernetes/pull/2383)

## 3.27.0 (May 9, 2023)

- Change destroy operation to use foreground cascading delete (https://github.com/pulumi/pulumi-kubernetes/pull/2379)

## 3.26.0 (May 1, 2023)

- Do not await during .get or import operations (https://github.com/pulumi/pulumi-kubernetes/pull/2373)

## 3.25.0 (April 11, 2023)
- Update Kubernetes to v1.27.0

### New resources:

    authentication.k8s.io/v1beta1.SelfSubjectReview
    authentication.k8s.io/v1beta1.SelfSubjectReviewPatch
    certificates.k8s.io/v1alpha1.ClusterTrustBundle
    certificates.k8s.io/v1alpha1.ClusterTrustBundleList
    certificates.k8s.io/v1alpha1.ClusterTrustBundlePatch
    networking.k8s.io/v1alpha1.IPAddress
    networking.k8s.io/v1alpha1.IPAddressList
    networking.k8s.io/v1alpha1.IPAddressPatch
    resource.k8s.io/v1alpha2.PodSchedulingContext
    resource.k8s.io/v1alpha2.PodSchedulingContextList
    resource.k8s.io/v1alpha2.PodSchedulingContextPatch
    resource.k8s.io/v1alpha2.ResourceClaim
    resource.k8s.io/v1alpha2.ResourceClaimList
    resource.k8s.io/v1alpha2.ResourceClaimPatch
    resource.k8s.io/v1alpha2.ResourceClaimTemplate
    resource.k8s.io/v1alpha2.ResourceClaimTemplateList
    resource.k8s.io/v1alpha2.ResourceClaimTemplatePatch
    resource.k8s.io/v1alpha2.ResourceClass
    resource.k8s.io/v1alpha2.ResourceClassList
    resource.k8s.io/v1alpha2.ResourceClassPatch

### Resources moved from v1alpha1 to v1alpha2
- "kubernetes:resource.k8s.io/v1alpha1:ResourceClaimTemplateList"
- "kubernetes:resource.k8s.io/v1alpha1:ResourceClassList"
- "kubernetes:resource.k8s.io/v1alpha1:ResourceClassPatch"
- "kubernetes:resource.k8s.io/v1alpha1:ResourceClaimList"
- "kubernetes:resource.k8s.io/v1alpha1:ResourceClass"
- "kubernetes:resource.k8s.io/v1alpha1:ResourceClaimTemplate"
- "kubernetes:resource.k8s.io/v1alpha1:ResourceClaimTemplatePatch"
- "kubernetes:resource.k8s.io/v1alpha1:ResourceClaim"
- "kubernetes:resource.k8s.io/v1alpha1:ResourceClaimPatch"

### Resources moved from v1beta1 to v1

- "kubernetes:storage.k8s.io/v1beta1:CSIStorageCapacity"
- "kubernetes:storage.k8s.io/v1beta1:CSIStorageCapacityPatch"
- "kubernetes:storage.k8s.io/v1beta1:CSIStorageCapacityList"

### Resources renamed
- "kubernetes:resource.k8s.io/v1alpha1:PodSchedulingList"
  - Renamed to kubernetes:resource.k8s.io/v1alpha2:PodSchedulingContextList
- "kubernetes:resource.k8s.io/v1alpha1:PodSchedulingPatch"
  - Renamed to kubernetes:resource.k8s.io/v1alpha2:PodSchedulingContextPatch
- "kubernetes:resource.k8s.io/v1alpha1:PodScheduling"
  - Renamed to kubernetes:resource.k8s.io/v1alpha2:PodSchedulingContext

### New Features
- Allow instantiation of kustomize.Directory with a not fully configured provider (https://github.com/pulumi/pulumi-kubernetes/pull/2347)

## 3.24.3 (April 6, 2023)

- Handle CSA to SSA field manager conflicts (https://github.com/pulumi/pulumi-kubernetes/pull/2354)

## 3.24.2 (March 16, 2023)

- Update Pulumi Java SDK to v0.8.0 (https://github.com/pulumi/pulumi-kubernetes/pull/2337)
- Remove empty keys when merging unstructured resources for diffing (https://github.com/pulumi/pulumi-kubernetes/pull/2332)

## 3.24.1 (February 16, 2023)

- Move `invoke_yaml_decode` into ConfigGroup for python (https://github.com/pulumi/pulumi-kubernetes/pull/2317)
- Upgrade to latest helm dependency (https://github.com/pulumi/pulumi-kubernetes/pull/2318)

## 3.24.0 (February 6, 2023)

- Fix unencrypted secrets in the state `outputs` after `Secret.get` (https://github.com/pulumi/pulumi-kubernetes/pull/2300)
- Upgrade to latest helm and k8s client dependencies (https://github.com/pulumi/pulumi-kubernetes/pull/2292)
- Fix await status for Job and Pod (https://github.com/pulumi/pulumi-kubernetes/pull/2299)

## 3.23.1 (December 19, 2022)

- Add `PULUMI_K8S_ENABLE_PATCH_FORCE` env var support (https://github.com/pulumi/pulumi-kubernetes/pull/2260)
- Add link to resolution guide for SSA conflicts (https://github.com/pulumi/pulumi-kubernetes/pull/2265)
- Always set a field manager name to avoid conflicts in Client-Side Apply mode (https://github.com/pulumi/pulumi-kubernetes/pull/2271)

## 3.23.0 (December 8, 2022)

- Expose the allowNullValues boolean as an InputProperty so that it can be set in SDKs (https://github.com/pulumi/pulumi-kubernetes/pull/2255)
- Update Kubernetes support to Kubernetes v1.26.0 (https://github.com/pulumi/pulumi-kubernetes/pull/2230)

## 3.22.2 (November 30, 2022)

- Add allowNullValues boolean option to pass Null values through helm configs without having them
  scrubbed (https://github.com/pulumi/pulumi-kubernetes/issues/2089)
- Fix replacement behavior for immutable fields in SSA mode
  (https://github.com/pulumi/pulumi-kubernetes/issues/2235)
- For SSA conflicts, add a note to the error message about how to resolve
  (https://github.com/pulumi/pulumi-kubernetes/issues/2235)
- Make room for the `resource` API in Kubernetes 1.26.0 by qualifying the type of the same name in
  .NET code templates (https://github.com/pulumi/pulumi-kubernetes/pull/2237)

## 3.22.1 (October 26, 2022)

Note: Enabling SSA mode by default was causing problems for a number of users, so we decided to revert this change.
We plan to re-enable this as the default behavior in the next major (`v4.0.0`) release with additional documentation
about the expected differences.

- Revert: Enable Server-Side Apply mode by default (https://github.com/pulumi/pulumi-kubernetes/pull/2216)

## 3.22.0 (October 21, 2022)

Important Note -- This release changes the Provider default to enable Server-Side Apply mode. This change is
backward compatible, and should not require further action from users. The `enableServerSideApply` flag is
still present, so you may explicitly opt out if you run into any problems using one of the following methods:

1. Set the [enableServerSideApply](https://www.pulumi.com/registry/packages/kubernetes/api-docs/provider/#enable_server_side_apply_python)  parameter to `false` on your Provider resource.
2. Set the environment variable `PULUMI_K8S_ENABLE_SERVER_SIDE_APPLY="false"`
3. Set the stack config `pulumi config set kubernetes:enableServerSideApply false`

See the [how-to guide](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/managing-resources-with-server-side-apply/) for additional information about using Server-Side Apply with Pulumi's Kubernetes provider.

- Fix values precedence in helm release (https://github.com/pulumi/pulumi-kubernetes/pull/2191)
- Enable Server-Side Apply mode by default (https://github.com/pulumi/pulumi-kubernetes/pull/2206)

## 3.21.4 (September 22, 2022)

New tag to fix a publishing error for the Java SDK

## 3.21.3 (September 22, 2022)

- Fix Helm Chart preview with unconfigured provider (C#) (https://github.com/pulumi/pulumi-kubernetes/issues/2162)
- Load default kubeconfig if not specified in provider (https://github.com/pulumi/pulumi-kubernetes/issues/2180)
- Skip computing a preview for Patch resources (https://github.com/pulumi/pulumi-kubernetes/issues/2182)
- [sdk/python] Handle CRDs with status field input (https://github.com/pulumi/pulumi-kubernetes/issues/2183)
- Upgrade Kubernetes and Helm dependencies (https://github.com/pulumi/pulumi-kubernetes/issues/2186)

## 3.21.2 (September 1, 2022)

- Fix yaml bug resulting in `TypeError: Cannot read properties of undefined` (https://github.com/pulumi/pulumi-kubernetes/pull/2156)

## 3.21.1 (August 31, 2022)

- Update Helm and Kubernetes module dependencies (https://github.com/pulumi/pulumi-kubernetes/pull/2152)
- Automatically fill in .Capabilities in Helm Charts (https://github.com/pulumi/pulumi-kubernetes/pull/2155)

## 3.21.0 (August 23, 2022)

- Update Kubernetes support to Kubernetes v1.25.0 (https://github.com/pulumi/pulumi-kubernetes/pull/2129)

Breaking change note --
Kubernetes v1.25 dropped a few alpha and beta fields from the API, so the following fields are no longer available in
the provider SDKs:

* Type "kubernetes:batch/v1beta1:CronJobSpec" dropped property "timeZone"
* Type "kubernetes:batch/v1beta1:CronJobStatus" dropped property "lastSuccessfulTime"
* Type "kubernetes:discovery.k8s.io/v1beta1:ForZone" was dropped
* Type "kubernetes:discovery.k8s.io/v1beta1:Endpoint" dropped property "hints"
* Type "kubernetes:discovery.k8s.io/v1beta1:EndpointHints" dropped
* Type "kubernetes:policy/v1beta1:PodDisruptionBudgetStatus" dropped property "conditions"

## 3.20.5 (August 16, 2022)

- Update autonaming to use NewUniqueName for deterministic update plans. (https://github.com/pulumi/pulumi-kubernetes/pull/2137)
- Another fix for managed-by label in SSA mode. (https://github.com/pulumi/pulumi-kubernetes/pull/2140)

## 3.20.4 (August 15, 2022)

- Fix Helm charts being ignored by policy packs. (https://github.com/pulumi/pulumi-kubernetes/pull/2133)
- Fixes to allow import of helm release (https://github.com/pulumi/pulumi-kubernetes/pull/2136)
- Keep managed-by label in SSA mode if already present (https://github.com/pulumi/pulumi-kubernetes/pull/2138)

## 3.20.3 (August 9, 2022)

- Add chart v2 deprecation note to schema/docs (https://github.com/pulumi/pulumi-kubernetes/pull/2114)
- Add a descriptive message for an invalid Patch delete (https://github.com/pulumi/pulumi-kubernetes/pull/2111)
- Fix erroneous resourceVersion diff for CRDs managed with SSA (https://github.com/pulumi/pulumi-kubernetes/pull/2121)
- Update C# YAML GetResource implementation to compile with .NET v6 (https://github.com/pulumi/pulumi-kubernetes/pull/2122)
- Change .metadata.name to optional for all Patch resources (https://github.com/pulumi/pulumi-kubernetes/pull/2126)
- Fix field names in CRD schemas (https://github.com/pulumi/pulumi-kubernetes/pull/2128)

## 3.20.2 (July 25, 2022)

- Add Java SDK (https://github.com/pulumi/pulumi-kubernetes/pull/2096)
- Fix ServiceAccount readiness logic for k8s v1.24+ (https://github.com/pulumi/pulumi-kubernetes/issues/2099)

## 3.20.1 (July 19, 2022)

- Update the provider and tests to use Go 1.18. (https://github.com/pulumi/pulumi-kubernetes/pull/2073)
- Fix Helm Chart not working with Crossguard (https://github.com/pulumi/pulumi-kubernetes/pull/2057)
- Handle ignoreChanges for Server-Side Apply mode (https://github.com/pulumi/pulumi-kubernetes/pull/2074)

## 3.20.0 (July 12, 2022)

- Implement Server-Side Apply mode (https://github.com/pulumi/pulumi-kubernetes/pull/2029)
- Add Patch resources to all SDKs (https://github.com/pulumi/pulumi-kubernetes/pull/2043) (https://github.com/pulumi/pulumi-kubernetes/pull/2068)
- Add awaiter for service-account-token secret (https://github.com/pulumi/pulumi-kubernetes/pull/2048)
- Add Java packages overrides to schema (https://github.com/pulumi/pulumi-kubernetes/pull/2055)

## 3.19.4 (June 21, 2022)

- Use fully-qualified resource name for generating manifests, to avoid conflicts (https://github.com/pulumi/pulumi-kubernetes/pull/2007)
- Upgrade helm and k8s client-go module dependencies (https://github.com/pulumi/pulumi-kubernetes/pull/2008)
- Allow a user to opt-in to removing resources from Pulumi state when a cluster is unreachable (https://github.com/pulumi/pulumi-kubernetes/pull/2037)

## 3.19.3 (June 8, 2022)

- Fix a bug where the specified provider was not used for some operations on kustomize, helm, and yaml resources (https://github.com/pulumi/pulumi-kubernetes/pull/2005)

## 3.19.2 (May 25, 2022)

### Deprecations

- The `kubernetes:helm/v2:Chart` API is deprecated in this update and will be removed in a future release. The
  `kubernetes:helm/v3:Chart` resource is backward compatible, so changing the import path should not cause any resource
  updates.
- The `enableReplaceCRD` option on the Provider is deprecated in the update and will be removed in a future release.
  The behavior formerly enabled by this option is now default, and this option is ignored by the provider.

### Improvements

- Deprecate helm/v2:Chart resources (https://github.com/pulumi/pulumi-kubernetes/pull/1990)
- Don't use the last-applied-configuration annotation for CRDs (https://github.com/pulumi/pulumi-kubernetes/pull/1882)

## 3.19.1 (May 4, 2022)

- Upgrade pulumi/pulumi deps to v3.31.1 (https://github.com/pulumi/pulumi-kubernetes/pull/1980)

## 3.19.0 (May 3, 2022)

Note: The `kubernetes:storage.k8s.io/v1alpha1:CSIStorageCapacity` API was removed in this update.

- Update Kubernetes support to Kubernetes v1.24.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1911)

## 3.18.3 (April 21, 2022)
- Fix fetching remote yaml files (https://github.com/pulumi/pulumi-kubernetes/pull/1962)
- Support attach
  [#1977](https://github.com/pulumi/pulumi-kubernetes/pull/1977)

## 3.18.2 (April 6, 2022)
- Only add keyring default value when verification is turned on (https://github.com/pulumi/pulumi-kubernetes/pull/1961)
  Regression introduced in 3.18.1
- Fix the DaemonSet name on diff which prevented pulumi to replace the resource (https://github.com/pulumi/pulumi-kubernetes/pull/1951)

## 3.18.1 (April 5, 2022)
- Fix autonaming panic for helm release (https://github.com/pulumi/pulumi-kubernetes/pull/1953)
  This change also adds support for deterministic autonaming through sequence numbers to the kubernetes provider.

## 3.18.0 (March 31, 2022)
- Pass provider options to helm invokes in Python, Go and TS (https://github.com/pulumi/pulumi-kubernetes/pull/1919)
- Fix panic in helm release update() (https://github.com/pulumi/pulumi-kubernetes/pull/1948)

## 3.17.0 (March 14, 2022)
-  Make ConfigMaps mutable unless marked explicitly (enabled with provider config option) (https://github.com/pulumi/pulumi-kubernetes/pull/1926)
   *NOTE*: With this change, once `enableConfigMapMutable` is enabled, all ConfigMaps will be seen as mutable. In this mode, you can opt-in to the previous replacement behavior for a particular ConfigMap by setting its `replaceOnChanges` resource option to `[".binaryData", ".data"]`.
   By default, the provider will continue to treat ConfigMaps as immutable, and will replace them if the `binaryData` or `data` properties are changed.

## 3.16.0 (February 16, 2022)
- Bump to v3.8.0 of Helm (https://github.com/pulumi/pulumi-kubernetes/pull/1892)
- [Helm/Release][Helm/V3] Add initial support for OCI registries (https://github.com/pulumi/pulumi-kubernetes/pull/1892)

## 3.15.2 (February 9, 2022)
- Infer default namespace from kubeconfig when not configured via the provider (https://github.com/pulumi/pulumi-kubernetes/pull/1896)
- Fix an error handling bug in await logic (https://github.com/pulumi/pulumi-kubernetes/pull/1899)

## 3.15.1 (February 2, 2022)
- [Helm/Release] Add import docs (https://github.com/pulumi/pulumi-kubernetes/pull/1893)

## 3.15.0 (January 27, 2022)
- [Helm/Release] Remove beta warnings (https://github.com/pulumi/pulumi-kubernetes/pull/1885)
- [Helm/Release] Handle partial failure during create/update (https://github.com/pulumi/pulumi-kubernetes/pull/1880)
- [Helm/Release] Improve support for helm release operations when cluster is unreachable (https://github.com/pulumi/pulumi-kubernetes/pull/1886)
- [Helm/Release] Add examples to API reference docs and sdks (https://github.com/pulumi/pulumi-kubernetes/pull/1887)
- Fix detailed diff for server-side apply (https://github.com/pulumi/pulumi-kubernetes/pull/1873)
- Update to latest pulumi dependencies (https://github.com/pulumi/pulumi-kubernetes/pull/1888)

## 3.14.1 (January 18, 2022)

- Disable last-applied-configuration annotation for replaced CRDs (https://github.com/pulumi/pulumi-kubernetes/pull/1868)
- Fix Provider config diffs (https://github.com/pulumi/pulumi-kubernetes/pull/1869)
- Fix replace for named resource using server-side diff (https://github.com/pulumi/pulumi-kubernetes/pull/1870)
- Fix import for Provider using server-side diff (https://github.com/pulumi/pulumi-kubernetes/pull/1872)

## 3.14.0 (January 12, 2022)

- Fix panic for deletions from virtual fields in Helm Release (https://github.com/pulumi/pulumi-kubernetes/pull/1850)
- Switch Pod and Job await logic to external lib (https://github.com/pulumi/pulumi-kubernetes/pull/1856)
- Upgrade kubernetes provider module deps (https://github.com/pulumi/pulumi-kubernetes/pull/1861)

## 3.13.0 (January 7, 2022)
- Change await log type to cloud-ready-check lib (https://github.com/pulumi/pulumi-kubernetes/pull/1855)
- Populate inputs from live state for imports (https://github.com/pulumi/pulumi-kubernetes/pull/1846)
- Elide last-applied-configuration annotation when server-side support is enabled (https://github.com/pulumi/pulumi-kubernetes/pull/1863)
- Fix panic for deletions from virtual fields in Helm Release (https://github.com/pulumi/pulumi-kubernetes/pull/1850)

## 3.12.2 (January 5, 2022)

- Relax ingress await restrictions (https://github.com/pulumi/pulumi-kubernetes/pull/1832)
- Exclude nil entries from values (https://github.com/pulumi/pulumi-kubernetes/pull/1845)

## 3.12.1 (December 9, 2021)

- Helm Release: Helm Release imports support (https://github.com/pulumi/pulumi-kubernetes/pull/1818)
- Helm Release: fix username fetch option (https://github.com/pulumi/pulumi-kubernetes/pull/1824)
- Helm Release: Use URN name as base for autonaming, Drop warning, fix default value for
  keyring (https://github.com/pulumi/pulumi-kubernetes/pull/1826)
- Helm Release: Add support for loading values from yaml files (https://github.com/pulumi/pulumi-kubernetes/pull/1828)

- Fix CRD upgrades (https://github.com/pulumi/pulumi-kubernetes/pull/1819)

## 3.12.0 (December 7, 2021)

- Add support for k8s v1.23.0. (https://github.com/pulumi/pulumi-kubernetes/pull/1681)

## 3.11.0 (December 6, 2021)

Breaking change note:

[#1817](https://github.com/pulumi/pulumi-kubernetes/pull/1817) removed the deprecated providers/Provider
resource definition from the Go SDK. Following this change, use the Provider resource at
`github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes` instead.

- Helm Release: Make RepositoryOpts optional (https://github.com/pulumi/pulumi-kubernetes/pull/1806)
- Helm Release: Support local charts (https://github.com/pulumi/pulumi-kubernetes/pull/1809)
- Update pulumi dependencies v3.19.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1816)
- [go/sdk] Remove deprecated providers/Provider resource (https://github.com/pulumi/pulumi-kubernetes/pull/1817)

## 3.10.1 (November 19, 2021)

- Remove unused helm ReleaseType type (https://github.com/pulumi/pulumi-kubernetes/pull/1805)
- Fix Helm Release Panic "Helm uninstall returned information" (https://github.com/pulumi/pulumi-kubernetes/pull/1807)

## 3.10.0 (November 12, 2021)

- Add await support for networking.k8s.io/v1 variant of ingress (https://github.com/pulumi/pulumi-kubernetes/pull/1795)
- Schematize overlay types (https://github.com/pulumi/pulumi-kubernetes/pull/1793)

## 3.9.0 (November 5, 2021)

- [sdk/python] Add ready attribute to await Helm charts (https://github.com/pulumi/pulumi-kubernetes/pull/1782)
- [sdk/go] Add ready attribute to await Helm charts (https://github.com/pulumi/pulumi-kubernetes/pull/1784)
- [sdk/dotnet] Add ready attribute to await Helm charts (https://github.com/pulumi/pulumi-kubernetes/pull/1785)
- [sdk/python] Update CustomResource python implementation to pickup snake-case updates (https://github.com/pulumi/pulumi-kubernetes/pull/1786)
- Update pulumi dependencies v3.16.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1790)

## 3.8.3 (October 29, 2021)

- Add env variable lookup for k8s client settings (https://github.com/pulumi/pulumi-kubernetes/pull/1777)
- Fix diff logic for CRD kinds with the same name as a built-in (https://github.com/pulumi/pulumi-kubernetes/pull/1779)

## 3.8.2 (October 18, 2021)

- [sdk/python] Relax PyYaml dependency to allow upgrade to PyYaml 6.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1768)
- [go/sdk] Add missing types for deprecated Provider (https://github.com/pulumi/pulumi-kubernetes/pull/1771)

## 3.8.1 (October 8, 2021)

- Fix error for helm.Release previews with computed values (https://github.com/pulumi/pulumi-kubernetes/pull/1760)
- Don't require values for helm.Release (https://github.com/pulumi/pulumi-kubernetes/pull/1761)

## 3.8.0 (October 6, 2021)

Breaking change note:

[#1751](https://github.com/pulumi/pulumi-kubernetes/pull/1751) moved the Helm Release (beta) Provider options into a
complex type called `helmReleaseSettings`. Following this change, you can set these options in the following ways:

1. As arguments to a first-class Provider
   ```typescript
   new k8s.Provider("test", { helmReleaseSettings: { driver: "secret" } });
   ```

1. Stack configuration for the default Provider
   ```
   pulumi config set --path kubernetes:helmReleaseSettings.driver "secret"
   ```

1. As environment variables
   ```
   EXPORT PULUMI_K8S_HELM_DRIVER="secret"
   ```

- [sdk/dotnet] Fix creation of CustomResources (https://github.com/pulumi/pulumi-kubernetes/pull/1741)
- Always override namespace for helm release operations (https://github.com/pulumi/pulumi-kubernetes/pull/1747)
- Add k8s client tuning settings to Provider (https://github.com/pulumi/pulumi-kubernetes/pull/1748)
- Nest helm.Release Provider settings (https://github.com/pulumi/pulumi-kubernetes/pull/1751)
- Change await logic client to use target apiVersion on updates (https://github.com/pulumi/pulumi-kubernetes/pull/1758)

## 3.7.3 (September 30, 2021)
- Use helm release's namespace on templates where namespace is left unspecified (https://github.com/pulumi/pulumi-kubernetes/pull/1733)
- Upgrade Helm dependency to v3.7.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1742)
- Helm Release: Await deletion if skipAwait is unset or atomic is specific (https://github.com/pulumi/pulumi-kubernetes/pull/1742)

## 3.7.2 (September 17, 2021)
- Fix handling of charts with empty manifests (https://github.com/pulumi/pulumi-kubernetes/pull/1717)
- Use existing helm template logic to populate manifests instead of relying on `dry-run` support (https://github.com/pulumi/pulumi-kubernetes/pull/1718)

## 3.7.1 (September 10, 2021)
- Don't replace PVC on .spec.resources.requests or .limits change. (https://github.com/pulumi/pulumi-kubernetes/pull/1705)
    - *NOTE*: User's will now need to use the `replaceOnChanges` resource option for PVCs if modifying requests or limits to trigger replacement

## 3.7.0 (September 3, 2021)
- Add initial support for a Helm release resource - `kubernetes:helm.sh/v3:Release. Currently available in Beta (https://github.com/pulumi/pulumi-kubernetes/pull/1677)

## 3.6.3 (August 23, 2021)

- [sdk/go] Re-add deprecated Provider file (https://github.com/pulumi/pulumi-kubernetes/pull/1687)

## 3.6.2 (August 20, 2021)

- Fix environment variable name in disable Helm hook warnings message (https://github.com/pulumi/pulumi-kubernetes/pull/1683)

## 3.6.1 (August 19, 2021)

- [sdk/python] Fix wait for metadata in `yaml._parse_yaml_object`. (https://github.com/pulumi/pulumi-kubernetes/pull/1675)
- Fix diff logic for server-side apply mode (https://github.com/pulumi/pulumi-kubernetes/pull/1679)
- Add option to disable Helm hook warnings (https://github.com/pulumi/pulumi-kubernetes/pull/1682)
- For renderToYamlDirectory, treat an empty directory as unset (https://github.com/pulumi/pulumi-kubernetes/pull/1678)

## 3.6.0 (August 4, 2021)

The following breaking changes are part of the Kubernetes v1.22 update:
- The alpha `EphemeralContainers` kind [has been removed](https://github.com/kubernetes/kubernetes/pull/101034)
- [.NET SDK] `Networking.V1Beta1.IngressClassParametersReferenceArgs` -> `Core.V1.TypedLocalObjectReferenceArgs`

- Update Helm and client-go deps (https://github.com/pulumi/pulumi-kubernetes/pull/1662)
- Add support for k8s v1.22.0. (https://github.com/pulumi/pulumi-kubernetes/pull/1551)

## 3.5.2 (July 29, 2021)

- Update pulumi dependencies v3.7.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1651)
- Update pulumi dependencies v3.9.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1660)

## 3.5.1 (July 14, 2021)

- Use shared informer for await logic for all resources (https://github.com/pulumi/pulumi-kubernetes/pull/1647)

## 3.5.0 (June 30, 2021)

- Update pulumi dependencies v3.5.1 (https://github.com/pulumi/pulumi-kubernetes/pull/1623)
- Skip cluster connectivity check in yamlRenderMode (https://github.com/pulumi/pulumi-kubernetes/pull/1629)
- Handle different namespaces for server-side diff (https://github.com/pulumi/pulumi-kubernetes/pull/1631)
- Handle auto-named namespaces for server-side diff (https://github.com/pulumi/pulumi-kubernetes/pull/1633)
- *Revert* Fix hanging updates for deployment await logic (https://github.com/pulumi/pulumi-kubernetes/pull/1596)
- Use shared informer for await logic for deployments (https://github.com/pulumi/pulumi-kubernetes/pull/1639)

## 3.4.1 (June 24, 2021)

- *Revert* Fix hanging updates for deployment await logic (https://github.com/pulumi/pulumi-kubernetes/pull/1596)

## 3.4.0 (June 17, 2021)

- Add skipAwait option to helm.v3 SDKs. (https://github.com/pulumi/pulumi-kubernetes/pull/1603)
- Add skipAwait option to YAML SDKs. (https://github.com/pulumi/pulumi-kubernetes/pull/1610)
- [sdk/go] `ConfigGroup` now respects explicit provider instances when parsing YAML. (https://github.com/pulumi/pulumi-kubernetes/pull/1601)
- Fix hanging updates for deployment await logic (https://github.com/pulumi/pulumi-kubernetes/pull/1596)
- Fix auto-naming bug for computed names (https://github.com/pulumi/pulumi-kubernetes/pull/1618)
- [sdk/python] Revert pulumi dependency to <4.0.0. (https://github.com/pulumi/pulumi-kubernetes/pull/1619)
- [sdk/dotnet] Unpin System.Collections.Immutable dependency (https://github.com/pulumi/pulumi-kubernetes/pull/1621)

## 3.3.1 (June 8, 2021)

- [sdk/python] Fix YAML regression by pinning pulumi dependency to <3.4.0. (https://github.com/pulumi/pulumi-kubernetes/pull/1605)

## 3.3.0 (May 26, 2021)

- Automatically mark Secret data as Pulumi secrets. (https://github.com/pulumi/pulumi-kubernetes/pull/1577)
- Update pulumi dependency (https://github.com/pulumi/pulumi-kubernetes/pull/1588)
    - [codegen] Automatically encrypt secret input parameters (https://github.com/pulumi/pulumi/pull/7128)
    - [sdk/python] Nondeterministic import ordering fix from (https://github.com/pulumi/pulumi/pull/7126)

## 3.2.0 (May 19, 2021)

- Allow opting out of CRD rendering for Helm v3 by specifying `SkipCRDRendering` argument to Helm charts. (https://github.com/pulumi/pulumi-kubernetes/pull/1572)
- Add replaceUnready annotation for Jobs. (https://github.com/pulumi/pulumi-kubernetes/pull/1575)
- Read live Job state for replaceUnready check. (https://github.com/pulumi/pulumi-kubernetes/pull/1578)
- Fix diff panic for malformed kubeconfig. (https://github.com/pulumi/pulumi-kubernetes/pull/1581)

## 3.1.2 (May 12, 2021)

- Update pulumi dependencies to fix Python regression from 3.1.1 (https://github.com/pulumi/pulumi-kubernetes/pull/1573)

## 3.1.1 (May 5, 2021)

- Avoid circular dependencies in NodeJS SDK modules (https://github.com/pulumi/pulumi-kubernetes/pull/1558)
- Update pulumi dependencies v3.2.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1564)

## 3.1.0 (April 29, 2021)

- Update Helm to v3.5.4 and client-go to v0.20.4 (https://github.com/pulumi/pulumi-kubernetes/pull/1536)
- In Python helmv2 and helmv3, and Node.js helmv3, no longer pass `Chart` resource options to child resources explicitly.
  (https://github.com/pulumi/pulumi-kubernetes/pull/1539)

## 3.0.0 (April 19, 2021)

- Depend on Pulumi 3.0, which includes improvements to Python resource arguments and key translation, Go SDK performance,
  Node SDK performance, general availability of Automation API, and more.

- Update pulumi dependency (https://github.com/pulumi/pulumi-kubernetes/pull/1521)
    - [sdk/go] Fix Go resource registrations (https://github.com/pulumi/pulumi/pull/6641)
    - [sdk/python] Support `<Resource>Args` classes (https://github.com/pulumi/pulumi/pull/6525)

- Do not return an ID when previewing the creation of a resource. (https://github.com/pulumi/pulumi-kubernetes/pull/1526)

## 2.9.1 (April 12, 2021)

- Fix refresh for render-yaml resources (https://github.com/pulumi/pulumi-kubernetes/pull/1523)
- Behavior change: Error on refresh for an unreachable cluster. (https://github.com/pulumi/pulumi-kubernetes/pull/1522)

## 2.9.0 (April 8, 2021)

- Add support for k8s v1.21.0. (https://github.com/pulumi/pulumi-kubernetes/pull/1449)
- [sdk/go] Fix plugin versioning for invoke calls (https://github.com/pulumi/pulumi-kubernetes/pull/1520)

## 2.8.4 (March 29, 2021)

- Ensure using `PULUMI_KUBERNETES_MANAGED_BY_LABEL` doesn't cause diffs on further stack updates (https://github.com/pulumi/pulumi-kubernetes/pull/1508)
- [sdk/ts] Update CustomResource to match current codegen (https://github.com/pulumi/pulumi-kubernetes/pull/1510)

## 2.8.3 (March 17, 2021)

- Fix bug where rendering manifests results in files being overwritten by subsequent resources with the same kind and name, but different namespace (https://github.com/pulumi/pulumi-kubernetes/pull/1429)
- Update pulumi dependency to fix python Resource.get() functions (https://github.com/pulumi/pulumi-kubernetes/pull/1480)
- Upgrade to Go1.16 (https://github.com/pulumi/pulumi-kubernetes/pull/1486)
- Adding arm64 plugin builds (https://github.com/pulumi/pulumi-kubernetes/pull/1490)
- Fix bug preventing helm chart being located in bitnami repo (https://github.com/pulumi/pulumi-kubernetes/pull/1491)

## 2.8.2 (February 23, 2021)

-   Postpone the removal of admissionregistration/v1beta1, which has been retargeted at 1.22 (https://github.com/pulumi/pulumi-kubernetes/pull/1474)
-   Change k8s API removals from error to warning (https://github.com/pulumi/pulumi-kubernetes/pull/1475)

## 2.8.1 (February 12, 2021)

-   Skip Helm test hook resources by default (https://github.com/pulumi/pulumi-kubernetes/pull/1467)
-   Ensure no panic when a kubernetes provider is used with an incompatible resource type (https://github.com/pulumi/pulumi-kubernetes/pull/1469)
-   Allow users to set `PULUMI_KUBERNETES_MANAGED_BY_LABEL` environment variable to control `app.kubernetes.io/managed-by` label (https://github.com/pulumi/pulumi-kubernetes/pull/1471)

## 2.8.0 (February 3, 2021)

Note: This release fixes a bug with the Helm v3 SDK that omitted any chart resources that included a hook annotation.
If you have existing charts deployed with the v3 SDK that include hook resources, the next update will create these
resources.

-   [Go SDK] Fix bug with v1/List in YAML parsing (https://github.com/pulumi/pulumi-kubernetes/pull/1457)
-   Fix bug rendering Helm v3 resources that include hooks (https://github.com/pulumi/pulumi-kubernetes/pull/1459)
-   Print warning for Helm resources using unsupported hooks (https://github.com/pulumi/pulumi-kubernetes/pull/1460)

## 2.7.8 (January 27, 2021)

-   Update pulumi dependency to remove unused Go types (https://github.com/pulumi/pulumi-kubernetes/pull/1450)

## 2.7.7 (January 20, 2021)

-   Expand allowed Python pyyaml dependency versions (https://github.com/pulumi/pulumi-kubernetes/pull/1435)

## 2.7.6 (January 13, 2021)

-   Upgrade helm and k8s deps (https://github.com/pulumi/pulumi-kubernetes/pull/1414)
-   Fix bug with Go Helm v3 transformation marshaling (https://github.com/pulumi/pulumi-kubernetes/pull/1420)

## 2.7.5 (December 16, 2020)

-   Add enum for kubernetes.core.v1.Service.Spec.Type (https://github.com/pulumi/pulumi-kubernetes/pull/1408)
-   Update pulumi deps to v2.15.5 (https://github.com/pulumi/pulumi-kubernetes/pull/1402)
-   Fix Go resource Input/Output methods (https://github.com/pulumi/pulumi-kubernetes/pull/1406)

## 2.7.4 (December 8, 2020)

-   Add support for k8s v1.20.0. (https://github.com/pulumi/pulumi-kubernetes/pull/1330)

## 2.7.3 (December 3, 2020)

-   Replace workload resources if any field in `.spec.selector` changes. (https://github.com/pulumi/pulumi-kubernetes/pull/1387)
-   Update pulumi deps to v2.15.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1393)
-   Add package/module registration for NodeJS and Python (https://github.com/pulumi/pulumi-kubernetes/pull/1394)

## 2.7.2 (November 19, 2020)

-   Fixed a gRPC error for larger Helm charts in the .NET SDK [#4224](https://github.com/pulumi/pulumi/issues/4224)
-   Update pulumi deps to v2.14.0 (https://github.com/pulumi/pulumi-kubernetes/pull/1385)

## 2.7.1 (November 18, 2020)

-   Error on delete if cluster is unreachable (https://github.com/pulumi/pulumi-kubernetes/pull/1379)
-   Update pulumi deps to v2.13.2 (https://github.com/pulumi/pulumi-kubernetes/pull/1383)

## 2.7.0 (November 12, 2020)

-   Add support for previewing Create and Update operations for API servers that support dry-run (https://github.com/pulumi/pulumi-kubernetes/pull/1355)
-   Fix panic introduced in #1355 (https://github.com/pulumi/pulumi-kubernetes/pull/1368)
-   (NodeJS) Add ready attribute to await Helm charts (https://github.com/pulumi/pulumi-kubernetes/pull/1364)
-   Update Helm to v3.4.0 and client-go to v1.19.2 (https://github.com/pulumi/pulumi-kubernetes/pull/1360)
-   Update Helm to v3.4.1 and client-go to v1.19.3 (https://github.com/pulumi/pulumi-kubernetes/pull/1365)
-   Fix panic when provider is given kubeconfig as an object instead of string (https://github.com/pulumi/pulumi-kubernetes/pull/1373)
-   Fix concurrency issue in Helm + .NET SDK [#1311](https://github.com/pulumi/pulumi-kubernetes/issues/1311) and [#1374](https://github.com/pulumi/pulumi-kubernetes/issues/1374)

## 2.6.3 (October 12, 2020)

-   Revert Helm v2 deprecation warnings (https://github.com/pulumi/pulumi-kubernetes/pull/1352)

## 2.6.2 (October 7, 2020)

## Important Note

Helm v2 support is [EOL](https://helm.sh/blog/helm-v2-deprecation-timeline/), and will no longer be supported upstream
as of next month. Furthermore, the stable/incubator chart repos will likely
[stop working](https://github.com/helm/charts#deprecation-timeline) after November 13, 2020. Deprecation warnings have
been added for any usage of Pulumi's `helm.v2` API, and this API will be removed at a future date. Our `helm.v3` API is
backward compatible, so you should be able to update without disruption to existing resources.

### Bug Fixes

-   Set plugin version for Go SDK invoke calls (https://github.com/pulumi/pulumi-kubernetes/pull/1325)
-   Python: Fix generated examples and docs to prefer input/output classes (https://github.com/pulumi/pulumi-kubernetes/pull/1346)

### Improvements

-   Update Helm v3 mod to v3.3.2 (https://github.com/pulumi/pulumi-kubernetes/pull/1326)
-   Update Helm v3 mod to v3.3.3 (https://github.com/pulumi/pulumi-kubernetes/pull/1328)
-   Change error to warning if internal autoname annotation is set (https://github.com/pulumi/pulumi-kubernetes/pull/1337)
-   Deprecate Helm v2 SDKs (https://github.com/pulumi/pulumi-kubernetes/pull/1344)

## 2.6.1 (September 16, 2020)

### Bug Fixes

-   Fix Python type hints for lists (https://github.com/pulumi/pulumi-kubernetes/pull/1313)
-   Fix Python type hints for integers (https://github.com/pulumi/pulumi-kubernetes/pull/1317)
-   Fix Helm v3 default namespace handling (https://github.com/pulumi/pulumi-kubernetes/pull/1323)

### Improvements

-   Update Helm v3 mod to v3.3.1 (https://github.com/pulumi/pulumi-kubernetes/pull/1320)

## 2.6.0 (September 10, 2020)

Note: There is a minor breaking change in the .NET SDK for Helm v3. As part of the switch to using native
Helm libraries in #1291, the Helm.V3.Chart class no longer inherits from the ChartBase class. Most users should
not be affected by this change.

### Bug Fixes

-   Upgrade version of pyyaml to fix a [security vulnerability](https://nvd.nist.gov/vuln/detail/CVE-2019-20477) (https://github.com/pulumi/pulumi-kubernetes/pull/1230)
-   Fix Helm api-versions handling in all SDKs. (https://github.com/pulumi/pulumi-kubernetes/pull/1307)

### Improvements

-   Update .NET Helm v3 to use native client. (https://github.com/pulumi/pulumi-kubernetes/pull/1291)
-   Update Go Helm v3 to use native client. (https://github.com/pulumi/pulumi-kubernetes/pull/1296)
-   Python: Allow type annotations on transformation functions. (https://github.com/pulumi/pulumi-kubernetes/pull/1298)

## 2.5.1 (September 2, 2020)

### Bug Fixes

-   Fix regression of .get methods in NodeJS SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/1285)

### Improvements

-   Upgrade to Pulumi v2.9.0, which adds type annotations and input/output classes to Python (https://github.com/pulumi/pulumi-kubernetes/pull/1276)
-   Switch Helm v3 logic to use native library. (https://github.com/pulumi/pulumi-kubernetes/pull/1263)
-   Bump python requests version dependency. (https://github.com/pulumi/pulumi-kubernetes/pull/1274)
-   Update NodeJS Helm v3 to use native client. (https://github.com/pulumi/pulumi-kubernetes/pull/1279)
-   [sdk/nodejs] Remove unneccessary constructor overloads. (https://github.com/pulumi/pulumi-kubernetes/pull/1286)

## 2.5.0 (August 26, 2020)

### Improvements

-   Add support for k8s v1.19.0. (https://github.com/pulumi/pulumi-kubernetes/pull/996)
-   Handle kubeconfig contents or path in provider. (https://github.com/pulumi/pulumi-kubernetes/pull/1255)
-   Add type annotations to Python SDK for API Extensions, Helm, Kustomize, and YAML. (https://github.com/pulumi/pulumi-kubernetes/pull/1259)
-   Update k8s package deps to v0.18.8. (https://github.com/pulumi/pulumi-kubernetes/pull/1265)
-   Move back to upstream json-patch module. (https://github.com/pulumi/pulumi-kubernetes/pull/1266)

## 2.4.3 (August 14, 2020)

### Bug Fixes

-   Rename Python's yaml.ConfigFile file_id parameter to file. (https://github.com/pulumi/pulumi-kubernetes/pull/1248)

### Improvements

-   Remove the ComponentStatus resource type. (https://github.com/pulumi/pulumi-kubernetes/pull/1234)

## 2.4.2 (August 3, 2020)

### Bug Fixes

-   Fix server-side diff when immutable fields change. (https://github.com/pulumi/pulumi-kubernetes/pull/988)
-   Update json-patch mod to fix hangs on pulumi update. (https://github.com/pulumi/pulumi-kubernetes/pull/1223)

## 2.4.1 (July 24, 2020)

### Bug Fixes

-   Handle networking/v1beta1 Ingress resources. (https://github.com/pulumi/pulumi-kubernetes/pull/1221)

### Improvements

-   Add NodeJS usage examples for Helm, Kustomize, and YAML resources. (https://github.com/pulumi/pulumi-kubernetes/pull/1205)
-   Add Python usage examples for Helm, Kustomize, and YAML resources. (https://github.com/pulumi/pulumi-kubernetes/pull/1209)
-   Add v3 Helm package for Go SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/1211)
-   Add Go usage examples for Helm, Kustomize, and YAML resources. (https://github.com/pulumi/pulumi-kubernetes/pull/1212)
-   Add yaml.ConfigGroup to Python SDK. (https://github.com/pulumi/pulumi-kubernetes/pull/1217)
-   Add C# usage examples for Helm, Kustomize, and YAML resources. (https://github.com/pulumi/pulumi-kubernetes/pull/1213)

## 2.4.0 (July 7, 2020)

### Bug Fixes

-   Fix error parsing Helm version (https://github.com/pulumi/pulumi-kubernetes/pull/1170)
-   Fix prometheus-operator test to wait for the CRD to be ready before use (https://github.com/pulumi/pulumi-kubernetes/pull/1172)
-   Fix suppress deprecation warnings flag (https://github.com/pulumi/pulumi-kubernetes/pull/1189)
-   Set additionalSecretOutputs on Secret data fields (https://github.com/pulumi/pulumi-kubernetes/pull/1194)

### Improvements

-   Set supported environment variables in SDK Provider classes (https://github.com/pulumi/pulumi-kubernetes/pull/1166)
-   Python SDK updated to align with other Pulumi Python SDKs. (https://github.com/pulumi/pulumi-kubernetes/pull/1160)
-   Add support for Kustomize. (https://github.com/pulumi/pulumi-kubernetes/pull/1178)
-   Implement GetSchema to enable example and import code generation. (https://github.com/pulumi/pulumi-kubernetes/pull/1181)
-   Only show deprecation messages when new API versions exist in current cluster version (https://github.com/pulumi/pulumi-kubernetes/pull/1182)

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
