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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/pulumi/pulumi-kubernetes/provider/cmd/crd2pulumi/nodejs"
	"github.com/pulumi/pulumi/sdk/go/common/util/contract"
	unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Language is the SDK language.
type Language string

const (
	DotNet Language = "dotnet"
	Go     Language = "go"
	NodeJS Language = "nodejs"
	Python Language = "python"
)

func main() {
	flag.Usage = func() {
		const usageFormat = "Usage: %s <language> <resourcedefinition.yaml> [output path]"
		_, err := fmt.Fprintf(flag.CommandLine.Output(), usageFormat, os.Args[0])
		contract.IgnoreError(err)
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		return
	}

	language, yamlPath := Language(args[0]), args[1]

	// Read the YAML file, unmarshal it into an unstruct.Unstructured, and parse
	// it into a map[string]pschema.ObjectTypeSpec
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading the yaml file: %v\n", err)
		os.Exit(-1)
	}
	crd, err := UnmarshalYaml(yamlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling yaml: %v\n", err)
		os.Exit(-1)
	}
	types := GetTypes(crd)

	// User can either specify their own output path, or use the default one
	var outputPath string
	if len(args) > 2 {
		outputPath = args[2]
	} else {
		plural, _, err := unstruct.NestedString(crd.Object, "spec", "names", "plural")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting plural name: %v\n", err)
			os.Exit(-1)
		}
		outputPath = getDefaultOutputPath(yamlPath, plural, language)
	}

	switch language {
	case NodeJS:
		code, err := nodejs.GenerateTypeScriptTypes(types)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", errors.Wrap(err, "generating nodejs types"))
			os.Exit(-1)
		}
		err = ioutil.WriteFile(outputPath, code, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", errors.Wrap(err, "outputting to file"))
			os.Exit(-1)
		}
	default:
		panic(fmt.Sprintf("Unrecognized language '%s'", language))
	}
}

func getDefaultOutputPath(yamlPath string, plural string, language Language) string {
	var extension string
	switch language {
	case NodeJS:
		extension = "ts"
	case DotNet:
		extension = "cs"
	case Python:
		extension = "py"
	case Go:
		extension = "go"
	}
	outputFileName := fmt.Sprintf("%s.%s", plural, extension)
	return path.Join(filepath.Dir(yamlPath), outputFileName)
}
