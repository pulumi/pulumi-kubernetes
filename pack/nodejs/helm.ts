import { execSync } from "child_process";
import * as fs from "fs";

import * as k8s from "./index";
import * as pulumi from "@pulumi/pulumi";
import * as shell from "shell-quote";
import * as tmp from "tmp";
import * as path from "./path";
import * as nodepath from "path";

export namespace v2 {
    export interface ChartOpts {
        repo: string;
        chart: string;
        version: string;

        namespace?: string;
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
    export class Chart extends k8s.yaml.CollectionComponentResource {
        constructor(releaseName: string, config: ChartOpts, opts?: pulumi.ComponentResourceOptions) {
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
                const chart = path.quotePath(nodepath.join(chartDir.name, config.chart));
                const release = shell.quote([releaseName]);
                const values = path.quotePath(overrides.name);
                const namespaceArg = config.namespace ? `--namespace ${shell.quote([config.namespace])}` : "";
                const yamlStream = execSync(
                    `helm template ${chart} --name ${release} --values ${values} ${namespaceArg}`
                ).toString();
                this.resources = k8s.yaml.parse({
                    yaml: [yamlStream],
                    transformations: config.transformations || [],
                }, { parent: this });
            } catch (e) {
                // Shed stack trace, only emit the error.
                throw new pulumi.RunError(e.toString());
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

        // For arguments that are not paths to files, it is sufficent to use shell.quote to quote the arguments.
        // However, for arguments that are actual paths to files we use path.quotePath (note that path here is
        // not the node path builtin module). This ensures proper escaping of paths on Windows.
        if (opts.version !== undefined)     { flags.push(`--version ${shell.quote([opts.version])}`);         }
        if (opts.caFile !== undefined)      { flags.push(`--ca-file ${path.quotePath(opts.caFile)}`);          }
        if (opts.certFile !== undefined)    { flags.push(`--cert-file ${path.quotePath(opts.certFile)}`);      }
        if (opts.keyFile !== undefined)     { flags.push(`--key-file ${path.quotePath(opts.keyFile)}`);        }
        if (opts.destination !== undefined) { flags.push(`--destination ${path.quotePath(opts.destination)}`); }
        if (opts.keyring !== undefined)     { flags.push(`--keyring ${path.quotePath(opts.keyring)}`);         }
        if (opts.password !== undefined)    { flags.push(`--password ${shell.quote([opts.password])}`);       }
        if (opts.repo !== undefined)        { flags.push(`--repo ${shell.quote([opts.repo])}`);               }
        if (opts.untardir !== undefined)    { flags.push(`--untardir ${path.quotePath(opts.untardir)}`);       }
        if (opts.username !== undefined)    { flags.push(`--username ${shell.quote([opts.username])}`);       }
        if (opts.home !== undefined)        { flags.push(`--home ${path.quotePath(opts.home)}`);               }
        if (opts.devel === true)            { flags.push(`--devel`);                                          }
        if (opts.prov === true)             { flags.push(`--prov`);                                           }
        if (opts.verify === true)           { flags.push(`--verify`);                                         }
    }
    execSync(`helm fetch ${shell.quote([chart])} ${flags.join(" ")}`);
}
