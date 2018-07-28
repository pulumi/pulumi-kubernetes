import { execSync } from "child_process";
import * as fs from "fs";

import * as k8s from "./index";
import * as pulumi from "@pulumi/pulumi";
import * as shell from "shell-quote";
import * as tmp from "tmp";
import * as yaml from "js-yaml";

export namespace v2 {
    export interface ChartOpts {
        repo: string;
        chart: string;
        version: string;

        values?: any;
        transformations?: ((o: any) => void)[];
        fetchOpts?: FetchOpts;
    }

    // Chart is a component representing a collection of resources described by an arbitrary Helm
    // Chart. The Chart can be fetched from any source that is accessible to the `helm` command
    // line. Values in the `values.yml` file can be overridden using `ChartOpts.values` (equivalent
    // to `--set` or having multiple `values.yml` files). Objects can be tranformed arbitrarily by
    // supplying callbacks to `ChartOpts.transformations`.
    //
    // `Chart` does not use Tiller. The Chart specified is copied and expanded locally; any values
    // that would be retrieved in-cluster would be assigned fake values, and none of Tiller's
    // server-side validity testing is executed.
    //
    // The semantics of `update` on a Chart are identical to those of Helm and kubectl; for example,
    // unlike a "normal" Pulumi program, updating a ConfigMap does not trigger a cascading update
    // among Deployments that reference it.
    export class Chart extends pulumi.ComponentResource {
        public readonly resources: {[key: string]: pulumi.CustomResource};

        constructor(releaseName: string, config: ChartOpts, opts?: pulumi.ResourceOptions) {
            super("kubernetes:helm.sh/v2:Chart", releaseName, config, opts);

            // Create temporary directories and files to hold chart data and override values.
            const overrides = tmp.fileSync({postfix: ".yaml"});
            const chartDir = tmp.dirSync({unsafeCleanup: true});

            try {
                // Fetch chart.
                fetch(`${config.repo}/${config.chart}`,
                    {destination: chartDir.name, version: config.version});

                // Write overrides file.
                const data = JSON.stringify(config.values || {}, undefined, "  ");
                fs.writeFileSync(overrides.name, data);

                // Does not require Tiller. From the `helm template` documentation:
                //
                // >  Render chart templates locally and display the output.
                // >
                // > This does not require Tiller. However, any values that would normally be
                // > looked up or retrieved in-cluster will be faked locally. Additionally, none
                // > of the server-side testing of chart validity (e.g. whether an API is supported)
                // > is done.
                const chart = `${shell.quote([chartDir.name])}/${shell.quote([config.chart])}`;
                const release = shell.quote([releaseName]);
                const values = shell.quote([overrides.name]);
                const yamlStream = execSync(
                    `helm template ${chart} --name ${release} --values ${values}`
                ).toString();
                const resourcesObjects = yaml.safeLoadAll(yamlStream);
                this.resources = fromList(resourcesObjects, config.transformations || [], { ...opts, parent: this });
            } finally {
                // Clean up temporary files and directories.
                chartDir.removeCallback()
                overrides.removeCallback()
            }
        }
    }
}

export interface FetchOpts {
    // Specific version of a chart. Without this, the latest version is fetched.
    version?: string;

    // Verify certificates of HTTPS-enabled servers using this CA bundle.
    caFile?: string;

    // Identify HTTPS client using this SSL certificate file.
    certFile?: string;

    // Identify HTTPS client using this SSL key file.
    keyFile?: string;

    // Location to write the chart. If this and tardir are specified, tardir is appended to this
    // (default ".").
    destination?: string;

    // Keyring containing public keys (default "/Users/alex/.gnupg/pubring.gpg").
    keyring?: string;

    // Chart repository password.
    password?: string;

    // Chart repository url where to locate the requested chart.
    repo?: string;

    // If untar is specified, this flag specifies the name of the directory into which the chart is
    // expanded (default ".").
    untardir?: string;

    // Chart repository username.
    username?: string;

    // Location of your Helm config. Overrides $HELM_HOME (default "/Users/alex/.helm").
    home?: string;

    // Use development versions, too. Equivalent to version '>0.0.0-0'. If --version is set, this is
    // ignored.
    devel?: boolean;

    // Fetch the provenance file, but don't perform verification.
    prov?: boolean;

    // If set to false, will leave the chart as a tarball after downloading.
    untar?: boolean;

    // Verify the package against its signature.
    verify?: boolean;
}

