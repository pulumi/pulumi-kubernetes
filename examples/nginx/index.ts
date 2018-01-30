import * as kubernetes from "@pulumi/kubernetes";

kubernetes.config.host = process.env.host;

// Create an nginx pod
let nginxcontainer = new kubernetes.Pod("nginx", {
        metadata: [
                {
                        name: "nginx",
                        labels: {
                                app: "nginx"
                        }
                }
        ],
        spec: [{
                container: [{
                        image: "nginx:1.7.9",
                        name: "nginx",
                        port: [{
                                containerPort: 80
                        }]
                }]
        }]
});

// Create an nginxvolume
let nginxvolume = new kubernetes.PersistentVolume("redis", {
        metadata: [{
                name: "nginxvolume"
        }],
        spec: [{
                capacity: {
                        storage: "10Gi"
                },
                accessModes: ["ReadWriteMany"],
                persistentVolumeSource: [{
                        gcePersistentDisk: [{
                                pdName: "test-123"
                        }]
                }]
        }]
});

// create a redis pod
let redispod = new kubernetes.Pod("redis", {
        metadata: [{
                name: "redis"
        }],
        spec: [{
                container: [{
                        name: "redis",
                        image: "redis",
                        volumeMount: [{
                                name: "redis-persistent-storage",
                                mountPath: "/data/redis"
                        }]
                }],
                volume: [{
                        name: "redis-persistent-storage",
                        emptyDir: [{
                                medium: ""
                        }]
                }]
        }],
});
