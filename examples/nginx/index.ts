import * as kubernetes from "@pulumi/kubernetes";

var fs = require("fs");

kubernetes.config.clientCertificate = fs.readFileSync("./kube/cert.pem");
kubernetes.config.clientKey = fs.readFileSync("./kube/key.pem");
kubernetes.config.host = process.env.host;// "https://192.168.99.100:8443"
kubernetes.config.clusterCaCertificate = fs.readFileSync("./kube/ca.pem");

let mypodmetadata = {
        name: "nginx",
        labels: {
                app: "nginx"
        },
};

let nginxcontainer = new kubernetes.core.Pod("nginx", {
        metadata: [mypodmetadata],
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

let myvolumemetadata = {
        name: "nginxvolume"
};

let nginxvolume = new kubernetes.core.PersistentVolume("redis", {
        metadata: [myvolumemetadata],
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

let redispodmetadata = {
        name: "redis"
};

let redispod = new kubernetes.core.Pod("redis", {
        metadata: [redispodmetadata],
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
