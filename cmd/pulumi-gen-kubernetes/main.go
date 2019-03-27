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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pulumi/pulumi-kubernetes/pkg/gen"
)

func main() {
	if len(os.Args) < 5 {
		log.Fatal("Usage: gen <language> <swagger-file> <template-dir> <out-dir>")
	}

	language := os.Args[1]

	swagger, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(swagger, &data)
	if err != nil {
		panic(err)
	}

	templateDir := os.Args[3]
	outdir := fmt.Sprintf("%s/%s", os.Args[4], language)

	switch language {
	case "nodejs":
		writeNodeJSClient(data, outdir, templateDir)
	case "python":
		writePythonClient(data, outdir, templateDir)
	default:
		panic(fmt.Sprintf("Unrecognized language '%s'", language))
	}
}

func writeNodeJSClient(data map[string]interface{}, outdir, templateDir string) {
	inputAPIts, ouputAPIts, providerts, helmts, indexts, packagejson, groupsts, err := gen.NodeJSClient(
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

	err = ioutil.WriteFile(fmt.Sprintf("%s/provider.ts", outdir), []byte(providerts), 0777)
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

	err = ioutil.WriteFile(fmt.Sprintf("%s/helm.ts", outdir), []byte(helmts), 0777)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/index.ts", outdir), []byte(indexts), 0777)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/package.json", outdir), []byte(packagejson), 0777)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s/package.json\n", outdir)
	fmt.Println(err)
}

func writePythonClient(data map[string]interface{}, outdir, templateDir string) {
	err := gen.PythonClient(data, templateDir, func(initPy string) error {
		return ioutil.WriteFile(
			fmt.Sprintf("%s/pulumi_kubernetes/__init__.py", outdir), []byte(initPy), 0777)
	}, func(group, initPy string) error {
		path := fmt.Sprintf("%s/pulumi_kubernetes/%s", outdir, group)

		err := os.MkdirAll(path, 0700)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(fmt.Sprintf("%s/__init__.py", path), []byte(initPy), 0777)
	}, func(group, version, initPy string) error {
		path := fmt.Sprintf("%s/pulumi_kubernetes/%s/%s", outdir, group, version)

		err := os.MkdirAll(path, 0700)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(fmt.Sprintf("%s/__init__.py", path), []byte(initPy), 0777)
	}, func(group, version, kind, kindPy string) error {
		path := fmt.Sprintf("%s/pulumi_kubernetes/%s/%s/%s.py", outdir, group, version, kind)
		return ioutil.WriteFile(path, []byte(kindPy), 0777)
	}, func(casingPy string) error {
		return ioutil.WriteFile(
			fmt.Sprintf("%s/pulumi_kubernetes/tables.py", outdir), []byte(casingPy), 0777)
	}, func(yamlPy string) error {
		return ioutil.WriteFile(
			fmt.Sprintf("%s/pulumi_kubernetes/yaml.py", outdir), []byte(yamlPy), 0777)
	})
	if err != nil {
		panic(err)
	}
}
