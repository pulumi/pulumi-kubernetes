import * as pulumi from "@pulumi/pulumi";

// ACSParameters are parameters for an ACS Engine deployment
export interface ACSParameters {
    orchestratorProfile: {
        orchestratorType: string;
    };
    masterProfile?: {
        count: number;
        dnsPrefix: string;
        vmSize: string;
    };
    agentPoolProfiles?: Array<{
        name: string;
        count: number;
        vmSize: string;
        availabilityProfile: string;
    }>;
    linuxProfile?: {
        adminUsername: string;
        ssh: {
            publicKeys: Array<{
                keyData: string;
            }>;
        };
    };
    servicePrincipalProfile?: {
        clientId: string;
        secret: string;
    };
}

// ACSEngineArgs are the arguments to create an ACSEngine resource
export interface ACSEngineArgs {
    parameters: pulumi.Input<ACSParameters>;
    location: pulumi.Input<string>;
}

interface ACSInputs {
    parameters: ACSParameters;
    location: string;
}

interface ACSOutputs {
    armTemplate: any;
    armParameters: any;
    kubeconfig: any;
    apiModel: {
        apiVersion: string;
        properties: ACSParameters;
    };
}

// ACSEngine is wrapper around the `acs-engine` utility for generating ARM templates for Container Services including
// Kubernetes and Mesosphere.
export class ACSEngine extends pulumi.dynamic.Resource {
    public armTemplate: pulumi.Output<any>;
    public armParameters: pulumi.Output<any>;
    public kubeconfig: pulumi.Output<any>;
    constructor(name: string, args: ACSEngineArgs, opts?: pulumi.ResourceOptions) {
        const provider: pulumi.dynamic.ResourceProvider = {
            async check(olds: ACSOutputs, news: ACSInputs): Promise<pulumi.dynamic.CheckResult> {
                return {
                    inputs: news,
                };
            },
            async diff(id: pulumi.ID, olds: ACSOutputs, news: ACSInputs): Promise<pulumi.dynamic.DiffResult> {
                return {};
            },
            async create(inputs: ACSInputs): Promise<pulumi.dynamic.CreateResult> {
                const outs = await acsEngineGenerate(inputs.parameters, inputs.location);
                const id = await sha1hash(JSON.stringify(outs));
                return {
                    id: id,
                    outs: outs,
                };
            },
            async update(id: pulumi.ID, olds: ACSOutputs, news: ACSInputs): Promise<pulumi.dynamic.UpdateResult> {
                const props = olds.apiModel.properties;
                deepmerge(props, news.parameters);
                const outs = await acsEngineGenerate(props, news.location);
                return {
                    outs: outs,
                }
            },
            async delete(id: pulumi.ID, props: ACSOutputs): Promise<void> {
                // No cleanup needed.
            },
        }
        super(provider, name, {...args, armTemplate: undefined, armParameters: undefined, kubeconfig: undefined}, opts);
    }
}

// Compute a sha1hash of the given string
async function sha1hash(str: string): Promise<string> {
    const crypto = await import("crypto")
    const shaSum = crypto.createHash("sha1");
    shaSum.update(str);
    return shaSum.digest("hex");
} 

// Recursively merge properties of source into the target object
function deepmerge(target: any, source: any) {
    for (let key in source) {
        if (typeof(source[key]) == "object") {
            deepmerge(target[key], source[key]);
        } else {
            target[key] = source[key];
        }
    }
}

// Generate 
async function acsEngineGenerate(props: ACSParameters, location: string): Promise<ACSOutputs> {
    const tmp = await import("tmp");
    const path = await import("path");
    const fs = await import("fs");
    const { exec } = await import("child_process");

    // Create a directory for model inputs/outputs
    const dir = tmp.dirSync().name;
    console.log(dir);
    // Create the model inputs
    const modelPath = path.join(dir, "model.json");
    const model = {
        apiVersion: "vlabs",
        properties: props,
    };
    console.log(JSON.stringify(model, undefined, " "));
    fs.writeFileSync(modelPath, JSON.stringify(model));
    const outputPath = await new Promise<string>((resolve, reject) => {
        // Generate the ARM templates using `acs-engine`
        const outputPath = path.join(dir, "output");
        exec(`acs-engine generate ${modelPath} --output-directory ${outputPath}`, {
            cwd: dir,
        }, (err, stdout, stderr) => {
            if (err) {
                console.log(stdout);
                console.error(stderr);
                reject(err);
            }
            resolve(outputPath);
        });
    });
    const templateBody = fs.readFileSync(path.join(outputPath, "azuredeploy.json")).toString();
    const parametersString = fs.readFileSync(path.join(outputPath, "azuredeploy.parameters.json")).toString();
    const apiModelBody = fs.readFileSync(path.join(outputPath, "apimodel.json")).toString();
    const parametersObj = JSON.parse(parametersString);
    const parametersBody = JSON.stringify(parametersObj["parameters"]);
    const kubeconfig = fs.readFileSync(path.join(outputPath, "kubeconfig", `kubeconfig.${location}.json`)).toString();
    return {
        armTemplate: JSON.parse(templateBody),
        armParameters: JSON.parse(parametersBody),
        apiModel: JSON.parse(apiModelBody),
        kubeconfig: JSON.parse(kubeconfig),
    };
}