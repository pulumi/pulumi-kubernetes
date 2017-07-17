// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package main

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/pulumi/lumi/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/pulumi/terraform-bridge/pkg/tfbridge/providers"
)

func main() {
	if err := newTFGenCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(-1)
	}
}

func newTFGenCmd() *cobra.Command {
	var logToStderr bool
	var outDir string
	var overlaysDir string
	var quiet bool
	var verbose int
	cmd := &cobra.Command{
		Use:   "lumi-tfgen tf-provider [providers]",
		Short: "The Lumi TFGen compiler generates Lumi metadata from a Terraform provider",
		Long: "The Lumi TFGen compiler generates Lumi metadata from a Terraform provider.\n" +
			"\n" +
			"By default, TFGen will load all known Terraform providers.  If you wish to only\n" +
			"generate a subset, specify them on the command line as [parameters].\n" +
			"\n" +
			"The tool will load the provider from your $PATH, inspect its contents dynamically,\n" +
			"and generate all of the Lumi metadata necessary to consume the resources.\n" +
			"\n" +
			"Note that there is no Lumi provider code required, because the standard\n" +
			"lumi-tfbridge-provider plugin works against all Terraform provider plugins.\n",
		Run: cmdutil.RunFunc(func(cmd *cobra.Command, args []string) error {
			// Let's generate some code!
			g := newGenerator()
			if len(args) == 0 {
				return errors.Errorf("Missing required 'tf-provider' argument (the module whose package to generate)")
			}
			pkg := args[0]
			if err := g.Generate(pkg, providers.Providers[pkg], outDir, overlaysDir); err != nil {
				return err
			}
			// If we succeeded at generate, but there were errors, exit unsuccessfully.
			if !cmdutil.Diag().Success() {
				os.Exit(-2)
			}
			return nil
		}),
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			glog.Flush()
		},
	}

	cmd.PersistentFlags().BoolVar(
		&logToStderr, "logtostderr", false, "Log to stderr instead of to files")
	cmd.PersistentFlags().StringVarP(
		&outDir, "out", "o", "", "Save generated package metadata to this directory")
	cmd.PersistentFlags().StringVar(
		&overlaysDir, "overlays", "", "Use the target directory for overlays rather than the default of overlays/")
	cmd.PersistentFlags().BoolVarP(
		&quiet, "quiet", "q", false, "Suppress non-error output progress messages")
	cmd.PersistentFlags().IntVarP(
		&verbose, "verbose", "v", 0, "Enable verbose logging (e.g., v=3); anything >3 is very verbose")

	return cmd
}