// Retrieve a package from a package repository, and download it locally.
//
// This is useful for fetching packages to inspect, modify, or repackage. It can also be used to
// perform cryptographic verification of a chart without installing the chart.
//
// There are options for unpacking the chart after download. This will create a directory for the
// chart and uncompress into that directory.
//
// If the `verify` option is specified, the requested chart MUST have a provenance file, and MUST
// pass the verification process. Failure in any part of this will result in an error, and the chart
// will not be saved locally.
export function fetch(chart: string, opts?: FetchOpts) {
    const flags: string[] = [];
    if (opts !== undefined) {
        // Untar by default.
        if(opts.untar !== false) { flags.push(`--untar`); }

        if (opts.version !== undefined)     { flags.push(`--version ${shell.quote([opts.version])}`);         }
        if (opts.caFile !== undefined)      { flags.push(`--ca-file ${shell.quote([opts.caFile])}`);          }
        if (opts.certFile !== undefined)    { flags.push(`--cert-file ${shell.quote([opts.certFile])}`);      }
        if (opts.keyFile !== undefined)     { flags.push(`--key-file ${shell.quote([opts.keyFile])}`);        }
        if (opts.destination !== undefined) { flags.push(`--destination ${shell.quote([opts.destination])}`); }
        if (opts.keyring !== undefined)     { flags.push(`--keyring ${shell.quote([opts.keyring])}`);         }
        if (opts.password !== undefined)    { flags.push(`--password ${shell.quote([opts.password])}`);       }
        if (opts.repo !== undefined)        { flags.push(`--repo ${shell.quote([opts.repo])}`);               }
        if (opts.untardir !== undefined)    { flags.push(`--untardir ${shell.quote([opts.untardir])}`);       }
        if (opts.username !== undefined)    { flags.push(`--username ${shell.quote([opts.username])}`);       }
        if (opts.home !== undefined)        { flags.push(`--home ${shell.quote([opts.home])}`);               }
        if (opts.devel === true)            { flags.push(`--devel`);                                          }
        if (opts.prov === true)             { flags.push(`--prov`);                                           }
        if (opts.verify === true)           { flags.push(`--verify`);                                         }
    }
    execSync(`helm fetch ${shell.quote([chart])} ${flags.join(" ")}`);
}


