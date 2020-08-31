// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/pulumi/pulumi-kubernetes/provider/v2/cmd/crd2pulumi/gen"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "crd2pulumi <language> <crd file path> [output directory]",
		Short: "A tool that generates strongly-typed k8s CRDS",
		Long: `crd2pulumi is a CLI tool that generates strongly-typed
		Kubernetes CRD resources to use in Pulumi programs.`,
		Example: "crd2pulumi nodejs crontab.yaml crds/\ncrd2pulumi go cert-manager/crd-orders.yaml",
		Version: "1.0.1",
		Args: func(cmd *cobra.Command, args []string) error {
			err := cobra.RangeArgs(2, 3)(cmd, args)
			if err != nil {
				return err
			}

			language := args[0]
			if !gen.IsValidLanguage(language) {
				return errors.New("unsupported language " + language)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			language, yamlPath, outputDir := args[0], args[1], ""
			if len(args) > 2 {
				outputDir = args[2]
			}
			force, _ := cmd.Flags().GetBool("force")

			err := gen.Generate(language, yamlPath, outputDir, force)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(-1)
			}

			fmt.Printf("Successfully generated code for %s\n", yamlPath)
		},
	}
)

var Force bool

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Force, "force", "f", false, "forcefully overwite existing files")
}
