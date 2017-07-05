// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package main

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/pulumi/lumi/pkg/util/cmdutil"
)

func main() {
	if err := newTFGenCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(-1)
	}
}

func newTFGenCmd() *cobra.Command {
	var logToStderr bool
	var out string
	var quiet bool
	var verbose int
	cmd := &cobra.Command{
		Use:   "lumi-tfgen tf-provider",
		Short: "The Lumi TFGen compiler generates Lumi metadata from a Terraform provider",
		Long: "The Lumi TFGen compiler generates Lumi metadata from a Terraform provider.\n" +
			"\n" +
			"The tool will load the provider from your $PATH, inspect its contents dynamically,\n" +
			"and generate all of the Lumi metadata necessary to consume the resources.\n" +
			"\n" +
			"Note that there is no Lumi provider code required, because the standard\n" +
			"lumi-tfbridge-provider plugin works against all Terraform provider plugins.\n",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmdutil.InitLogging(logToStderr, verbose, true)
		},
		Run: cmdutil.RunFunc(func(cmd *cobra.Command, args []string) error {
			// Let's generate some code!
			g := newGenerator()
			return g.Generate(terraformProviders(), out)
		}),
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			glog.Flush()
		},
	}

	cmd.PersistentFlags().BoolVar(
		&logToStderr, "logtostderr", false, "Log to stderr instead of to files")
	cmd.PersistentFlags().StringVar(
		&out, "out", "", "Save generated package metadata to this directory")
	cmd.PersistentFlags().BoolVarP(
		&quiet, "quiet", "q", false, "Suppress non-error output progress messages")
	cmd.PersistentFlags().IntVarP(
		&verbose, "verbose", "v", 0, "Enable verbose logging (e.g., v=3); anything >3 is very verbose")

	return cmd
}
