// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"github.com/pulumi/lumi/pkg/resource"
)

// A composite key is a key that contains multiple pieces: a resource ID string followed by a series of key/value
// pairs.  This is used when Terraform needs multiple keys in order to lookup the associated resource state, versus
// just a single string ID.  There are even places where Terraform *could* have used a single ID, but chose not to,
// in which case we also need to use a composite.  This is unfortunate but lets us continue to implement Get/Query
// routines on the Lumi side, which depend on singular keys, while still giving Terraform what it needs.

// A few constants for delimiters used when encoding/parsing composite keys (e.g., "id;k1=v1;k2=v2").
const (
	compSep    = ";"
	compKeySep = "="
)

// isCompositeKey returns true if the ID is a composite key.
func isCompositeKey(id resource.ID) bool {
	return strings.Contains(string(id), compSep)
}

// createCompositeKey creates a single string ID out of the keys plus state.
func createCompositeKey(keys []string, state *terraform.InstanceState,
	custom map[string]SchemaInfo) (resource.ID, error) {
	// Make a copy of the keys and always emit them in sorted order.
	var sortedKeys []string
	sortedKeys = append(sortedKeys, keys...)
	sort.Strings(sortedKeys)

	// Now start with the ID and then append all keys from the state.
	id := state.ID
	for _, key := range sortedKeys {
		value, has := state.Attributes[key]
		if !has {
			return "", errors.Errorf("Terraform state value for composite key '%v' is missing", key)
		}
		luminame, _ := getInfoFromTerraformName(key, custom)
		id += fmt.Sprintf("%v%v%v%v", compSep, luminame, compKeySep, value)
	}
	return resource.ID(id), nil
}

// parseCompositeKey takes an ID and, if it is composite, parses out the components and returns and resulting ID part
// plus the composite keys in the resulting map.
func parseCompositeKey(id resource.ID, custom map[string]SchemaInfo) (string, map[string]string, error) {
	sid := string(id)

	// If there wasn't any composite key marker; just return the entire ID as the resulting string.
	cix := strings.Index(sid, compSep)
	if cix == -1 {
		return sid, nil, nil
	}

	// Otherwise, parse out the pieces.
	parsedID := sid[:cix] // extract the ID part.
	sid = sid[cix+1:]     // and skip to the remaining key/value parts.
	parsedKeys := make(map[string]string)
	for len(sid) > 0 {
		cix = strings.Index(sid, compSep)
		var keyval string
		if cix == -1 {
			keyval = sid // it's all a key/value.
			sid = ""     // empty the string so we exit afterwards.
		} else {
			keyval = sid[:cix] // extract this key/value.
			sid = sid[cix+1:]  // and skip over it for subsequent parsing.
		}

		// Look for the = delimiter.
		kvix := strings.Index(keyval, compKeySep)
		if kvix == -1 {
			return "", nil,
				errors.Errorf("Malformed composite key ID string '%v'; missing '=' in %v", id, keyval)
		}
		key := keyval[:kvix]
		tfkey, _ := getInfoFromLumiName(resource.PropertyKey(key), custom) // return the Terraform name.
		val := keyval[kvix+1:]
		parsedKeys[tfkey] = val
	}

	return parsedID, parsedKeys, nil
}
