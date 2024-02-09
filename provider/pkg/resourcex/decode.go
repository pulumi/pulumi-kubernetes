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
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

// DecodeValues decodes a property map into a JSON-like structure containing only values.
// Unknown values are decoded as nil, both in maps and arrays.
// Secrets are collapsed into their underlying values.
func DecodeValues(props resource.PropertyMap) interface{} {
	return decodeM(props)
}

// decodeV returns a mapper-compatible object map, suitable for deserialization into structures.
func decodeM(props resource.PropertyMap) map[string]interface{} {
	obj := make(map[string]interface{})
	for _, k := range props.StableKeys() {
		key := string(k)
		obj[key] = decodeV(props[k])
	}
	return obj
}

// decodeV returns a mapper-compatible object map, suitable for deserialization into structures.
func decodeV(v resource.PropertyValue) interface{} {
	if v.IsNull() {
		return nil
	} else if v.IsBool() {
		return v.BoolValue()
	} else if v.IsNumber() {
		return v.NumberValue()
	} else if v.IsString() {
		return v.StringValue()
	} else if v.IsArray() {
		arr := make([]interface{}, len(v.ArrayValue()))
		for i := 0; i < len(v.ArrayValue()); i++ {
			arr[i] = decodeV(v.ArrayValue()[i])
		}
		return arr
	} else if v.IsAsset() {
		contract.Failf("unsupported value type '%v'", v.TypeString())
		return nil
	} else if v.IsArchive() {
		contract.Failf("unsupported value type '%v'", v.TypeString())
		return nil
	} else if v.IsComputed() {
		return nil // zero value for unknowns
	} else if v.IsOutput() {
		if !v.OutputValue().Known {
			return nil // zero value for unknowns
		}
		return decodeV(v.OutputValue().Element)
	} else if v.IsSecret() {
		return decodeV(v.SecretValue().Element)
	} else if v.IsResourceReference() {
		contract.Failf("unsupported value type '%v'", v.TypeString())
		return nil
	} else if v.IsObject() {
		return decodeM(v.ObjectValue())
	} else {
		contract.Failf("unexpected value type '%v'", v.TypeString())
		return nil
	}
}
