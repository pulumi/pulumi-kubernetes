// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

// --------------------------------------------------------------------------
// MariaDB library.
// --------------------------------------------------------------------------

const appLabels = {app: "apache"};
const defaults = {
    name: "mariadb",
    namespace: "default",
    labels: appLabels,
    metricsEnabled: true,

    mariaConfig: {
        user: "",
        db: "",
        cnf: `[mysqld]
        innodb_buffer_pool_size=2G`,
    },

    mariaRootPassword: "password",
    mariaStorageClassName: "",
    mariaPvAccessMode: "ReadWriteOnce",
    mariaPvSize: "8Gi",

    image: "bitnami/mariadb:10.1.26-r2",
    imagePullPolicy: "IfNotPresent",
    serviceType: "ClusterIP",
    persistence: {
        accessMode: "ReadWriteOnce",
        size: "8Gi",
    },
    resources: {
        requests: {
            memory: "256Mi",
            cpu: "250m",
        },
    },
    metrics: {
        image: "prom/mysqld-exporter",
        imageTag: "v0.10.0",
        imagePullPolicy: "IfNotPresent",
        resources: {},
        annotations: {
            "prometheus.io/scrape": "true",
            "prometheus.io/port": "9104",
        },
    },
};

export function makeService(
    namespace: string, name: string, metricsEnabled=true, labels: object={app:name},
    selector: object={app:name},
): k8s.core.v1.Service {
    const ports = [
        {
            name: "mysql",
            port: 3306,
            targetPort: "mysql",
        },
    ];

    if (metricsEnabled) {
        ports.push({
            name: "metrics",
            port: 9104,
            targetPort: "metrics",
        })
    }

    return new k8s.core.v1.Service(
        name,
        {
            metadata: {
                name: name,
                labels: labels,
                ...metricsEnabled && {
                    annotations: {
                        "prometheus.io/scrape": "true",
                        "prometheus.io/port": "9104",
                    },
                },
            },
            spec: {
                type: "ClusterIP",
                ports: ports,
                selector: labels,
            },
        });
}

export function makeSecret(
    namespace: string, name: string, mariaRootPassword: string, labels: object={app:name},
): k8s.core.v1.Secret {
    return new k8s.core.v1.Secret(
        name,
        {
            metadata: {
                name: name,
                namespace: namespace,
                labels: labels,
            },
            type: "Opaque",
            data: {
                "mariadb-root-password": Buffer.from(mariaRootPassword).toString('base64'),
            },
        });
}

export function makeConfigMap(
    namespace:string , name: string, labels: object={app:name},
): k8s.core.v1.ConfigMap {
    return new k8s.core.v1.ConfigMap(
        name,
        {
            metadata: {
                name: name,
                namespace: namespace,
                labels: labels,
            },
            data: {
                "my.cnf": defaults.mariaConfig.cnf,
            },
        });
}

export function makePersistentVolumeClaim(
    namespace: string, name: string, storageClassName="", labels: object={app:name},
): k8s.core.v1.PersistentVolumeClaim {
    return new k8s.core.v1.PersistentVolumeClaim(
        name,
        {
            metadata: {
                name: name,
                namespace: namespace,
                labels: labels,
            },
            spec: {
                accessModes: [
                    defaults.mariaPvAccessMode,
                ],
                resources: {
                    requests: {
                        storage: defaults.mariaPvSize,
                    },
                },
                ...defaults.mariaStorageClassName != null && {
                    storageClassName: defaults.mariaStorageClassName
                },
            },
        });
}

namespace deployment {
    export function makePersistent(
        namespace: string, name: string, passwordSecretName: string, mariaConfig=defaults.mariaConfig,
        metricsEnabled=true, existingClaim=name, labels={app:name}, configMapName=name,
    ): k8s.apps.v1.Deployment {
        const volume = {
            name: "data",
            persistentVolumeClaim: {
                claimName: existingClaim
            }
        };
        const dep = base(namespace, name, passwordSecretName, mariaConfig, metricsEnabled,
            existingClaim, labels, configMapName);
            (<any>dep).spec.template.spec.volumes.push(volume);
            (<any>dep).spec.template.spec.containers.map((c: any) => {
                if (c.name == name) {
                    c.volumeMounts.push({
                        name: "data",
                        mountPath: "/bitnami/mariadb",
                    });
                }
                return c;
            })
            return dep;
    }

    export function makeNonPersistent(
        namespace: string, name: string, passwordSecretName: string,
        mariaConfig=defaults.mariaConfig, metricsEnabled=true, existingClaim=name,
        labels: object={app:name}, configMapName=name,
    ): k8s.apps.v1.Deployment {
        return base(namespace, name, passwordSecretName, mariaConfig, metricsEnabled, existingClaim,
            labels, configMapName);
    }

