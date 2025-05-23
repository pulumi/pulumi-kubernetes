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

package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

// UserAgent is how the provider identifies itself to the API server.
var UserAgent string

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("unable to read build info")
	}
	clientGoVersion := "unknown"
	for _, dep := range bi.Deps {
		if dep.Path != "k8s.io/client-go" {
			continue
		}
		clientGoVersion = dep.Version
	}
	version := "dev"
	if Version != "" {
		version = Version
	}
	UserAgent = fmt.Sprintf("%s/%s (%s/%s) client-go/%s",
		"pulumi-kubernetes",
		version,
		runtime.GOOS,
		runtime.GOARCH,
		clientGoVersion,
	)
}
