import * as k8s from "@pulumi/kubernetes";
import * as pulumi from "@pulumi/pulumi";
import * as random from "@pulumi/random";

const config = new pulumi.Config();

export interface AppConfig {
    namespace: pulumi.Input<string>;
    appName: string;

    metricsEnabled: pulumi.Input<boolean>;

    storageClassName?: pulumi.Input<string>;

    mariaConfig: pulumi.Input<{
        rootPassword: pulumi.Input<string>;
        user: pulumi.Input<string>;
        db: pulumi.Input<string>;
        password: pulumi.Input<string>;
        cnf: pulumi.Input<string>;
    }>;

    image: pulumi.Input<string>;
    imagePullPolicy: pulumi.Input<string>;
    serviceType: pulumi.Input<string>;

    persistence: pulumi.Input<{
        accessMode: pulumi.Input<string>;
        size: pulumi.Input<string>;
    }>;

    resources: pulumi.Input<{
        requests: pulumi.Input<{
            cpu: pulumi.Input<string>;
            memory: pulumi.Input<string>;
        }>;
    }>;

    metrics: pulumi.Input<{
        image: pulumi.Input<string>;
        imageTag: pulumi.Input<string>;
        imagePullPolicy: pulumi.Input<string>;
        resources: pulumi.Input<{}>;
        annotations: pulumi.Input<{
            "prometheus.io/scrape": "true";
            "prometheus.io/port": "9104";
        }>;
    }>;
}

export const appConfig: AppConfig = {
    namespace:
        config.get("namespace") ||
        new k8s.core.v1.Namespace("mariadb-ns").metadata.name,
    appName: config.get("appName") || "mariadb-app",

    metricsEnabled: config.getBoolean("metricsEnabled") || true,

    mariaConfig: {
        rootPassword:
            config.get("mariadbRootPassword") ||
            new random.RandomString("rootPassword", {
                length: 16,
                special: true
            }).result,
        user: config.get("mariadbUser") || "admin",
        db: config.get("mariadbTableName") || "admin-db",
        password:
            config.get("mariadbPassword") ||
            new random.RandomString("password", {
                length: 16,
                special: true
            }).result,
        cnf: `[mysqld]
        innodb_buffer_pool_size=2G`
    },

    image: config.get("image") || "bitnamilegacy/mariadb:12.0.2",
    imagePullPolicy: config.get("imagePullPolicy") || "IfNotPresent",
    serviceType: config.get("serviceType") || "ClusterIP",
    persistence: {
        accessMode: config.get("mariadbPvAccess") || "ReadWriteOnce",
        size: config.get("mariadbPvSize") || "8Gi"
    },
    resources: {
        requests: {
            memory: "256Mi",
            cpu: "250m"
        }
    },
    metrics: {
        image: "prom/mysqld-exporter",
        imageTag: "v0.10.0",
        imagePullPolicy: "IfNotPresent",
        resources: {},
        annotations: {
            "prometheus.io/scrape": "true",
            "prometheus.io/port": "9104"
        }
    }
};