    const secure = (passwordSecretName: string) => [
        {
            name: "MARIADB_ROOT_PASSWORD",
            valueFrom: {
                secretKeyRef: {
                    name: passwordSecretName,
                    key: "mariadb-root-password",
                },
            },
        },
        {
            name: "MARIADB_PASSWORD",
            valueFrom: {
                secretKeyRef: {
                    name: passwordSecretName,
                    key: "mariadb-password",
                },
            },
        },
    ];

    const insecure = (passwordSecretName: string) => [
        {
            name: "ALLOW_EMPTY_PASSWORD",
            value: "yes",
        },
        {
            name: "MARIADB_PASSWORD",
            valueFrom: {
                secretKeyRef: {
                    name: passwordSecretName,
                    key: "mariadb-password",
                },
            },
        }
    ];

    const base = (
        namespace:string, name: string, passwordSecretName: string,
        mariaConfig: any, metricsEnabled: boolean, existingClaim: string,
        labels: object, configMapName: string,
    ): k8s.apps.v1.Deployment => {
        const metricsContainer =
        metricsEnabled == false
        ? []
        : [
            {
                name: "metrics",
                image: `${defaults.metrics.image}:${defaults.metrics.imageTag}`,
                imagePullPolicy: defaults.metrics.imagePullPolicy,
                env: [
                    {
                        name: "MARIADB_ROOT_PASSWORD",
                        valueFrom: {
                            secretKeyRef: {
                                name: name,
                                key: "mariadb-root-password",
                            },
                        },
                    },
                ],
                command: [ 'sh', '-c', 'DATA_SOURCE_NAME="root:$MARIADB_ROOT_PASSWORD@(localhost:3306)/" /bin/mysqld_exporter' ],
                ports: [
                    {
                        name: "metrics",
                        containerPort: 9104,
                    },
                ],
                livenessProbe: {
                    httpGet: {
                        path: "/metrics",
                        port: "metrics",
                    },
                    initialDelaySeconds: 15,
                    timeoutSeconds: 5,
                },
                readinessProbe: {
                    httpGet: {
                        path: "/metrics",
                        port: "metrics",
                    },
                    initialDelaySeconds: 5,
                    timeoutSeconds: 1,
                },
                resources: defaults.metrics.resources,
            },
        ];

        return new k8s.apps.v1.Deployment(
            name,
            {
                metadata: {
                    name: name,
                    namespace: namespace,
                    labels: labels,
                },
                spec: {
                    selector: {
                        matchLabels: labels,
                    },
                    template: {
                        metadata: {
                            namespace: namespace,
                            labels: labels,
                        },
                        spec: {
                            containers: [
                                {
                                    name: "mariadb",
                                    image: defaults.image,
                                    imagePullPolicy: defaults.imagePullPolicy,
                                    env: [
                                        ...secure(passwordSecretName),
                                        {
                                            name: "MARIADB_USER",
                                            value: mariaConfig.user
                                        },
                                        {
                                            name: "MARIADB_DATABASE",
                                            value:  mariaConfig.db
                                        },
                                    ],
                                    ports: [
                                        {
                                            name: "mysql",
                                            containerPort: 3306,
                                        },
                                    ],
                                    livenessProbe: {
                                        exec: {
                                            command: [
                                                "mysqladmin",
                                                "ping",
                                            ],
                                        },
                                        initialDelaySeconds: 30,
                                        timeoutSeconds: 5,
                                    },
                                    readinessProbe: {
                                        exec: {
                                            command: [
                                                "mysqladmin",
                                                "ping",
                                            ],
                                        },
                                        initialDelaySeconds: 5,
                                        timeoutSeconds: 1,
                                    },
                                    resources: defaults.resources,
                                    volumeMounts: [
                                        {
                                            name: "config",
                                            mountPath: "/bitnami/mariadb/conf/my_custom.cnf",
                                            subPath: "my.cnf",
                                        },
                                    ],
                                },
                                ...metricsContainer
                            ],
                            volumes: [
                                {
                                    name: "config",
                                    configMap: {
                                        name: configMapName,
                                    },
                                }
                            ],
                        },
                    },
                },
            });
    };
}

// --------------------------------------------------------------------------
// Example app.
// --------------------------------------------------------------------------

const cm = makeConfigMap("default", "mariadb-app");
const d = deployment.makeNonPersistent("default", "mariadb-app", "mariadb-app");
// const pvc = persistentVolumeClaim.make("default", "mariadb-app");
const s = makeSecret("default", "mariadb-app", "mariaRootPassword");
const svc = makeService("default", "mariadb-app")
