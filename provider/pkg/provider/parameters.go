// Copyright 2016-2023, Pulumi Corporation.
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

package provider

import (
	"errors"

	flag "github.com/spf13/pflag"
)

type CRDParameters struct {
	// Package name comes in via the key property on a provider request object
	PackageVersion string
	YamlPaths      []string
}

func ParseCrdArgs(args []string) (*CRDParameters, error) {
	var crdPackageVersion string
	var yamlPaths []string

	flags := flag.NewFlagSet("crdargs", flag.PanicOnError)
	flags.StringVarP(&crdPackageVersion, "version", "v", "", "The version of the CRD package.")
	err := flags.Parse(args)
	if err != nil {
		panic(err)
	}
	yamlPaths = flags.Args()

	if crdPackageVersion == "" {
		return nil, errors.New("package version must be provided")
	}

	if len(yamlPaths) == 0 {
		return nil, errors.New("no locations of yaml files given")
	}

	return &CRDParameters{
		PackageVersion: crdPackageVersion,
		YamlPaths:      yamlPaths,
	}, nil
}
