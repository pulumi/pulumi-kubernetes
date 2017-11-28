// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfgen

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
	"github.com/pulumi/pulumi/pkg/util/contract"
	"github.com/spf13/cobra"

	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
)

// Main executes the TFGen process for the given package pkg and provider prov.
func Main(pkg string, prov tfbridge.ProviderInfo) {
	if err := newTFGenCmd(pkg, prov).Execute(); err != nil {
		_, fmterr := fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		contract.IgnoreError(fmterr)
		os.Exit(-1)
	}
}

func newTFGenCmd(pkg string, prov tfbridge.ProviderInfo) *cobra.Command {
	var logToStderr bool
	var outDir string
	var overlaysDir string
	var quiet bool
	var verbose int
	cmd := &cobra.Command{
		Use:   "lumi-tfgen-" + pkg,
		Short: "The Pulumi TFGen compiler generates Pulumi Fabric metadata from a Terraform provider",
		Long: "The Pulumi TFGen compiler generates Pulumi Fabric metadata from a Terraform provider.\n" +
			"\n" +
			"\n" +
			"The tool will load the provider from your $PATH, inspect its contents dynamically,\n" +
			"and generate all of the Pulumi metadata necessary to consume the resources.\n" +
			"\n" +
			"Note that there is no custom Pulumi Fabric provider code required, because the standard\n" +
			"lumi-tfbridge-provider plugin works against all Terraform provider plugins.\n",
		Run: cmdutil.RunFunc(func(cmd *cobra.Command, args []string) error {
			// Let's generate some code!
			g := newGenerator(pkg, prov)
			if err := g.Generate(outDir, overlaysDir); err != nil {
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
