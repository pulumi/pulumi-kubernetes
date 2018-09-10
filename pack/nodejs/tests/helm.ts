import * as assert from "assert";
import * as helm from "../helm";
const helmSort = helm.v2.helmSort;

function makeKinds(kinds: string[]): { kind: string }[] {
    return kinds.map(kind => {
        return { kind: kind };
    });
}

// Simple Fischer-Yates implementation. Taken from StackOverflow[1].
//
// [1]: https://stackoverflow.com/a/6274381
function shuffle<T>(a: T[]): T[] {
    for (let i = a.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [a[i], a[j]] = [a[j], a[i]];
    }
    return a;
}

describe("helmSort", () => {
    it("is the identity function for the empty array", () => {
        assert.deepEqual([].sort(helmSort), []);
    });

    it("is the identity function for single-element array", () => {
        assert.deepEqual(makeKinds(["Namespace"]).sort(helmSort), makeKinds(["Namespace"]));
    });

    it("correctly sorts duplicates", () => {
        assert.deepEqual(
            makeKinds(["Namespace", "LimitRange", "Namespace"]).sort(helmSort),
            makeKinds(["Namespace", "Namespace", "LimitRange"])
        );
    });

    it("puts unknown kinds last", () => {
        assert.deepEqual(
            makeKinds(["UNKNOWN KIND", "LimitRange", "Namespace"]).sort(helmSort),
            makeKinds(["Namespace", "LimitRange", "UNKNOWN KIND"])
        );
    });

    it("sorts a shuffled array", () => {
        const shuffled = makeKinds(
            shuffle([
                "Namespace",
                "ResourceQuota",
                "LimitRange",
                "PodSecurityPolicy",
                "Secret",
                "ConfigMap",
                "StorageClass",
                "PersistentVolume",
                "PersistentVolumeClaim",
                "ServiceAccount",
                "CustomResourceDefinition",
                "ClusterRole",
                "ClusterRoleBinding",
                "Role",
                "RoleBinding",
                "Service",
                "DaemonSet",
                "Pod",
                "ReplicationController",
                "ReplicaSet",
                "Deployment",
                "StatefulSet",
                "Job",
                "CronJob",
                "Ingress",
                "APIService",
                "UNKNOWN_KIND"
            ])
        );

        const sorted = makeKinds([
            "Namespace",
            "ResourceQuota",
            "LimitRange",
            "PodSecurityPolicy",
            "Secret",
            "ConfigMap",
            "StorageClass",
            "PersistentVolume",
            "PersistentVolumeClaim",
            "ServiceAccount",
            "CustomResourceDefinition",
            "ClusterRole",
            "ClusterRoleBinding",
            "Role",
            "RoleBinding",
            "Service",
            "DaemonSet",
            "Pod",
            "ReplicationController",
            "ReplicaSet",
            "Deployment",
            "StatefulSet",
            "Job",
            "CronJob",
            "Ingress",
            "APIService",
            "UNKNOWN_KIND"
        ]);

        assert.notDeepEqual(shuffled, sorted);
        assert.deepEqual(shuffled.sort(helmSort), sorted);
    });
});
