// Copyright 2016-2019, Pulumi Corporation.
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

package cluster

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
)

// ServerVersion captures k8s major.minor version in a parsed form
type ServerVersion struct {
	Major, Minor int
}

func (v ServerVersion) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// Compare returns -1/0/+1 iff v is less than / equal / greater than input version.
func (v ServerVersion) Compare(version ServerVersion) int {
	a := v.Major
	b := version.Major

	if a == b {
		a = v.Minor
		b = version.Minor
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

// TryGetServerVersion attempts to retrieve the server version from k8s.
// Returns the configured default version in case this fails.
func TryGetServerVersion(cdi discovery.CachedDiscoveryInterface) ServerVersion {
	defaultSV := ServerVersion{
		Major: 1,
		Minor: 14,
	}

	if sv, err := cdi.ServerVersion(); err == nil {
		if v, err := parseVersion(sv); err == nil {
			return v
		}

		return defaultSV
	}

	return defaultSV
}

// gitVersion captures k8s major.minor.patch version in a parsed form
type gitVersion struct {
	Major, Minor, Patch int
}

func (gv gitVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", gv.Major, gv.Minor, gv.Patch)
}

func parseGitVersion(versionString string) (gitVersion, error) {
	// Format v0.0.0(-master+$Format:%h$)
	gitVersionRe := regexp.MustCompile(`v([0-9]+).([0-9]+).([0-9]+).*`)

	parsedVersion := gitVersionRe.FindStringSubmatch(versionString)
	if len(parsedVersion) != 4 {
		err := fmt.Errorf("unable to parse git version %q", versionString)
		return gitVersion{}, err
	}

	var gv gitVersion
	var err error
	gv.Major, err = strconv.Atoi(parsedVersion[1])
	if err != nil {
		return gitVersion{}, err
	}
	gv.Minor, err = strconv.Atoi(parsedVersion[2])
	if err != nil {
		return gitVersion{}, err
	}
	gv.Patch, err = strconv.Atoi(parsedVersion[3])
	if err != nil {
		return gitVersion{}, err
	}

	return gv, nil
}

// parseVersion parses version.Info into a serverVersion struct
func parseVersion(v *version.Info) (ServerVersion, error) {
	fallbackToGitVersion := false

	major, err := strconv.Atoi(v.Major)
	if err != nil {
		fallbackToGitVersion = true
	}

	// trim "+" in minor version (happened on GKE)
	v.Minor = strings.TrimSuffix(v.Minor, "+")

	minor, err := strconv.Atoi(v.Minor)
	if err != nil {
		fallbackToGitVersion = true
	}

	if fallbackToGitVersion {
		gv, err := parseGitVersion(v.GitVersion)
		if err != nil {
			return ServerVersion{}, err
		}

		return ServerVersion{Major: gv.Major, Minor: gv.Minor}, nil
	}

	return ServerVersion{Major: major, Minor: minor}, nil
}
