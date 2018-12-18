## 0.18.1 (Unreleased)

* Implement incremental status updates for `StatefulSet`
  (https://github.com/pulumi/pulumi-kubernetes/pull/307)

## 0.18.0 (Released December 4, 2018)

### Major changes

* Allow Helm Charts to have `pulumi.Input` in their `values`
  (https://github.com/pulumi/pulumi-kubernetes/pull/241)

### Improvements

* Retry REST calls to Kubernetes if they fail, greatly improving resiliance against resorce
  operation ordering problems.
* Add support for creating CRDs and CRs in the same app
  (https://github.com/pulumi/pulumi-kubernetes/pull/271,
  https://github.com/pulumi/pulumi-kubernetes/pull/280)
* Add incremental await for logic for `Ingress`
  (https://github.com/pulumi/pulumi-kubernetes/pull/283)
* Allow users to specify a Chart's source any way they can do it from the CLI
  (https://github.com/pulumi/pulumi-kubernetes/pull/284)
* "Fix" "bug" that cases Pulumi to crash if there is a duplicate key in a YAML template, to conform
  with Helm's behavior (https://github.com/pulumi/pulumi-kubernetes/pull/289)
* Emit better error when the API server is unreachable
  (https://github.com/pulumi/pulumi-kubernetes/pull/291)
* Add support for Kubernetes v0.12.* (https://github.com/pulumi/pulumi-kubernetes/pull/293)
* Fix bug that spuriously requires `.metadata.name` to be specified in Kubernetes list types
  (_e.g._, `v1/List`) (https://github.com/pulumi/pulumi-kubernetes/pull/294,
  https://github.com/pulumi/pulumi-kubernetes/pull/296)
* Add Kubernetes v0.13.* support (https://github.com/pulumi/pulumi-kubernetes/pull/306)
* Improve error message when `Service` fails to initialized
  (https://github.com/pulumi/pulumi-kubernetes/pull/309)
* Fix bug that causes us to erroneously report `Pod`'s owner
  (https://github.com/pulumi/pulumi-kubernetes/pull/311)
