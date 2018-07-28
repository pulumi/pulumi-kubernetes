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
	if len(os.Args) < 4 {
		log.Fatal("Usage: gen <swagger-file> <template-dir> <out-dir>")
	}

	swagger, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(swagger, &data)
	if err != nil {
		panic(err)
	}

	templateDir := os.Args[2]
	inputAPIts, ouputAPIts, providerts, helmts, packagejson, err := gen.NodeJSClient(data, templateDir)
	if err != nil {
		panic(err)
	}

	outdir := os.Args[3]
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

	err = ioutil.WriteFile(fmt.Sprintf("%s/helm.ts", outdir), []byte(helmts), 0777)
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
