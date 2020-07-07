package main

import (
	"fmt"
	"io"
)

const pulumiName = "pulumi"
const kubernetesName = "k8s"

// Maps the OpenAPI types to TypeScript types
var types = map[string]string{
	"integer": "number",
	"number":  "number",
	"string":  "string",
	"boolean": "boolean",
}

// Wraps a type around the Pulumi input type. For example, inputType("number")
// returns "pulumi.Input<number>" if const pulumiName = "pulumi".
func inputType(T string) string {
	return fmt.Sprintf("%s.Input<%s>", pulumiName, T)
}

// Returns the modulePath relative to the variable for the imported Pulumi
// module. For example, if const pulumiName = "pulumi", then
// fromPulumi("CustomResourceOptions") = "pulumi.CustomResourceOptions"
func fromPulumi(modulePath string) string {
	return fmt.Sprintf("%s.%s", pulumiName, modulePath)
}

// Returns the modulePath relative to the variable for the imported
// pulumi-kubernetes module. For example, if const kubernetesName = "k8s", then
// fromKubernetes("apiextensions.CustomResource") =
// "k8s.apiextensions.CustomResource"
func fromKubernetes(modulePath string) string {
	return fmt.Sprintf("%s.%s", kubernetesName, modulePath)
}

func (r ResourceDefinition) GenerateNodeJS(w io.Writer) {
	// Print the import statements
	fmt.Fprintf(w, "import * as %s from \"@pulumi/pulumi\";\n", pulumiName)
	fmt.Fprintf(w, "import * as %s from \"@pulumi/kubernetes\";\n\n", kubernetesName)

	// Print the custom resource definition class
	fmt.Fprintf(w, "export class %s extends %s {\n", r.definitionName(), fromKubernetes("apiextensions.v1.CustomResourceDefinition"))
	fmt.Fprintf(w, "\tconstructor(name: string, opts?: %s) {\n", fromPulumi("CustomResourceOptions"))
	fmt.Fprintf(w, "\t\tsuper(name, %s, opts)\n", r.json)
	fmt.Fprintf(w, "\t}\n}\n\n")

	// Print the custom resource class
	fmt.Fprintf(w, "export class %s extends %s {\n", r.kind, fromKubernetes("apiextensions.CustomResource"))
	fmt.Fprintf(w, "\tconstructor(name: string, args?: %s, opts?: %s) {\n", r.argsName(), fromPulumi("CustomResourceOptions"))
	fmt.Fprintf(w, "\t\tsuper(name, { apiVersion: \"%s\", kind: \"%s\", ...args }, opts)\n", r.apiVersion(), r.kind)
	fmt.Fprintf(w, "\t}\n}\n\n")

	// Print the custom resource args interface
	fmt.Fprintf(w, "export interface %s {\n", r.argsName())
	fmt.Fprintf(w, "\tmetadata?: %s;\n", inputType(fromKubernetes("types.input.meta.v1.ObjectMeta")))
	fmt.Fprintf(w, "\tspec?: %s;\n", r.specArgsName())
	fmt.Fprintf(w, "}\n\n")

	// Print the custom resource spec args interface
	fmt.Fprintf(w, "export interface %s {\n", r.specArgsName())
	for _, field := range r.fields {
		fmt.Fprintf(w, "\t%s", field.name)
		if !field.required {
			fmt.Fprintf(w, "?")
		}
		fmt.Fprintf(w, ": %s;\n", inputType(types[field.openType]))
	}
	fmt.Fprintf(w, "}\n")
}
