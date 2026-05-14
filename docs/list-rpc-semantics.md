# Kubernetes List Semantics

## Summary

[#4332](https://github.com/pulumi/pulumi-kubernetes/issues/4332) implements
the Pulumi provider `List` RPC for Kubernetes resources. The provider resolves
a Pulumi type token to a `GroupVersionKind`, maps it through the discovery
RESTMapper, and streams results from the K8s dynamic client.

All generated non-`Patch` K8s API resources are listable.

Merged as [#4333](https://github.com/pulumi/pulumi-kubernetes/pull/4333).

## Design

The K8s dynamic client already speaks the paging, label-selector, and
field-selector vocabulary Pulumi's `List` RPC needs, so the provider mostly
translates.

- `provider/pkg/provider/provider.go`: `kubeProvider.List` resolves the
  token to a GVK and REST mapping, applies query filters, and streams results.
- `provider/pkg/provider/list_pagination.go`: provider-owned continuation
  state (K8s cursor + remaining budget).
- `provider/pkg/gen/schema.go`: emits `listInputs` for every listable kind.

### Schema

A resource is listable if it is not nested and not a `Patch` resource.
Patch tokens have their `Patch` suffix stripped at runtime and route through
the underlying API resource.

Every listable resource gets the same four `listInputs`:

- `name`
- `labelSelector`
- `fieldSelector`
- `namespace` (omitted for known cluster-scoped kinds via
  `kinds.Kind(...).Namespaced()`)

All four are optional.

## Use Cases

### Import ID discovery

Pick a kind, scan returned IDs, hand the selected ID to `pulumi import`.
Kubernetes IDs are `namespace/name` (or just `name` for cluster-scoped),
so they are inherently recognizable.

### Scoped child discovery

`kubernetes:core/v1:Pod` with `namespace: kube-system` lists only pods
in `kube-system`.

### Bulk inventory and bulk import tooling

Higher-level tools can enumerate IDs page by page and decide what to do
with them. `List` returns only importable IDs, not full state; callers
that need state must follow up with `Read`.

### Pulumi program inputs

Programs can use `List` results inside a stack, but the response shape
(`id` + `name` only) limits this to ID-driven selection. Programs that
select a resource by tag, label, or property should use data sources or
`getResource`-style helpers.

## Filter / Query parameters

Four query fields, all optional, all strings:

| Field           | Maps to                                | Notes |
|-----------------|----------------------------------------|-------|
| `name`          | `fieldSelector=metadata.name=<value>` | Combined with any explicit `fieldSelector` via `,`. |
| `labelSelector` | `metav1.ListOptions.LabelSelector`     | Full K8s selector grammar (`key=value`, `!key`, `key in (a,b)`, etc.). |
| `fieldSelector` | `metav1.ListOptions.FieldSelector`     | Server-side fields only; varies per kind. |
| `namespace`     | `dynamic.Interface.Namespace(...)`     | Rejected with `InvalidArgument` for cluster-scoped kinds. |

The query also accepts a `metadata.{namespace,name}` fallback. Top-level
keys win when both are set.

`labelSelector` and `fieldSelector` are interpreted server-side and passed
through unmodified. Syntax errors surface as K8s API errors.

## Page size and continuation

The protocol sends the same `limit` on every paginated call, so the
provider tracks how many results it has already delivered.

The continuation token is a base64-encoded JSON blob carrying:

- `k8sContinue`: K8s's own cursor for the next page.
- `remaining`: items still allowed under the session-wide `limit`
  (`*int64`; see [Edge Cases](#edge-cases)).

Per call:

1. Decode the incoming token. Empty means first call.
2. Compute `effectiveLimit = min(page_size, remaining)`, treating `0` as
   "no cap" on the K8s side.
3. Call `resourceClient.List(...)` with `Continue=k8sContinue` and
   `Limit=effectiveLimit`.
4. Stream a `Result` per item.
5. Update `remaining`, build the next state, and emit a `Continuation`
   message only if there's anything left to do.

## Edge cases

### `Remaining *int64` instead of `int64`

The provider needs to distinguish three states:

| State                  | Encoding         | Meaning |
|------------------------|------------------|---------|
| No cap (limit unset)   | `nil`            | Stream until K8s says we're done. |
| Cap with budget        | `*N` where N > 0 | N items still allowed. |
| Cap exhausted          | `*0`             | Stop; do not call K8s again. |

A plain `int64` collapses the first and third states into `0`, the proto3
zero-value problem. The pointer disambiguates and survives JSON round-trips
through the continuation token.

`effectiveLimit` panics on `*0`; the caller is required to stop before
invoking it.

### Patch resources

`XPatch` tokens are not separate API resources; they share the underlying
kind. `List` strips the `Patch` suffix from `gvk.Kind` and routes through
the underlying resource.

### Cluster-scoped kinds

The schema omits `namespace` from `listInputs` for known cluster-scoped
kinds. If a `namespace` query arrives anyway, the provider rejects it
with `InvalidArgument`.

### Stale discovery cache

A freshly installed CRD will not yet be in the RESTMapper. On a mapping
miss, the provider resets the discovery cache and retries once. A second
miss surfaces as the underlying error.

### Empty results

K8s returns an empty list when nothing matches. The provider sends zero
`Result` messages and no `Continuation`.

## Non-Goals

### Helm Release List

`kubernetes:helm.sh/v3:Release` wraps the Helm Go client, not a K8s API
resource. Listing Releases requires a separate code path through
`action.NewList(cfg).Run()` and belongs in its own issue.

### Helm Chart, kustomize, yaml overlays

These are generator resources that expand client-side into K8s manifests.
They have no persistent cluster-side identity. Users wanting inventory
should list the K8s resources the generators produced.

### `CustomResource` overlay List

Listing arbitrary CRD instances already works via the direct GVK token
(`kubernetes:cert-manager.io/v1:Certificate`); runtime resolves through
the RESTMapper.

