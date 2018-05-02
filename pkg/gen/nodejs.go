package gen

import (
	"fmt"

	"github.com/cbroglie/mustache"
	providerVersion "github.com/pulumi/pulumi-kubernetes/pkg/version"
)

// --------------------------------------------------------------------------

// Main interface.

// --------------------------------------------------------------------------

// NodeJSClient will generate a Pulumi Kubernetes provider client SDK for nodejs.
func NodeJSClient(
	swagger map[string]interface{}, templateDir string,
) (string, string, string, error) {
	definitions := swagger["definitions"].(map[string]interface{})

	groupsSlice := createGroups(definitions, api)

	apits, err := mustache.RenderFile(fmt.Sprintf("%s/api.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return "", "", "", err
	}

	groupsSlice = createGroups(definitions, provider)

	providerts, err := mustache.RenderFile(fmt.Sprintf("%s/provider.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return "", "", "", err
	}

	packagejson, err := mustache.RenderFile(fmt.Sprintf("%s/package.json.mustache", templateDir),
		map[string]interface{}{
			"ProviderVersion": providerVersion.Version,
		})
	if err != nil {
		return "", "", "", err
	}

	return apits, providerts, packagejson, nil
}
