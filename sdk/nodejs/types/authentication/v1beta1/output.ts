// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../../../types/input";
import * as outputs from "../../../types/output";
import * as enums from "../../../types/enums";
import * as utilities from "../../../utilities";

/**
 * TokenReviewSpec is a description of the token authentication request.
 */
export interface TokenReviewSpec {
    /**
     * Audiences is a list of the identifiers that the resource server presented with the token identifies as. Audience-aware token authenticators will verify that the token was intended for at least one of the audiences in this list. If no audiences are provided, the audience will default to the audience of the Kubernetes apiserver.
     */
    audiences: string[];
    /**
     * Token is the opaque bearer token.
     */
    token: string;
}

/**
 * TokenReviewSpec is a description of the token authentication request.
 */
export interface TokenReviewSpecPatch {
    /**
     * Audiences is a list of the identifiers that the resource server presented with the token identifies as. Audience-aware token authenticators will verify that the token was intended for at least one of the audiences in this list. If no audiences are provided, the audience will default to the audience of the Kubernetes apiserver.
     */
    audiences: string[];
    /**
     * Token is the opaque bearer token.
     */
    token: string;
}

/**
 * TokenReviewStatus is the result of the token authentication request.
 */
export interface TokenReviewStatus {
    /**
     * Audiences are audience identifiers chosen by the authenticator that are compatible with both the TokenReview and token. An identifier is any identifier in the intersection of the TokenReviewSpec audiences and the token's audiences. A client of the TokenReview API that sets the spec.audiences field should validate that a compatible audience identifier is returned in the status.audiences field to ensure that the TokenReview server is audience aware. If a TokenReview returns an empty status.audience field where status.authenticated is "true", the token is valid against the audience of the Kubernetes API server.
     */
    audiences: string[];
    /**
     * Authenticated indicates that the token was associated with a known user.
     */
    authenticated: boolean;
    /**
     * Error indicates that the token couldn't be checked
     */
    error: string;
    /**
     * User is the UserInfo associated with the provided token.
     */
    user: outputs.authentication.v1beta1.UserInfo;
}

/**
 * TokenReviewStatus is the result of the token authentication request.
 */
export interface TokenReviewStatusPatch {
    /**
     * Audiences are audience identifiers chosen by the authenticator that are compatible with both the TokenReview and token. An identifier is any identifier in the intersection of the TokenReviewSpec audiences and the token's audiences. A client of the TokenReview API that sets the spec.audiences field should validate that a compatible audience identifier is returned in the status.audiences field to ensure that the TokenReview server is audience aware. If a TokenReview returns an empty status.audience field where status.authenticated is "true", the token is valid against the audience of the Kubernetes API server.
     */
    audiences: string[];
    /**
     * Authenticated indicates that the token was associated with a known user.
     */
    authenticated: boolean;
    /**
     * Error indicates that the token couldn't be checked
     */
    error: string;
    /**
     * User is the UserInfo associated with the provided token.
     */
    user: outputs.authentication.v1beta1.UserInfoPatch;
}

/**
 * UserInfo holds the information about the user needed to implement the user.Info interface.
 */
export interface UserInfo {
    /**
     * Any additional information provided by the authenticator.
     */
    extra: {[key: string]: string[]};
    /**
     * The names of groups this user is a part of.
     */
    groups: string[];
    /**
     * A unique value that identifies this user across time. If this user is deleted and another user by the same name is added, they will have different UIDs.
     */
    uid: string;
    /**
     * The name that uniquely identifies this user among all active users.
     */
    username: string;
}

/**
 * UserInfo holds the information about the user needed to implement the user.Info interface.
 */
export interface UserInfoPatch {
    /**
     * Any additional information provided by the authenticator.
     */
    extra: {[key: string]: string[]};
    /**
     * The names of groups this user is a part of.
     */
    groups: string[];
    /**
     * A unique value that identifies this user across time. If this user is deleted and another user by the same name is added, they will have different UIDs.
     */
    uid: string;
    /**
     * The name that uniquely identifies this user among all active users.
     */
    username: string;
}

