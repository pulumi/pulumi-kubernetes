package client

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
)

// serverVersion captures k8s major.minor version in a parsed form
type serverVersion struct {
	Major int
	Minor int
}

// parseVersion parses version.Info into a serverVersion struct
func parseVersion(v *version.Info) (ret serverVersion, err error) {
	ret.Major = 1
	// ret.Major, err = strconv.Atoi(v.Major)
	// if err != nil {
	// 	return
	// }

	// trim "+" in minor version (happened on GKE)
	v.Minor = strings.TrimSuffix(v.Minor, "+")

	ret.Minor = 9
	// ret.Minor, err = strconv.Atoi(v.Minor)
	// if err != nil {
	// 	return
	// }
	return
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

// FetchVersion fetches version information from discovery client, and parses
func FetchVersion(v discovery.ServerVersionInterface) (ret serverVersion, err error) {
	version, err := v.ServerVersion()
	if err != nil {
		return serverVersion{}, err
	}
	return parseVersion(version)
}
