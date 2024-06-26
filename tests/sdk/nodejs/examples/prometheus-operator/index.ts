import * as k8s from "@pulumi/kubernetes";
import * as pulumi from "@pulumi/pulumi";

// PrometheusOperatorArgs are the options to configure on the CoreOS
// PrometheusOperator.
interface PrometheusOperatorArgs {
    version?: string;
}

// PrometheusOperator implements the CoreOS Prometheus Operator.
class PrometheusOperator extends pulumi.ComponentResource {
    public readonly configFile: k8s.yaml.ConfigFile;
    public readonly service: pulumi.Output<k8s.core.v1.Service>;
    constructor(
        name: string,
        args: PrometheusOperatorArgs,
        opts?: pulumi.ComponentResourceOptions,
    ) {
        super('pulumi:monitoring/v1:PrometheusOperator', name, {}, opts);

        this.configFile = new k8s.yaml.ConfigFile(name, {
            file: 'https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/v0.63.0/bundle.yaml',
        }, {parent: this});

        this.service = this.configFile.getResource("v1/Service", "default", "prometheus-operator");
    }
}

// Create the Prometheus Operator.
const prometheusOperator = new PrometheusOperator("prometheus", {});

// Create the Prometheus Operator ServiceMonitor.
const myMonitoring = prometheusOperator.service.apply(service => {
    return new k8s.apiextensions.CustomResource('my-monitoring', {
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
                    // start with the following
                    relabelings: [
                        {
                            regex: '(.*)',
                            targetLabel: 'stackdriver',
                            replacement: 'true',
                            action: 'replace'
                        }
                    ],
                    // try to add the following in replacement of above in steps/step1.ts
                    // metricRelabelings: [
                    //   {
                    //     sourceLabels: ['__name__'],
                    //     regex: 'typhoon_(.*)',
                    //     targetLabel: 'stackdriver',
                    //     replacement: 'true',
                    //     action: 'replace'
                    //   }
                    // ]
                },
            ],
        },
    }, {dependsOn: service});
})
export const myMonitoringName = myMonitoring.id;
