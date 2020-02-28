// Copyright 2016-2018, Pulumi Corporation.
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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi-kubernetes/pkg/gen"
)

// This is the URL for the v1.17.0 swagger spec. This is the last version of the spec containing the following
// deprecated resources:
// - extensions/v1beta1/*
// - apps/v1beta1/*
// - apps/v1beta2/*
// Since these resources will continue to be important to users for the foreseeable future, we will merge in
// newer specs on top of this spec so that these resources continue to be available in our SDKs.
const Swagger117Url = "https://raw.githubusercontent.com/kubernetes/kubernetes/v1.17.0/api/openapi-spec/swagger.json"
const Swagger117FileName = "swagger-v1.17.0.json"

func main() {
	if len(os.Args) < 5 {
		log.Fatal("Usage: gen <language> <swagger-file> <template-dir> <out-dir>")
	}

	language := os.Args[1]

	swagger, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	swaggerDir := filepath.Dir(os.Args[2])

	legacySwaggerPath := filepath.Join(swaggerDir, Swagger117FileName)
	err = DownloadFile(legacySwaggerPath, Swagger117Url)
	if err != nil {
		panic(err)
	}
	legacySwagger, err := ioutil.ReadFile(legacySwaggerPath)
	if err != nil {
		panic(err)
	}
	mergedSwagger := mergeSwaggerSpecs(legacySwagger, swagger)
	data := mergedSwagger.(map[string]interface{})

	templateDir := os.Args[3]
	outdir := fmt.Sprintf("%s/%s", os.Args[4], language)

	switch language {
	case "nodejs":
		writeNodeJSClient(data, outdir, templateDir)
	case "python":
		writePythonClient(data, outdir, templateDir)
	case "dotnet":
		writeDotnetClient(data, outdir, templateDir)
	default:
		panic(fmt.Sprintf("Unrecognized language '%s'", language))
	}
}