function fromList(
    objs: any[], transforms?: ((o: any) => void)[], opts?: pulumi.ResourceOptions,
): {[key: string]: pulumi.CustomResource} {
    const resources: {[key: string]: pulumi.CustomResource} = {};
    for (const obj of objs) {
        if (obj == null) {
            continue;
        }

        for (const t of transforms || []) {
            t(obj);
        }

        const kind = obj["kind"];
        const apiVersion = obj["apiVersion"];
        if (kind == null || apiVersion == null) {
            continue;
        }
        const namespace = obj["metadata"]["namespace"] || "default";
        const name = obj["metadata"]["name"];
        switch (`${apiVersion}/${kind}`) {
            case "admissionregistration.k8s.io/v1alpha1/InitializerConfiguration":
                resources[`admissionregistration.k8s.io/v1alpha1/InitializerConfiguration::${namespace}/${name}`] =
                    new k8s.admissionregistration.v1alpha1.InitializerConfiguration(name, obj, opts);
                break;
            case "admissionregistration.k8s.io/v1alpha1/InitializerConfigurationList":
                resources[`admissionregistration.k8s.io/v1alpha1/InitializerConfigurationList::${namespace}/${name}`] =
                    new k8s.admissionregistration.v1alpha1.InitializerConfigurationList(name, obj, opts);
                break;
            case "admissionregistration.k8s.io/v1beta1/MutatingWebhookConfiguration":
                resources[`admissionregistration.k8s.io/v1beta1/MutatingWebhookConfiguration::${namespace}/${name}`] =
                    new k8s.admissionregistration.v1beta1.MutatingWebhookConfiguration(name, obj, opts);
                break;
            case "admissionregistration.k8s.io/v1beta1/MutatingWebhookConfigurationList":
                resources[`admissionregistration.k8s.io/v1beta1/MutatingWebhookConfigurationList::${namespace}/${name}`] =
                    new k8s.admissionregistration.v1beta1.MutatingWebhookConfigurationList(name, obj, opts);
                break;
            case "admissionregistration.k8s.io/v1beta1/ValidatingWebhookConfiguration":
                resources[`admissionregistration.k8s.io/v1beta1/ValidatingWebhookConfiguration::${namespace}/${name}`] =
                    new k8s.admissionregistration.v1beta1.ValidatingWebhookConfiguration(name, obj, opts);
                break;
            case "admissionregistration.k8s.io/v1beta1/ValidatingWebhookConfigurationList":
                resources[`admissionregistration.k8s.io/v1beta1/ValidatingWebhookConfigurationList::${namespace}/${name}`] =
                    new k8s.admissionregistration.v1beta1.ValidatingWebhookConfigurationList(name, obj, opts);
                break;
            case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition":
                resources[`apiextensions.k8s.io/v1beta1/CustomResourceDefinition::${namespace}/${name}`] =
                    new k8s.apiextensions.v1beta1.CustomResourceDefinition(name, obj, opts);
                break;
            case "apiextensions.k8s.io/v1beta1/CustomResourceDefinitionList":
                resources[`apiextensions.k8s.io/v1beta1/CustomResourceDefinitionList::${namespace}/${name}`] =
                    new k8s.apiextensions.v1beta1.CustomResourceDefinitionList(name, obj, opts);
                break;
            case "apiregistration/v1beta1/APIService":
                resources[`apiregistration/v1beta1/APIService::${namespace}/${name}`] =
                    new k8s.apiregistration.v1beta1.APIService(name, obj, opts);
                break;
            case "apiregistration/v1beta1/APIServiceList":
                resources[`apiregistration/v1beta1/APIServiceList::${namespace}/${name}`] =
                    new k8s.apiregistration.v1beta1.APIServiceList(name, obj, opts);
                break;
            case "apps/v1/ControllerRevision":
                resources[`apps/v1/ControllerRevision::${namespace}/${name}`] =
                    new k8s.apps.v1.ControllerRevision(name, obj, opts);
                break;
            case "apps/v1/ControllerRevisionList":
                resources[`apps/v1/ControllerRevisionList::${namespace}/${name}`] =
                    new k8s.apps.v1.ControllerRevisionList(name, obj, opts);
                break;
            case "apps/v1/DaemonSet":
                resources[`apps/v1/DaemonSet::${namespace}/${name}`] =
                    new k8s.apps.v1.DaemonSet(name, obj, opts);
                break;
            case "apps/v1/DaemonSetList":
                resources[`apps/v1/DaemonSetList::${namespace}/${name}`] =
                    new k8s.apps.v1.DaemonSetList(name, obj, opts);
                break;
            case "apps/v1/Deployment":
                resources[`apps/v1/Deployment::${namespace}/${name}`] =
                    new k8s.apps.v1.Deployment(name, obj, opts);
                break;
            case "apps/v1/DeploymentList":
                resources[`apps/v1/DeploymentList::${namespace}/${name}`] =
                    new k8s.apps.v1.DeploymentList(name, obj, opts);
                break;
            case "apps/v1/ReplicaSet":
                resources[`apps/v1/ReplicaSet::${namespace}/${name}`] =
                    new k8s.apps.v1.ReplicaSet(name, obj, opts);
                break;
            case "apps/v1/ReplicaSetList":
                resources[`apps/v1/ReplicaSetList::${namespace}/${name}`] =
                    new k8s.apps.v1.ReplicaSetList(name, obj, opts);
                break;
            case "apps/v1/StatefulSet":
                resources[`apps/v1/StatefulSet::${namespace}/${name}`] =
                    new k8s.apps.v1.StatefulSet(name, obj, opts);
                break;
            case "apps/v1/StatefulSetList":
                resources[`apps/v1/StatefulSetList::${namespace}/${name}`] =
                    new k8s.apps.v1.StatefulSetList(name, obj, opts);
                break;
            case "apps/v1beta1/ControllerRevision":
                resources[`apps/v1beta1/ControllerRevision::${namespace}/${name}`] =
                    new k8s.apps.v1beta1.ControllerRevision(name, obj, opts);
                break;
            case "apps/v1beta1/ControllerRevisionList":
                resources[`apps/v1beta1/ControllerRevisionList::${namespace}/${name}`] =
                    new k8s.apps.v1beta1.ControllerRevisionList(name, obj, opts);
                break;
            case "apps/v1beta1/Deployment":
                resources[`apps/v1beta1/Deployment::${namespace}/${name}`] =
                    new k8s.apps.v1beta1.Deployment(name, obj, opts);
                break;
            case "apps/v1beta1/DeploymentList":
                resources[`apps/v1beta1/DeploymentList::${namespace}/${name}`] =
                    new k8s.apps.v1beta1.DeploymentList(name, obj, opts);
                break;
            case "apps/v1beta1/DeploymentRollback":
                resources[`apps/v1beta1/DeploymentRollback::${namespace}/${name}`] =
                    new k8s.apps.v1beta1.DeploymentRollback(name, obj, opts);
                break;
            case "apps/v1beta1/Scale":
                resources[`apps/v1beta1/Scale::${namespace}/${name}`] =
                    new k8s.apps.v1beta1.Scale(name, obj, opts);
                break;
            case "apps/v1beta1/StatefulSet":
                resources[`apps/v1beta1/StatefulSet::${namespace}/${name}`] =
                    new k8s.apps.v1beta1.StatefulSet(name, obj, opts);
                break;
            case "apps/v1beta1/StatefulSetList":
                resources[`apps/v1beta1/StatefulSetList::${namespace}/${name}`] =
                    new k8s.apps.v1beta1.StatefulSetList(name, obj, opts);
                break;
            case "apps/v1beta2/ControllerRevision":
                resources[`apps/v1beta2/ControllerRevision::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.ControllerRevision(name, obj, opts);
                break;
            case "apps/v1beta2/ControllerRevisionList":
                resources[`apps/v1beta2/ControllerRevisionList::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.ControllerRevisionList(name, obj, opts);
                break;
            case "apps/v1beta2/DaemonSet":
                resources[`apps/v1beta2/DaemonSet::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.DaemonSet(name, obj, opts);
                break;
            case "apps/v1beta2/DaemonSetList":
                resources[`apps/v1beta2/DaemonSetList::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.DaemonSetList(name, obj, opts);
                break;
            case "apps/v1beta2/Deployment":
                resources[`apps/v1beta2/Deployment::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.Deployment(name, obj, opts);
                break;
            case "apps/v1beta2/DeploymentList":
                resources[`apps/v1beta2/DeploymentList::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.DeploymentList(name, obj, opts);
                break;
            case "apps/v1beta2/ReplicaSet":
                resources[`apps/v1beta2/ReplicaSet::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.ReplicaSet(name, obj, opts);
                break;
            case "apps/v1beta2/ReplicaSetList":
                resources[`apps/v1beta2/ReplicaSetList::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.ReplicaSetList(name, obj, opts);
                break;
            case "apps/v1beta2/Scale":
                resources[`apps/v1beta2/Scale::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.Scale(name, obj, opts);
                break;
            case "apps/v1beta2/StatefulSet":
                resources[`apps/v1beta2/StatefulSet::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.StatefulSet(name, obj, opts);
                break;
            case "apps/v1beta2/StatefulSetList":
                resources[`apps/v1beta2/StatefulSetList::${namespace}/${name}`] =
                    new k8s.apps.v1beta2.StatefulSetList(name, obj, opts);
                break;
            case "authentication.k8s.io/v1/TokenReview":
                resources[`authentication.k8s.io/v1/TokenReview::${namespace}/${name}`] =
                    new k8s.authentication.v1.TokenReview(name, obj, opts);
                break;
            case "authentication.k8s.io/v1beta1/TokenReview":
                resources[`authentication.k8s.io/v1beta1/TokenReview::${namespace}/${name}`] =
                    new k8s.authentication.v1beta1.TokenReview(name, obj, opts);
                break;
            case "authorization.k8s.io/v1/LocalSubjectAccessReview":
                resources[`authorization.k8s.io/v1/LocalSubjectAccessReview::${namespace}/${name}`] =
                    new k8s.authorization.v1.LocalSubjectAccessReview(name, obj, opts);
                break;
            case "authorization.k8s.io/v1/SelfSubjectAccessReview":
                resources[`authorization.k8s.io/v1/SelfSubjectAccessReview::${namespace}/${name}`] =
                    new k8s.authorization.v1.SelfSubjectAccessReview(name, obj, opts);
                break;
            case "authorization.k8s.io/v1/SelfSubjectRulesReview":
                resources[`authorization.k8s.io/v1/SelfSubjectRulesReview::${namespace}/${name}`] =
                    new k8s.authorization.v1.SelfSubjectRulesReview(name, obj, opts);
                break;
            case "authorization.k8s.io/v1/SubjectAccessReview":
                resources[`authorization.k8s.io/v1/SubjectAccessReview::${namespace}/${name}`] =
                    new k8s.authorization.v1.SubjectAccessReview(name, obj, opts);
                break;
            case "authorization.k8s.io/v1beta1/LocalSubjectAccessReview":
                resources[`authorization.k8s.io/v1beta1/LocalSubjectAccessReview::${namespace}/${name}`] =
                    new k8s.authorization.v1beta1.LocalSubjectAccessReview(name, obj, opts);
                break;
            case "authorization.k8s.io/v1beta1/SelfSubjectAccessReview":
                resources[`authorization.k8s.io/v1beta1/SelfSubjectAccessReview::${namespace}/${name}`] =
                    new k8s.authorization.v1beta1.SelfSubjectAccessReview(name, obj, opts);
                break;
            case "authorization.k8s.io/v1beta1/SelfSubjectRulesReview":
                resources[`authorization.k8s.io/v1beta1/SelfSubjectRulesReview::${namespace}/${name}`] =
                    new k8s.authorization.v1beta1.SelfSubjectRulesReview(name, obj, opts);
                break;
            case "authorization.k8s.io/v1beta1/SubjectAccessReview":
                resources[`authorization.k8s.io/v1beta1/SubjectAccessReview::${namespace}/${name}`] =
                    new k8s.authorization.v1beta1.SubjectAccessReview(name, obj, opts);
                break;
            case "autoscaling/v1/CrossVersionObjectReference":
                resources[`autoscaling/v1/CrossVersionObjectReference::${namespace}/${name}`] =
                    new k8s.autoscaling.v1.CrossVersionObjectReference(name, obj, opts);
                break;
            case "autoscaling/v1/HorizontalPodAutoscaler":
                resources[`autoscaling/v1/HorizontalPodAutoscaler::${namespace}/${name}`] =
                    new k8s.autoscaling.v1.HorizontalPodAutoscaler(name, obj, opts);
                break;
            case "autoscaling/v1/HorizontalPodAutoscalerList":
                resources[`autoscaling/v1/HorizontalPodAutoscalerList::${namespace}/${name}`] =
                    new k8s.autoscaling.v1.HorizontalPodAutoscalerList(name, obj, opts);
                break;
            case "autoscaling/v1/Scale":
                resources[`autoscaling/v1/Scale::${namespace}/${name}`] =
                    new k8s.autoscaling.v1.Scale(name, obj, opts);
                break;
            case "autoscaling/v2beta1/CrossVersionObjectReference":
                resources[`autoscaling/v2beta1/CrossVersionObjectReference::${namespace}/${name}`] =
                    new k8s.autoscaling.v2beta1.CrossVersionObjectReference(name, obj, opts);
                break;
            case "autoscaling/v2beta1/HorizontalPodAutoscaler":
                resources[`autoscaling/v2beta1/HorizontalPodAutoscaler::${namespace}/${name}`] =
                    new k8s.autoscaling.v2beta1.HorizontalPodAutoscaler(name, obj, opts);
                break;
            case "autoscaling/v2beta1/HorizontalPodAutoscalerList":
                resources[`autoscaling/v2beta1/HorizontalPodAutoscalerList::${namespace}/${name}`] =
                    new k8s.autoscaling.v2beta1.HorizontalPodAutoscalerList(name, obj, opts);
                break;
            case "batch/v1/Job":
                resources[`batch/v1/Job::${namespace}/${name}`] =
                    new k8s.batch.v1.Job(name, obj, opts);
                break;
            case "batch/v1/JobList":
                resources[`batch/v1/JobList::${namespace}/${name}`] =
                    new k8s.batch.v1.JobList(name, obj, opts);
                break;
            case "batch/v1beta1/CronJob":
                resources[`batch/v1beta1/CronJob::${namespace}/${name}`] =
                    new k8s.batch.v1beta1.CronJob(name, obj, opts);
                break;
            case "batch/v1beta1/CronJobList":
                resources[`batch/v1beta1/CronJobList::${namespace}/${name}`] =
                    new k8s.batch.v1beta1.CronJobList(name, obj, opts);
                break;
            case "batch/v2alpha1/CronJob":
                resources[`batch/v2alpha1/CronJob::${namespace}/${name}`] =
                    new k8s.batch.v2alpha1.CronJob(name, obj, opts);
                break;
            case "batch/v2alpha1/CronJobList":
                resources[`batch/v2alpha1/CronJobList::${namespace}/${name}`] =
                    new k8s.batch.v2alpha1.CronJobList(name, obj, opts);
                break;
            case "certificates.k8s.io/v1beta1/CertificateSigningRequest":
                resources[`certificates.k8s.io/v1beta1/CertificateSigningRequest::${namespace}/${name}`] =
                    new k8s.certificates.v1beta1.CertificateSigningRequest(name, obj, opts);
                break;
            case "certificates.k8s.io/v1beta1/CertificateSigningRequestList":
                resources[`certificates.k8s.io/v1beta1/CertificateSigningRequestList::${namespace}/${name}`] =
                    new k8s.certificates.v1beta1.CertificateSigningRequestList(name, obj, opts);
                break;
            case "v1/Binding":
                resources[`v1/Binding::${namespace}/${name}`] =
                    new k8s.core.v1.Binding(name, obj, opts);
                break;
            case "v1/ComponentStatus":
                resources[`v1/ComponentStatus::${namespace}/${name}`] =
                    new k8s.core.v1.ComponentStatus(name, obj, opts);
                break;
            case "v1/ComponentStatusList":
                resources[`v1/ComponentStatusList::${namespace}/${name}`] =
                    new k8s.core.v1.ComponentStatusList(name, obj, opts);
                break;
            case "v1/ConfigMap":
                resources[`v1/ConfigMap::${namespace}/${name}`] =
                    new k8s.core.v1.ConfigMap(name, obj, opts);
                break;
            case "v1/ConfigMapList":
                resources[`v1/ConfigMapList::${namespace}/${name}`] =
                    new k8s.core.v1.ConfigMapList(name, obj, opts);
                break;
            case "v1/Endpoints":
                resources[`v1/Endpoints::${namespace}/${name}`] =
                    new k8s.core.v1.Endpoints(name, obj, opts);
                break;
            case "v1/EndpointsList":
                resources[`v1/EndpointsList::${namespace}/${name}`] =
                    new k8s.core.v1.EndpointsList(name, obj, opts);
                break;
            case "v1/Event":
                resources[`v1/Event::${namespace}/${name}`] =
                    new k8s.core.v1.Event(name, obj, opts);
                break;
            case "v1/EventList":
                resources[`v1/EventList::${namespace}/${name}`] =
                    new k8s.core.v1.EventList(name, obj, opts);
                break;
            case "v1/LimitRange":
                resources[`v1/LimitRange::${namespace}/${name}`] =
                    new k8s.core.v1.LimitRange(name, obj, opts);
                break;
            case "v1/LimitRangeList":
                resources[`v1/LimitRangeList::${namespace}/${name}`] =
                    new k8s.core.v1.LimitRangeList(name, obj, opts);
                break;
            case "v1/Namespace":
                resources[`v1/Namespace::${namespace}/${name}`] =
                    new k8s.core.v1.Namespace(name, obj, opts);
                break;
            case "v1/NamespaceList":
                resources[`v1/NamespaceList::${namespace}/${name}`] =
                    new k8s.core.v1.NamespaceList(name, obj, opts);
                break;
            case "v1/Node":
                resources[`v1/Node::${namespace}/${name}`] =
                    new k8s.core.v1.Node(name, obj, opts);
                break;
            case "v1/NodeConfigSource":
                resources[`v1/NodeConfigSource::${namespace}/${name}`] =
                    new k8s.core.v1.NodeConfigSource(name, obj, opts);
                break;
            case "v1/NodeList":
                resources[`v1/NodeList::${namespace}/${name}`] =
                    new k8s.core.v1.NodeList(name, obj, opts);
                break;
            case "core/v1/ObjectReference":
                resources[`core/v1/ObjectReference::${namespace}/${name}`] =
                    new k8s.core.v1.ObjectReference(name, obj, opts);
                break;
            case "v1/PersistentVolume":
                resources[`v1/PersistentVolume::${namespace}/${name}`] =
                    new k8s.core.v1.PersistentVolume(name, obj, opts);
                break;
            case "v1/PersistentVolumeClaim":
                resources[`v1/PersistentVolumeClaim::${namespace}/${name}`] =
                    new k8s.core.v1.PersistentVolumeClaim(name, obj, opts);
                break;
            case "v1/PersistentVolumeClaimList":
                resources[`v1/PersistentVolumeClaimList::${namespace}/${name}`] =
                    new k8s.core.v1.PersistentVolumeClaimList(name, obj, opts);
                break;
            case "v1/PersistentVolumeList":
                resources[`v1/PersistentVolumeList::${namespace}/${name}`] =
                    new k8s.core.v1.PersistentVolumeList(name, obj, opts);
                break;
            case "v1/Pod":
                resources[`v1/Pod::${namespace}/${name}`] =
                    new k8s.core.v1.Pod(name, obj, opts);
                break;
            case "v1/PodList":
                resources[`v1/PodList::${namespace}/${name}`] =
                    new k8s.core.v1.PodList(name, obj, opts);
                break;
            case "v1/PodTemplate":
                resources[`v1/PodTemplate::${namespace}/${name}`] =
                    new k8s.core.v1.PodTemplate(name, obj, opts);
                break;
            case "v1/PodTemplateList":
                resources[`v1/PodTemplateList::${namespace}/${name}`] =
                    new k8s.core.v1.PodTemplateList(name, obj, opts);
                break;
            case "v1/ReplicationController":
                resources[`v1/ReplicationController::${namespace}/${name}`] =
                    new k8s.core.v1.ReplicationController(name, obj, opts);
                break;
            case "v1/ReplicationControllerList":
                resources[`v1/ReplicationControllerList::${namespace}/${name}`] =
                    new k8s.core.v1.ReplicationControllerList(name, obj, opts);
                break;
            case "v1/ResourceQuota":
                resources[`v1/ResourceQuota::${namespace}/${name}`] =
                    new k8s.core.v1.ResourceQuota(name, obj, opts);
                break;
            case "v1/ResourceQuotaList":
                resources[`v1/ResourceQuotaList::${namespace}/${name}`] =
                    new k8s.core.v1.ResourceQuotaList(name, obj, opts);
                break;
            case "v1/Secret":
                resources[`v1/Secret::${namespace}/${name}`] =
                    new k8s.core.v1.Secret(name, obj, opts);
                break;
            case "v1/SecretList":
                resources[`v1/SecretList::${namespace}/${name}`] =
                    new k8s.core.v1.SecretList(name, obj, opts);
                break;
            case "v1/Service":
                resources[`v1/Service::${namespace}/${name}`] =
                    new k8s.core.v1.Service(name, obj, opts);
                break;
            case "v1/ServiceAccount":
                resources[`v1/ServiceAccount::${namespace}/${name}`] =
                    new k8s.core.v1.ServiceAccount(name, obj, opts);
                break;
            case "v1/ServiceAccountList":
                resources[`v1/ServiceAccountList::${namespace}/${name}`] =
                    new k8s.core.v1.ServiceAccountList(name, obj, opts);
                break;
            case "v1/ServiceList":
                resources[`v1/ServiceList::${namespace}/${name}`] =
                    new k8s.core.v1.ServiceList(name, obj, opts);
                break;
            case "events.k8s.io/v1beta1/Event":
                resources[`events.k8s.io/v1beta1/Event::${namespace}/${name}`] =
                    new k8s.events.v1beta1.Event(name, obj, opts);
                break;
            case "events.k8s.io/v1beta1/EventList":
                resources[`events.k8s.io/v1beta1/EventList::${namespace}/${name}`] =
                    new k8s.events.v1beta1.EventList(name, obj, opts);
                break;
            case "extensions/v1beta1/DaemonSet":
                resources[`extensions/v1beta1/DaemonSet::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.DaemonSet(name, obj, opts);
                break;
            case "extensions/v1beta1/DaemonSetList":
                resources[`extensions/v1beta1/DaemonSetList::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.DaemonSetList(name, obj, opts);
                break;
            case "extensions/v1beta1/Deployment":
                resources[`extensions/v1beta1/Deployment::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.Deployment(name, obj, opts);
                break;
            case "extensions/v1beta1/DeploymentList":
                resources[`extensions/v1beta1/DeploymentList::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.DeploymentList(name, obj, opts);
                break;
            case "extensions/v1beta1/DeploymentRollback":
                resources[`extensions/v1beta1/DeploymentRollback::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.DeploymentRollback(name, obj, opts);
                break;
            case "extensions/v1beta1/Ingress":
                resources[`extensions/v1beta1/Ingress::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.Ingress(name, obj, opts);
                break;
            case "extensions/v1beta1/IngressList":
                resources[`extensions/v1beta1/IngressList::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.IngressList(name, obj, opts);
                break;
            case "extensions/v1beta1/NetworkPolicy":
                resources[`extensions/v1beta1/NetworkPolicy::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.NetworkPolicy(name, obj, opts);
                break;
            case "extensions/v1beta1/NetworkPolicyList":
                resources[`extensions/v1beta1/NetworkPolicyList::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.NetworkPolicyList(name, obj, opts);
                break;
            case "extensions/v1beta1/PodSecurityPolicy":
                resources[`extensions/v1beta1/PodSecurityPolicy::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.PodSecurityPolicy(name, obj, opts);
                break;
            case "extensions/v1beta1/PodSecurityPolicyList":
                resources[`extensions/v1beta1/PodSecurityPolicyList::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.PodSecurityPolicyList(name, obj, opts);
                break;
            case "extensions/v1beta1/ReplicaSet":
                resources[`extensions/v1beta1/ReplicaSet::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.ReplicaSet(name, obj, opts);
                break;
            case "extensions/v1beta1/ReplicaSetList":
                resources[`extensions/v1beta1/ReplicaSetList::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.ReplicaSetList(name, obj, opts);
                break;
            case "extensions/v1beta1/Scale":
                resources[`extensions/v1beta1/Scale::${namespace}/${name}`] =
                    new k8s.extensions.v1beta1.Scale(name, obj, opts);
                break;
            case "v1/APIGroup":
                resources[`v1/APIGroup::${namespace}/${name}`] =
                    new k8s.meta.v1.APIGroup(name, obj, opts);
                break;
            case "v1/APIGroupList":
                resources[`v1/APIGroupList::${namespace}/${name}`] =
                    new k8s.meta.v1.APIGroupList(name, obj, opts);
                break;
            case "v1/APIResourceList":
                resources[`v1/APIResourceList::${namespace}/${name}`] =
                    new k8s.meta.v1.APIResourceList(name, obj, opts);
                break;
            case "v1/APIVersions":
                resources[`v1/APIVersions::${namespace}/${name}`] =
                    new k8s.meta.v1.APIVersions(name, obj, opts);
                break;
            case "v1/DeleteOptions":
                resources[`v1/DeleteOptions::${namespace}/${name}`] =
                    new k8s.meta.v1.DeleteOptions(name, obj, opts);
                break;
            case "meta/v1/OwnerReference":
                resources[`meta/v1/OwnerReference::${namespace}/${name}`] =
                    new k8s.meta.v1.OwnerReference(name, obj, opts);
                break;
            case "v1/Status":
                resources[`v1/Status::${namespace}/${name}`] =
                    new k8s.meta.v1.Status(name, obj, opts);
                break;
            case "networking.k8s.io/v1/NetworkPolicy":
                resources[`networking.k8s.io/v1/NetworkPolicy::${namespace}/${name}`] =
                    new k8s.networking.v1.NetworkPolicy(name, obj, opts);
                break;
            case "networking.k8s.io/v1/NetworkPolicyList":
                resources[`networking.k8s.io/v1/NetworkPolicyList::${namespace}/${name}`] =
                    new k8s.networking.v1.NetworkPolicyList(name, obj, opts);
                break;
            case "policy/v1beta1/Eviction":
                resources[`policy/v1beta1/Eviction::${namespace}/${name}`] =
                    new k8s.policy.v1beta1.Eviction(name, obj, opts);
                break;
            case "policy/v1beta1/PodDisruptionBudget":
                resources[`policy/v1beta1/PodDisruptionBudget::${namespace}/${name}`] =
                    new k8s.policy.v1beta1.PodDisruptionBudget(name, obj, opts);
                break;
            case "policy/v1beta1/PodDisruptionBudgetList":
                resources[`policy/v1beta1/PodDisruptionBudgetList::${namespace}/${name}`] =
                    new k8s.policy.v1beta1.PodDisruptionBudgetList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1/ClusterRole":
                resources[`rbac.authorization.k8s.io/v1/ClusterRole::${namespace}/${name}`] =
                    new k8s.rbac.v1.ClusterRole(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1/ClusterRoleBinding":
                resources[`rbac.authorization.k8s.io/v1/ClusterRoleBinding::${namespace}/${name}`] =
                    new k8s.rbac.v1.ClusterRoleBinding(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1/ClusterRoleBindingList":
                resources[`rbac.authorization.k8s.io/v1/ClusterRoleBindingList::${namespace}/${name}`] =
                    new k8s.rbac.v1.ClusterRoleBindingList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1/ClusterRoleList":
                resources[`rbac.authorization.k8s.io/v1/ClusterRoleList::${namespace}/${name}`] =
                    new k8s.rbac.v1.ClusterRoleList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1/Role":
                resources[`rbac.authorization.k8s.io/v1/Role::${namespace}/${name}`] =
                    new k8s.rbac.v1.Role(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1/RoleBinding":
                resources[`rbac.authorization.k8s.io/v1/RoleBinding::${namespace}/${name}`] =
                    new k8s.rbac.v1.RoleBinding(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1/RoleBindingList":
                resources[`rbac.authorization.k8s.io/v1/RoleBindingList::${namespace}/${name}`] =
                    new k8s.rbac.v1.RoleBindingList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1/RoleList":
                resources[`rbac.authorization.k8s.io/v1/RoleList::${namespace}/${name}`] =
                    new k8s.rbac.v1.RoleList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1alpha1/ClusterRole":
                resources[`rbac.authorization.k8s.io/v1alpha1/ClusterRole::${namespace}/${name}`] =
                    new k8s.rbac.v1alpha1.ClusterRole(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1alpha1/ClusterRoleBinding":
                resources[`rbac.authorization.k8s.io/v1alpha1/ClusterRoleBinding::${namespace}/${name}`] =
                    new k8s.rbac.v1alpha1.ClusterRoleBinding(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1alpha1/ClusterRoleBindingList":
                resources[`rbac.authorization.k8s.io/v1alpha1/ClusterRoleBindingList::${namespace}/${name}`] =
                    new k8s.rbac.v1alpha1.ClusterRoleBindingList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1alpha1/ClusterRoleList":
                resources[`rbac.authorization.k8s.io/v1alpha1/ClusterRoleList::${namespace}/${name}`] =
                    new k8s.rbac.v1alpha1.ClusterRoleList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1alpha1/Role":
                resources[`rbac.authorization.k8s.io/v1alpha1/Role::${namespace}/${name}`] =
                    new k8s.rbac.v1alpha1.Role(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1alpha1/RoleBinding":
                resources[`rbac.authorization.k8s.io/v1alpha1/RoleBinding::${namespace}/${name}`] =
                    new k8s.rbac.v1alpha1.RoleBinding(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1alpha1/RoleBindingList":
                resources[`rbac.authorization.k8s.io/v1alpha1/RoleBindingList::${namespace}/${name}`] =
                    new k8s.rbac.v1alpha1.RoleBindingList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1alpha1/RoleList":
                resources[`rbac.authorization.k8s.io/v1alpha1/RoleList::${namespace}/${name}`] =
                    new k8s.rbac.v1alpha1.RoleList(name, obj, opts);
                break;
            case "rbac/v1alpha1/Subject":
                resources[`rbac/v1alpha1/Subject::${namespace}/${name}`] =
                    new k8s.rbac.v1alpha1.Subject(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1beta1/ClusterRole":
                resources[`rbac.authorization.k8s.io/v1beta1/ClusterRole::${namespace}/${name}`] =
                    new k8s.rbac.v1beta1.ClusterRole(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1beta1/ClusterRoleBinding":
                resources[`rbac.authorization.k8s.io/v1beta1/ClusterRoleBinding::${namespace}/${name}`] =
                    new k8s.rbac.v1beta1.ClusterRoleBinding(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1beta1/ClusterRoleBindingList":
                resources[`rbac.authorization.k8s.io/v1beta1/ClusterRoleBindingList::${namespace}/${name}`] =
                    new k8s.rbac.v1beta1.ClusterRoleBindingList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1beta1/ClusterRoleList":
                resources[`rbac.authorization.k8s.io/v1beta1/ClusterRoleList::${namespace}/${name}`] =
                    new k8s.rbac.v1beta1.ClusterRoleList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1beta1/Role":
                resources[`rbac.authorization.k8s.io/v1beta1/Role::${namespace}/${name}`] =
                    new k8s.rbac.v1beta1.Role(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1beta1/RoleBinding":
                resources[`rbac.authorization.k8s.io/v1beta1/RoleBinding::${namespace}/${name}`] =
                    new k8s.rbac.v1beta1.RoleBinding(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1beta1/RoleBindingList":
                resources[`rbac.authorization.k8s.io/v1beta1/RoleBindingList::${namespace}/${name}`] =
                    new k8s.rbac.v1beta1.RoleBindingList(name, obj, opts);
                break;
            case "rbac.authorization.k8s.io/v1beta1/RoleList":
                resources[`rbac.authorization.k8s.io/v1beta1/RoleList::${namespace}/${name}`] =
                    new k8s.rbac.v1beta1.RoleList(name, obj, opts);
                break;
            case "scheduling.k8s.io/v1alpha1/PriorityClass":
                resources[`scheduling.k8s.io/v1alpha1/PriorityClass::${namespace}/${name}`] =
                    new k8s.scheduling.v1alpha1.PriorityClass(name, obj, opts);
                break;
            case "scheduling.k8s.io/v1alpha1/PriorityClassList":
                resources[`scheduling.k8s.io/v1alpha1/PriorityClassList::${namespace}/${name}`] =
                    new k8s.scheduling.v1alpha1.PriorityClassList(name, obj, opts);
                break;
            case "settings.k8s.io/v1alpha1/PodPreset":
                resources[`settings.k8s.io/v1alpha1/PodPreset::${namespace}/${name}`] =
                    new k8s.settings.v1alpha1.PodPreset(name, obj, opts);
                break;
            case "settings.k8s.io/v1alpha1/PodPresetList":
                resources[`settings.k8s.io/v1alpha1/PodPresetList::${namespace}/${name}`] =
                    new k8s.settings.v1alpha1.PodPresetList(name, obj, opts);
                break;
            case "storage.k8s.io/v1/StorageClass":
                resources[`storage.k8s.io/v1/StorageClass::${namespace}/${name}`] =
                    new k8s.storage.v1.StorageClass(name, obj, opts);
                break;
            case "storage.k8s.io/v1/StorageClassList":
                resources[`storage.k8s.io/v1/StorageClassList::${namespace}/${name}`] =
                    new k8s.storage.v1.StorageClassList(name, obj, opts);
                break;
            case "storage.k8s.io/v1alpha1/VolumeAttachment":
                resources[`storage.k8s.io/v1alpha1/VolumeAttachment::${namespace}/${name}`] =
                    new k8s.storage.v1alpha1.VolumeAttachment(name, obj, opts);
                break;
            case "storage.k8s.io/v1alpha1/VolumeAttachmentList":
                resources[`storage.k8s.io/v1alpha1/VolumeAttachmentList::${namespace}/${name}`] =
                    new k8s.storage.v1alpha1.VolumeAttachmentList(name, obj, opts);
                break;
            case "storage.k8s.io/v1beta1/StorageClass":
                resources[`storage.k8s.io/v1beta1/StorageClass::${namespace}/${name}`] =
                    new k8s.storage.v1beta1.StorageClass(name, obj, opts);
                break;
            case "storage.k8s.io/v1beta1/StorageClassList":
                resources[`storage.k8s.io/v1beta1/StorageClassList::${namespace}/${name}`] =
                    new k8s.storage.v1beta1.StorageClassList(name, obj, opts);
                break;
            default:
            throw new Error(`Unrecognized resource type ${apiVersion}/${kind}`);
        }
    }
    return resources;
}
