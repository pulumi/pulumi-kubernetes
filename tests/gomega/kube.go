//nolint:govet,golint
package gomega

import (
	. "github.com/onsi/gomega/gstruct"
	gomegatypes "github.com/onsi/gomega/types"
)

func HaveSkipAwaitAnnotation() gomegatypes.GomegaMatcher {
	return MatchProps(IgnoreExtras, Props{
		"metadata": MatchObject(IgnoreExtras, Props{
			"annotations": MatchObject(IgnoreExtras, Props{
				"pulumi.com/skipAwait": MatchValue("true"),
			}),
		}),
	})
}
