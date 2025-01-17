// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../../../types/input";
import * as outputs from "../../../types/output";
import * as enums from "../../../types/enums";
import * as utilities from "../../../utilities";

/**
 * Specification defining the Helm chart repository to use.
 */
export interface RepositoryOpts {
    /**
     * The Repository's CA File
     */
    caFile?: pulumi.Input<string>;
    /**
     * The repository's cert file
     */
    certFile?: pulumi.Input<string>;
    /**
     * The repository's cert key file
     */
    keyFile?: pulumi.Input<string>;
    /**
     * Password for HTTP basic authentication
     */
    password?: pulumi.Input<string>;
    /**
     * Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
     */
    repo?: pulumi.Input<string>;
    /**
     * Username for HTTP basic authentication
     */
    username?: pulumi.Input<string>;
}
