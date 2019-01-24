// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package await

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi/pkg/diag"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
	"k8s.io/apimachinery/pkg/version"
)

// Format v0.0.0(-master+$Format:%h$)
var gitVersionRe = regexp.MustCompile(`v([0-9]+).([0-9]+).([0-9]+).*`)

// serverVersion captures k8s major.minor version in a parsed form
type serverVersion struct {
	Major, Minor int
}

// Compare returns -1/0/+1 iff v is less than / equal / greater than major.minor
func (v serverVersion) Compare(major, minor int) int {
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

func (v serverVersion) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// gitVersion captures k8s major.minor.patch version in a parsed form
type gitVersion struct {
	Major, Minor, Patch int
}

func (gv gitVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", gv.Major, gv.Minor, gv.Patch)
}

// DefaultVersion takes a wild guess (v1.9) at the version of a Kubernetes cluster.
func defaultVersion() serverVersion {
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
	return serverVersion{Major: 1, Minor: 9}
}

func parseGitVersion(versionString string) (gitVersion, error) {
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
func parseVersion(v *version.Info) (serverVersion, error) {
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
			return serverVersion{}, err
		}

		return serverVersion{Major: gv.Major, Minor: gv.Minor}, nil
	}

	return serverVersion{Major: major, Minor: minor}, nil
}

// canonicalizeDeploymentAPIVersion unifies the various pre-release apiVerion values for a
// Deployment into "apps/v1".
func canonicalizeDeploymentAPIVersion(ver string) string {
	switch ver {
	case "extensions/v1beta1", "apps/v1beta1", "apps/v1beta2", "apps/v1":
		// Canonicalize all of these to "apps/v1".
		return "apps/v1"
	default:
		// If the input version was not a version we understand, just return it as-is.
		return ver
	}
}

// canonicalizeStatefulSetAPIVersion unifies the various pre-release apiVersion values for a
// StatefulSet into "apps/v1".
func canonicalizeStatefulSetAPIVersion(ver string) string {
	switch ver {
	case "apps/v1beta1", "apps/v1beta2", "apps/v1":
		// Canonicalize all of these to "apps/v1".
		return "apps/v1"
	default:
		// If the input version was not a version we understand, just return it as-is.
		return ver
	}
}
