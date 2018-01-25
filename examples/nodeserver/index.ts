import * as kubernetes from "@pulumi/kubernetes";

var fs = require("fs");

kubernetes.config.clientCertificate = fs.readFileSync("./kube/cert.pem");
kubernetes.config.clientKey = fs.readFileSync("./kube/key.pem");
kubernetes.config.host = process.env.host; //"https://192.168.99.100"
kubernetes.config.clusterCaCertificate = fs.readFileSync("./kube/ca.pem");

let mypodmetadata = {
        name: "helloworld-pod",
        labels: {
                app: "hello-world"
        },
};

let mypod = new kubernetes.core.Pod("helloworldpod", {
        metadata: [mypodmetadata],
        spec: [{
                container: [{
                        image: "hello-node:v1",
                        name: "helloworld-container",
                        port: [{
                                containerPort: 8080
                        }]
                }]
        }]
});

let service = new kubernetes.core.Service("helloworldservice", {
        metadata: [{
                name: "helloworld-service"
        }],
        spec: [{
                selector: {
                        app: mypodmetadata.labels.app
                },
                port: [{
                        port: 8080,
                        targetPort: 8080
                }],
                type: "LoadBalancer"
        }]
});
