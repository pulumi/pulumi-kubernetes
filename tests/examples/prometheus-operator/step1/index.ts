import * as k8s from "@pulumi/kubernetes";
import * as pulumi from "@pulumi/pulumi";

// PrometheusOperatorArgs are the options to configure on the CoreOS
// PrometheusOperator.
interface PrometheusOperatorArgs {
    namespace: pulumi.Input<string>;
    version?: string;
}

// PrometheusOperator implements the CoreOS Prometheus Operator.
export class PrometheusOperator extends pulumi.ComponentResource {
    public readonly configFile: k8s.yaml.ConfigFile;
    constructor(
        name: string,
        args: PrometheusOperatorArgs,
        opts?: pulumi.ComponentResourceOptions,
    ) {
        super('pulumi:monitoring/v1:PrometheusOperator', name, {}, opts);

        this.configFile = new k8s.yaml.ConfigFile(
            name,
            {
                file: `https://github.com/coreos/prometheus-operator/raw/release-${args.version || '0.38'}/bundle.yaml`,
                transformations: [
                    obj => {
                        if (obj.metadata.namespace) {
                            obj.metadata.namespace = args.namespace;
                        }
                        if (obj.kind === 'ClusterRoleBinding') {
                            obj.subjects[0].namespace = args.namespace;
                        }
                    },
                ],
            }, { parent: this });
    }
}

// Create the Prometheus Operator.
const prometheusOperator = new PrometheusOperator("prometheus", {
    namespace: "default",
});

// Create the Prometheus Operator ServiceMonitor.
const myMonitoring = new k8s.apiextensions.CustomResource('my-monitoring', {
    apiVersion: 'monitoring.coreos.com/v1',
    kind: 'ServiceMonitor',
    spec: {
        selector: {
            matchLabels: { app: 'my-app' },
        },
        endpoints: [
            {
                port: 'http',
                interval: '65s',
                // removing the following in index.ts in favor of below 
                // relabelings: [
                //   {
                //     regex: '(.*)',
                //     targetLabel: 'stackdriver',
                //     replacement: 'true',
                //     action: 'replace'
                //   }
                // ],
                // add the following in replacement of above in index.ts
                metricRelabelings: [
                    {
                        sourceLabels: ['__name__'],
                        regex: 'typhoon_(.*)',
                        targetLabel: 'stackdriver',
                        replacement: 'true',
                        action: 'replace'
                    }
                ]
            },
        ],
    },
}, {dependsOn: prometheusOperator});
