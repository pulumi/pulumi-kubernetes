import * as azure from "@pulumi/azure";
import * as pulumi from "@pulumi/pulumi";
import * as acsengine from "./acsengine";

export interface ACSArgs {
    parameters: pulumi.Input<acsengine.ACSParameters>;
    resourceGroupName: pulumi.Input<string>;
    location: pulumi.Input<string>;
}

export class ACS extends pulumi.ComponentResource {
    public masterFQDN: pulumi.Output<string>;
    public kubeconfig: pulumi.Output<any>;
    constructor(name: string, args: ACSArgs, opts?: pulumi.ResourceOptions) {
        super("azurecontainerservice:acs:acs", name, args, opts);

        const engine = new acsengine.ACSEngine(name, {
            parameters: args.parameters,
            location: args.location,
        }, { parent: this });
        
        // Deploy the template+parameters
        const acs = new azure.core.TemplateDeployment(name, {
            resourceGroupName: args.resourceGroupName,
            deploymentMode: "Incremental",
            templateBody: engine.armTemplate.apply(JSON.stringify),
            parametersBody: engine.armParameters.apply(JSON.stringify),
        }, { parent: this });

        // Set outputs
        this.masterFQDN = acs.outputs.apply(outputs => outputs["masterFQDN"]);
        this.kubeconfig = engine.kubeconfig;
        this.registerOutputs({
            templateId: acs.id,
            templateOutputs: acs.outputs,
        });
    
    }
}
