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

package client

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi/pkg/diag"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
)

// Format v0.0.0(-master+$Format:%h$)
var gitVersionRe = regexp.MustCompile("v([0-9])+.([0-9])+.[0-9]+.*")

// ServerVersion captures k8s major.minor version in a parsed form
type ServerVersion struct {
	Major int
	Minor int
}

// DefaultVersion takes a wild guess (v1.9) at the version of a Kubernetes cluster.
func DefaultVersion() ServerVersion {
	cmdutil.Diag().Warningf(
		diag.Message("", "Cluster failed to report its version number; falling back to 1.9"), false)

	//
	// Fallback behavior to work around [1]. Some versions of minikube erroneously report a blank
	// `version.Info`, which will cause us to break. It is necessary for us to check this version for
	// `Delete`, because of bugs and quirks in various Kubernetes versions. Currently it is only
	// important that we know the version is above or below 1.5, so here we (hopefully) temporarily
	// choose to fall back to 1.9, which is what most people running minikube use out of the box.
	//
	// [1]: https://github.com/kubernetes/minikube/issues/2505
	//
	return ServerVersion{Major: 1, Minor: 9}
}

func parseGitVersion(gitVersion string) (ServerVersion, error) {
	parsedVersion := gitVersionRe.FindStringSubmatch(gitVersion)
	if len(parsedVersion) != 3 {
		return ServerVersion{}, fmt.Errorf("Unable to parse git version %s", gitVersion)
	}
	var ret ServerVersion
	var err error
	ret.Major, err = strconv.Atoi(parsedVersion[1])
	if err != nil {
		return ServerVersion{}, err
	}
	ret.Minor, err = strconv.Atoi(parsedVersion[2])
	if err != nil {
		return ServerVersion{}, err
	}
	return ret, nil
}

// parseVersion parses version.Info into a serverVersion struct
func parseVersion(v *version.Info) (ret ServerVersion, err error) {
	ret.Major, err = strconv.Atoi(v.Major)
	if err != nil {
		return parseGitVersion(v.GitVersion)
	}

	// trim "+" in minor version (happened on GKE)
	v.Minor = strings.TrimSuffix(v.Minor, "+")

	ret.Minor, err = strconv.Atoi(v.Minor)
	if err != nil {
		return parseGitVersion(v.GitVersion)
	}

	return ret, nil
}

// Compare returns -1/0/+1 iff v is less than / equal / greater than major.minor
func (v ServerVersion) Compare(major, minor int) int {
	a := v.Major
	b := major

	if a == b {
		a = v.Minor
		b = minor
	}

	var res int
	if a > b {
		res = 1
	} else if a == b {
		res = 0
	} else {
		res = -1
	}
	return res
}

func (v ServerVersion) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// FetchVersion fetches version information from discovery client, and parses
func FetchVersion(v discovery.ServerVersionInterface) (ret ServerVersion, err error) {
	version, err := v.ServerVersion()
	if err != nil {
		return ServerVersion{}, err
	}
	return parseVersion(version)
}
