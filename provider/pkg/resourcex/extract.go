// Copyright 2016-2024, Pulumi Corporation.
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

package resourcex

import (
	"fmt"
	"slices"

	"github.com/mitchellh/mapstructure"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type ExtractOptions struct {
	// TagName is the struct tag name to use for resource property names; defaults to "json".
	TagName string
	// RejectUnknowns produces an error (of type ContainsUnknownsError) if any unknowns are extracted.
	RejectUnknowns bool
}

type ExtractResult struct {
	// Dependencies is the set of resources that the extracted information depends on (known or unknown).
	Dependencies []resource.URN
	// ContainsUnknowns is true if the extracted information contains unknown values.
	ContainsUnknowns bool
	// ContainsSecrets is true if the extracted information contains secret values.
	ContainsSecrets bool
}

// Extract extracts information from a property map into a target struct.
// It returns a summary of the outputness and secretness of the extracted information.
func Extract(target interface{}, props resource.PropertyMap, opts ExtractOptions) (ExtractResult, error) {
	if opts.TagName == "" {
		opts.TagName = "json"
	}

	// decode the property map into a JSON-like structure containing only values.
	stripped := DecodeValues(props)

	// deserialize the JSON-like structure into a strongly typed struct.
	config := &mapstructure.DecoderConfig{
		Metadata:   &mapstructure.Metadata{},
		Result:     target,
		TagName:    opts.TagName,
		ZeroFields: true, // for arrays where an element itself is unknown, produces a nil element and requisite metadata.
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return ExtractResult{}, err
	}
	err = decoder.Decode(stripped)
	if err != nil {
		return ExtractResult{}, err
	}

	// use the decoder metadata to visit the properties that were decoded.
	// Note that unknown values are decoded as nils, and are visited.
	// Used properties ("Metadata.Key") are those that appear on the target struct and have a value.
	r := ExtractResult{}
	err = r.visit(resource.NewObjectProperty(props), config.Metadata.Keys)
	if err != nil {
		return ExtractResult{}, err
	}

	if opts.RejectUnknowns && r.ContainsUnknowns {
		return ExtractResult{}, NewContainsUnknownsError(r.Dependencies)
	}
	return r, nil
}

// visit extracts summary information about the given property paths.
func (r *ExtractResult) visit(props resource.PropertyValue, paths []string) error {
	for _, path := range paths {
		p, err := resource.ParsePropertyPathStrict(path)
		if err != nil {
			return fmt.Errorf("parse error for path %q: %v", path, err)
		}
		visitor := func(v resource.PropertyValue) {
			switch {
			case v.IsComputed():
				r.ContainsUnknowns = true
			case v.IsOutput():
				r.Dependencies = mergeDependencies(r.Dependencies, v.OutputValue().Dependencies...)
				if !v.OutputValue().Known {
					r.ContainsUnknowns = true
				}
				r.ContainsSecrets = r.ContainsSecrets || v.OutputValue().Secret
			case v.IsSecret():
				r.ContainsSecrets = true
			}
		}
		Traverse(props, p, visitor)
	}
	return nil
}

func mergeDependencies(slice []resource.URN, elems ...resource.URN) []resource.URN {
	for _, r := range elems {
		if !slices.Contains(slice, r) {
			slice = append(slice, r)
		}
	}
	return slice
}