func writeNodeJSClient(data map[string]interface{}, outdir, templateDir string) {
	inputAPIts, ouputAPIts, indexts, yamlts, packagejson, groupsts, err := gen.NodeJSClient(
		data, templateDir)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(outdir, 0700)
	if err != nil {
		panic(err)
	}

	typesDir := fmt.Sprintf("%s/types", outdir)
	err = os.MkdirAll(typesDir, 0700)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/input.ts", typesDir), []byte(inputAPIts), 0777)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/output.ts", typesDir), []byte(ouputAPIts), 0777)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/yaml/yaml.ts", outdir), []byte(yamlts), 0777)
	if err != nil {
		panic(err)
	}

	for groupName, group := range groupsts {
		groupDir := fmt.Sprintf("%s/%s", outdir, groupName)
		err = os.MkdirAll(groupDir, 0700)
		if err != nil {
			panic(err)
		}

		for versionName, version := range group.Versions {
			versionDir := fmt.Sprintf("%s/%s", groupDir, versionName)
			err = os.MkdirAll(versionDir, 0700)
			if err != nil {
				panic(err)
			}

			for kindName, kind := range version.Kinds {
				err = ioutil.WriteFile(fmt.Sprintf("%s/%s.ts", versionDir, kindName), []byte(kind), 0777)
				if err != nil {
					panic(err)
				}
			}

			err = ioutil.WriteFile(fmt.Sprintf("%s/%s.ts", versionDir, "index"), []byte(version.Index), 0777)
			if err != nil {
				panic(err)
			}
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.ts", groupDir, "index"), []byte(group.Index), 0777)
		if err != nil {
			panic(err)
		}
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/index.ts", outdir), []byte(indexts), 0777)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/package.json", outdir), []byte(packagejson), 0777)
	if err != nil {
		panic(err)
	}

	err = CopyFile(
		filepath.Join(templateDir, "CustomResource.ts"), filepath.Join(outdir, "apiextensions", "CustomResource.ts"))
	if err != nil {
		panic(err)
	}

	err = CopyFile(filepath.Join(templateDir, "README.md"), filepath.Join(outdir, "README.md"))
	if err != nil {
		panic(err)
	}

	err = CopyDir(filepath.Join(templateDir, "helm"), filepath.Join(outdir, "helm"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s/package.json\n", outdir)
	fmt.Println(err)
}

func writePythonClient(data map[string]interface{}, outdir, templateDir string) {
	sdkDir := filepath.Join(outdir, "pulumi_kubernetes")

	err := gen.PythonClient(data, templateDir,
		func(initPy string) error {
			return ioutil.WriteFile(filepath.Join(sdkDir, "__init__.py"), []byte(initPy), 0777)
		},
		func(group, initPy string) error {
			destDir := filepath.Join(sdkDir, group)

			err := os.MkdirAll(destDir, 0700)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(filepath.Join(destDir, "__init__.py"), []byte(initPy), 0777)
		},
		func(crBytes string) error {
			destDir := filepath.Join(sdkDir, "apiextensions")

			err := os.MkdirAll(destDir, 0700)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(
				filepath.Join(destDir, "CustomResource.py"),
				[]byte(crBytes), 0777)
		},
		func(group, version, initPy string) error {
			destDir := filepath.Join(sdkDir, group, version)

			err := os.MkdirAll(destDir, 0700)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(filepath.Join(destDir, "__init__.py"), []byte(initPy), 0777)
		},
		func(group, version, kind, kindPy string) error {
			destDir := filepath.Join(sdkDir, group, version, fmt.Sprintf("%s.py", kind))
			return ioutil.WriteFile(destDir, []byte(kindPy), 0777)
		},
		func(casingPy string) error {
			destDir := filepath.Join(sdkDir, "tables.py")
			return ioutil.WriteFile(destDir, []byte(casingPy), 0777)
		},
		func(yamlPy string) error {
			destDir := filepath.Join(sdkDir, "yaml.py")
			return ioutil.WriteFile(destDir, []byte(yamlPy), 0777)
		})
	if err != nil {
		panic(err)
	}

	err = CopyDir(filepath.Join(templateDir, "helm"), filepath.Join(sdkDir, "helm"))
	if err != nil {
		panic(err)
	}

	err = CopyFile(filepath.Join(templateDir, "README.md"), filepath.Join(sdkDir, "README.md"))
	if err != nil {
		panic(err)
	}
}

func writeDotnetClient(data map[string]interface{}, outdir, templateDir string) {

	inputAPIcs, ouputAPIcs, yamlcs, kindsCs, err := gen.DotnetClient(data, templateDir)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(outdir, 0700)
	if err != nil {
		panic(err)
	}

	typesDir := fmt.Sprintf("%s/Types", outdir)
	err = os.MkdirAll(typesDir, 0700)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/Input.cs", typesDir), []byte(inputAPIcs), 0777)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/Output.cs", typesDir), []byte(ouputAPIcs), 0777)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/Yaml/Yaml.cs", outdir), []byte(yamlcs), 0777)
	if err != nil {
		panic(err)
	}

	for path, contents := range kindsCs {
		filename := filepath.Join(outdir, path)
		err := os.MkdirAll(filepath.Dir(filename), 0700)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(filename, []byte(contents), 0777)
		if err != nil {
			panic(err)
		}
	}

	err = CopyFile(filepath.Join(templateDir, "README.md"), filepath.Join(outdir, "README.md"))
	if err != nil {
		panic(err)
	}

	err = CopyFile(filepath.Join(templateDir, "Utilities.cs"), filepath.Join(outdir, "Utilities.cs"))
	if err != nil {
		panic(err)
	}

	err = CopyFile(filepath.Join(templateDir, "Provider.cs"), filepath.Join(outdir, "Provider.cs"))
	if err != nil {
		panic(err)
	}

	err = CopyFile(filepath.Join(templateDir, "logo.png"), filepath.Join(outdir, "logo.png"))
	if err != nil {
		panic(err)
	}

	err = CopyFile(filepath.Join(templateDir, "Pulumi.Kubernetes.csproj"), filepath.Join(outdir, "Pulumi.Kubernetes.csproj"))
	if err != nil {
		panic(err)
	}
}
