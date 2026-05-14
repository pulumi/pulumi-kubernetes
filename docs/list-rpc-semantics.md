# Kubernetes List Semantics

## Summary

[#4332](https://github.com/pulumi/pulumi-kubernetes/issues/4332) implements
the Pulumi provider `List` RPC for Kubernetes resources. The provider resolves
a Pulumi type token to a `GroupVersionKind`, maps it through the discovery
RESTMapper, and streams results from the K8s dynamic client.

Every non-nested, non-`Patch` generated resource is listable. The 10 overlay
resources (Helm Chart/Release, kustomize, yaml, CustomResource) are out of
scope for this pass — they don't correspond to a single Kubernetes API
endpoint and are addressed under [Non-Goals](#non-goals).

Built on the runtime draft from [@Frassle](https://github.com/Frassle) on the
`fraser/list` branch; merged as [#4333](https://github.com/pulumi/pulumi-kubernetes/pull/4333).

## Design

The implementation is deliberately small: the K8s dynamic client already
speaks the same paging, label-selector, and field-selector vocabulary that
Pulumi's `List` RPC needs, so the provider mostly translates.

- `provider/pkg/provider/provider.go` — `kubeProvider.List` resolves the
  token → GVK → REST mapping, applies query filters, and streams results.
- `provider/pkg/provider/list_pagination.go` — provider-owned continuation
  state (K8s cursor + remaining budget).
- `provider/pkg/gen/schema.go` — emits `listInputs` for every listable kind.

### Schema

A resource is listable if it is not nested and not a `Patch` resource.
Patch tokens are stripped of their `Patch` suffix at runtime and routed to
the underlying API resource.

Every listable resource gets the same four `listInputs`:

- `name`
- `labelSelector`
- `fieldSelector`
- `namespace` (omitted for known cluster-scoped kinds via
  `kinds.Kind(...).Namespaced()`)

All four are optional. There is no required parent scope — K8s lets you list
namespaced kinds across all namespaces by omitting `namespace`.

## Use Cases

The Pulumi `List` RPC is designed around four use cases. The first three
apply directly to Kubernetes; the fourth is partially supported.

### Import ID discovery

Pick a kind, scan returned IDs, hand the selected ID to `pulumi import`.
Kubernetes IDs are `namespace/name` (or just `name` for cluster-scoped),
which are inherently recognizable, so this is the strongest use case.

### Scoped child discovery

`namespace` is the natural scope filter. `kubernetes:core/v1:Pod` with
`namespace: kube-system` lists only pods in `kube-system`. This is the
Kubernetes analog of AWS Native's "scoped child discovery," except the
scope is optional rather than required.

### Bulk inventory and bulk import tooling

Higher-level tools can enumerate IDs page by page and decide what to do
with them. `List` returns only importable IDs, not full state — callers
that need state must follow up with `Read`.

### Pulumi program inputs

Programs can use `List` results inside a stack, but the response shape
(`id` + `name` only) limits this to ID-driven selection. Programs that
need to select a resource by tag, label, or property should still use
data sources or `getResource`-style helpers.

## Filter / Query parameters

Four query fields, all optional, all strings:

| Field           | Maps to                                | Notes |
|-----------------|----------------------------------------|-------|
| `name`          | `fieldSelector=metadata.name=<value>` | Combined with any explicit `fieldSelector` via `,`. |
| `labelSelector` | `metav1.ListOptions.LabelSelector`     | Full K8s selector grammar (`key=value`, `!key`, `key in (a,b)`, etc.). |
| `fieldSelector` | `metav1.ListOptions.FieldSelector`     | Server-side fields only — varies per kind. |
| `namespace`     | `dynamic.Interface.Namespace(...)`     | Rejected with `InvalidArgument` for cluster-scoped kinds. |

For ergonomics, the query also accepts a `metadata.{namespace,name}`
fallback, since users often paste a partial resource manifest. Top-level
keys win when both are set.

K8s `LabelSelector` and `FieldSelector` are interpreted server-side; the
provider passes them through unmodified, so any error in selector syntax
surfaces as a K8s API error.

## Page size and continuation

The K8s API and the Pulumi protocol both have the same notion of paging,
but they don't compose for free: the protocol says "the engine sends the
same `limit` on every paginated call," so the provider has to remember how
many results it has already delivered.

Our continuation token is therefore not K8s's raw token. It is a
base64-encoded JSON blob carrying both:

- `k8sContinue` — K8s's own cursor for the next page;
- `remaining` — items still allowed under the session-wide `limit`
  (`*int64`, see [Edge Cases](#edge-cases) for why).

Per call:

1. Decode the incoming token (empty = first call).
2. Compute `effectiveLimit = min(page_size, remaining)`, treating `0` as
   "no cap" on the K8s side.
3. Call `resourceClient.List(...)` with `Continue=k8sContinue` and
   `Limit=effectiveLimit`.
4. Stream a `Result` per item.
5. Update `remaining`, build the next state, and emit a `Continuation`
   message only if there's anything left to do.

This matches the AWS Native model in intent — the provider owns the token,
the engine never sees the underlying API's token — but is simpler because
K8s's `Continue` is already an opaque server-side cursor; we don't need
session storage or in-memory buffering.

## Edge cases

### `Remaining *int64` — why a pointer instead of `int64`

The provider needs to distinguish three states:

| State                  | Encoding         | Meaning |
|------------------------|------------------|---------|
| No cap (limit unset)   | `nil`            | Stream until K8s says we're done. |
| Cap with budget        | `*N` where N > 0 | N items still allowed. |
| Cap exhausted          | `*0`             | Stop; do not call K8s again. |

A plain `int64` collapses the first and third states into `0`, which is
exactly the proto3 zero-value problem. The pointer is the cheapest
disambiguation that survives JSON round-trips through our continuation
token.

`effectiveLimit` panics on `*0` because it should never be invoked when
the cap is exhausted — the caller is required to stop first. This is a
programming-error guard, not user input validation.

### Patch resources

`XPatch` tokens are not separate API resources; they share the underlying
kind. `List` strips the `Patch` suffix from `gvk.Kind` and routes through
the underlying resource. There is a known corner case if K8s ever
introduces a real kind literally named `FooPatch`, tracked in the
follow-ups in `guinsnotes/list-rpc-followups.md`.

### Cluster-scoped kinds

The schema omits `namespace` from `listInputs` for known cluster-scoped
kinds. At runtime, if a `namespace` query arrives anyway (older clients,
hand-rolled gRPC), the provider rejects it with `InvalidArgument` rather
than silently ignoring it.

### Stale discovery cache

A freshly installed CRD will not yet be in the RESTMapper. On a mapping
miss, the provider resets the discovery cache and retries once. A second
miss surfaces as the underlying error.

### Empty results

K8s returns an empty list when nothing matches. The provider sends zero
`Result` messages and no `Continuation`, which the engine interprets as
"done."

## Non-Goals

### Helm Release List — future issue

`kubernetes:helm.sh/v3:Release` is not a K8s API resource; it wraps the
Helm Go client. Listing Releases is real user value (drift detection,
bulk import) but requires a separate code path through
`action.NewList(cfg).Run()`. Worth filing as its own issue.

### Helm Chart, kustomize, yaml overlays — category error

These are *generator* resources that expand client-side into K8s
manifests. They have no persistent cluster-side identity, so "list all
my Charts" doesn't have an obvious meaning. Users wanting inventory
should list the K8s resources the generators produced.

### `CustomResource` overlay List — duplicate path

Listing arbitrary CRD instances already works via the direct GVK token
(`kubernetes:cert-manager.io/v1:Certificate`) — runtime resolves through
the RESTMapper. Listing through the overlay would duplicate that path
with extra indirection.

### Property-level search

Pulumi's `List` returns importable IDs only — `id` and `name`. The provider
will not grow filtering by arbitrary resource properties; that role
belongs to higher-level tooling that calls `List` and then `Read`.
