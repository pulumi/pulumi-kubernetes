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
	"path/filepath"

	"github.com/pulumi/pulumi-kubernetes/provider/v2/cmd/crd2pulumi/gen"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	DotNet string = "dotnet"
	Go     string = "go"
	NodeJS string = "nodejs"
	Python string = "python"
)

const (
	DotNetPath string = "dotnetPath"
	GoPath     string = "goPath"
	NodeJSPath string = "nodejsPath"
	PythonPath string = "pythonPath"
)

var defaultOutputPath = "crds/"

const long = `crd2pulumi is a CLI tool that generates typed Kubernetes 
CustomResources to use in Pulumi programs, based on a
CustomResourceDefinition YAML schema.`

const example = `crd2pulumi --nodejs crontabs.yaml
crd2pulumi -dgnp crd-certificates.yaml crd-issuers.yaml crd-challenges.yaml
crd2pulumi --pythonPath=crds/python/istio --nodejsPath=crds/nodejs/istio crd-all.gen.yaml crd-mixer.yaml crd-operator.yaml

Notice that by just setting a language-specific output path (--pythonPath, --nodejsPath, etc) the code will
still get generated, so setting -p, -n, etc becomes unnecessary.
`

func getLanguageSettings(flags *pflag.FlagSet) gen.LanguageSettings {
	nodejs, _ := flags.GetBool(NodeJS)
	python, _ := flags.GetBool(Python)
	dotnet, _ := flags.GetBool(DotNet)
	golang, _ := flags.GetBool(Go)

	nodejsPath, _ := flags.GetString(NodeJSPath)
	pythonPath, _ := flags.GetString(PythonPath)
	dotnetPath, _ := flags.GetString(DotNetPath)
	goPath, _ := flags.GetString(GoPath)

	ls := gen.LanguageSettings{}
	if nodejsPath != "" {
		ls.NodeJSPath = &nodejsPath
	} else if nodejs {
		path := filepath.Join(defaultOutputPath, NodeJS)
		ls.NodeJSPath = &path
	}
	if pythonPath != "" {
		ls.PythonPath = &pythonPath
	} else if python {
		path := filepath.Join(defaultOutputPath, Python)
		ls.PythonPath = &path
	}
	if dotnetPath != "" {
		ls.DotNetPath = &dotnetPath
	} else if dotnet {
		path := filepath.Join(defaultOutputPath, DotNet)
		ls.DotNetPath = &path
	}
	if goPath != "" {
		ls.GoPath = &goPath
	} else if golang {
		path := filepath.Join(defaultOutputPath, Go)
		ls.GoPath = &path
	}
	return ls
}

var (
	rootCmd = &cobra.Command{
		Use:     "crd2pulumi [-dgnp] [--nodejsPath path] [--pythonPath path] [--dotnetPath path] [--goPath path] <crd1.yaml> [crd2.yaml ...]",
		Short:   "A tool that generates typed Kubernetes CustomResources",
		Long:    long,
		Example: example,
		Version: gen.Version,
		Args: func(cmd *cobra.Command, args []string) error {
			emptyLanguageSettings := gen.LanguageSettings{}
			if getLanguageSettings(cmd.Flags()) == emptyLanguageSettings {
				return errors.New("must specify at least one language")
			}

			err := cobra.MinimumNArgs(1)(cmd, args)
			if err != nil {
				return errors.New("must specify at least one CRD YAML file")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool("force")
			languageSettings := getLanguageSettings(cmd.Flags())

			err := gen.Generate(languageSettings, args, force)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(-1)
			}

			fmt.Println("Successfully generated code.")
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

var forceValue bool
var nodeJSValue, pythonValue, dotNetValue, goValue bool
var nodeJSPathValue, pythonPathValue, dotNetPathValue, goPathValue string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&nodeJSValue, NodeJS, "n", false, "generate NodeJS")
	rootCmd.PersistentFlags().BoolVarP(&pythonValue, Python, "p", false, "generate Python")
	rootCmd.PersistentFlags().BoolVarP(&dotNetValue, DotNet, "d", false, "generate .NET")
	rootCmd.PersistentFlags().BoolVarP(&goValue, Go, "g", false, "generate Go")

	rootCmd.PersistentFlags().StringVar(&nodeJSPathValue, NodeJSPath, "", "optional NodeJS output dir")
	rootCmd.PersistentFlags().StringVar(&pythonPathValue, PythonPath, "", "optional Python output dir")
	rootCmd.PersistentFlags().StringVar(&dotNetPathValue, DotNetPath, "", "optional .NET output dir")
	rootCmd.PersistentFlags().StringVar(&goPathValue, GoPath, "", "optional Go output dir")

	rootCmd.PersistentFlags().BoolVarP(&forceValue, "force", "f", false, "overwrite existing files")
}
