package await

// ------------------------------------------------------------------------------------------------

// Await logic for apps/v1beta1/StatefulSet, apps/v1beta2/StatefulSet,
// and apps/v1/StatefulSet.
//
// The goal of this code is to provide a fine-grained account of the status of a Kubernetes
// StatefulSet as it is being initialized. The idea is that if something goes wrong early, we want to
// alert the user so they can cancel the operation instead of waiting for timeout (~10 minutes).
//
// A StatefulSet is a construct that allows users to specify how to execute an update to a stateful
// application that is replicated some number of times in a cluster. When an application is updated,
// the StatefulSet will incrementally roll out the new version (according to the policy requested by
// the user). When the new application Pods becomes "live" (as specified by the liveness and
// readiness probes), the old Pods are killed and deleted.
//
// Because this resource abstracts over so much, the success conditions are fairly complex:
//
//   1. `.metadata.generation` in the current StatefulSet must have been incremented by the
//   	StatefulSet controller, i.e., it must not be equal to the generation number in the
//   	previous outputs.
//   2. `.status.updateRevision` matches `.status.currentRevision`.
//   3. `.status.currentReplicas` and `.status.readyReplicas` match the value of `.status.replicas`.
//
// The event loop depends on the following channels:
//
//   1. The StatefulSet channel, to which the Kubernetes API server will push every change
//      (additions, modifications, deletions) to any StatefulSet it knows about.
//   2. The Pod channel, which is the same idea as the StatefulSet channel, except it gets updates
//      to Pod objects. These are then aggregated and any errors are bundled together and
//      periodically reported to the user.
//   3. A timeout channel, which fires after some minutes.
//   4. A cancellation channel, with which the user can signal cancellation (e.g., using SIGINT).
//   5. A period channel, which is used to signal when we should display an aggregated report of
//      Pod errors we know about.
//
// The `statefulsetInitAwaiter` will synchronously process events from the union of all these
// channels. Any time the success conditions described above are reached, we will terminate
// the awaiter.
//
// The opportunity to display intermediate results will typically appear after a container in the
// Pod fails, (e.g., volume fails to mount, image fails to pull, exited with code 1, etc.).
//
// NB: Deleting a StatefulSet does not automatically delete any associated PersistentVolumes. We
//     may wish to address this case separately, but for now, PersistentVolumes are ignored by
//     the await logic. The await logic will still catch misconfiguration problems with
//     PersistentVolumeClaims because the related Pod will fail to go active, preventing success
//     condition 3 from being met.
//
// x-refs:
//   * https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/
//   * https://kubernetes.io/docs/tutorials/stateful-application/basic-stateful-set/
//   * https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#statefulset-v1-apps

// ------------------------------------------------------------------------------------------------

