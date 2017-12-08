// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/stretchr/testify/assert"
)

// TestTerraformInputs verifies that we translate Pulumi inputs into Terraform inputs.
func TestTerraformInputs(t *testing.T) {
	result, err := MakeTerraformInputs(
		nil, /*res*/
		nil, /*olds*/
		resource.NewPropertyMapFromMap(map[string]interface{}{
			"nilPropertyValue":    nil,
			"boolPropertyValue":   false,
			"numberPropertyValue": 42,
			"stringo":             "ognirts",
			"arrayPropertyValue":  []interface{}{"an array"},
			"objectPropertyValue": map[string]interface{}{
				"propertyA": "a",
				"propertyB": true,
			},
			"mapPropertyValue": map[string]interface{}{
				"propertyA": "a",
				"propertyB": true,
				"propertyC": map[string]interface{}{
					"nestedPropertyA": true,
				},
			},
			"nestedResources": []map[string]interface{}{{
				"configuration": map[string]interface{}{
					"configurationValue": true,
				},
			}},
		}),
		map[string]*schema.Schema{
			// Type mapPropertyValue as a map so that keys aren't mangled in the usual way.
			"map_property_value": {Type: schema.TypeMap},
			"nested_resources": {
				Type:     schema.TypeList,
				MaxItems: 1,
				// Embed a `*schema.Resource` to validate that type directed
				// walk of the schema successfully walks inside Resources as well
				// as Schemas.
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration": {Type: schema.TypeMap},
					},
				},
			},
		},
		map[string]*SchemaInfo{
			// Reverse map string_property_value to the stringo property.
			"string_property_value": {
				Name: "stringo",
			},
		},
		nil,   /* assets */
		false, /*defaults*/
		false, /*useRawNames*/
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"nil_property_value":    nil,
		"bool_property_value":   false,
		"number_property_value": 42,
		"string_property_value": "ognirts",
		"array_property_value":  []interface{}{"an array"},
		"object_property_value": map[string]interface{}{
			"property_a": "a",
			"property_b": true,
		},
		"map_property_value": map[string]interface{}{
			"propertyA": "a",
			"propertyB": true,
			"propertyC": map[string]interface{}{
				"nestedPropertyA": true,
			},
		},
		"nested_resources": []interface{}{
			map[string]interface{}{
				"configuration": map[string]interface{}{
					"configurationValue": true,
				},
			},
		},
	}, result)
}

// TestTerraformOutputs verifies that we translate Terraform outputs into Pulumi outputs.
func TestTerraformOutputs(t *testing.T) {
	result := MakeTerraformOutputs(
		map[string]interface{}{
			"nil_property_value":    nil,
			"bool_property_value":   false,
			"number_property_value": 42,
			"string_property_value": "ognirts",
			"array_property_value":  []interface{}{"an array"},
			"object_property_value": map[string]interface{}{
				"property_a": "a",
				"property_b": true,
			},
			"map_property_value": map[string]interface{}{
				"propertyA": "a",
				"propertyB": true,
				"propertyC": map[string]interface{}{
					"nestedPropertyA": true,
				},
			},
			"nested_resources": []interface{}{
				map[string]interface{}{
					"configuration": map[string]interface{}{
						"configurationValue": true,
					},
				},
			},
		},
		map[string]*schema.Schema{
			// Type mapPropertyValue as a map so that keys aren't mangled in the usual way.
			"map_property_value": {Type: schema.TypeMap},
			"nested_resources": {
				Type:     schema.TypeList,
				MaxItems: 1,
				// Embed a `*schema.Resource` to validate that type directed
				// walk of the schema successfully walks inside Resources as well
				// as Schemas.
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration": {Type: schema.TypeMap},
					},
				},
			},
		},
		map[string]*SchemaInfo{
			// Reverse map string_property_value to the stringo property.
			"string_property_value": {
				Name: "stringo",
			},
		},
		nil,   /* assets */
		false, /*useRawNames*/
	)
	assert.Equal(t, resource.NewPropertyMapFromMap(map[string]interface{}{
		"nilPropertyValue":    nil,
		"boolPropertyValue":   false,
		"numberPropertyValue": 42,
		"stringo":             "ognirts",
		"arrayPropertyValue":  []interface{}{"an array"},
		"objectPropertyValue": map[string]interface{}{
			"propertyA": "a",
			"propertyB": true,
		},
		"mapPropertyValue": map[string]interface{}{
			"propertyA": "a",
			"propertyB": true,
			"propertyC": map[string]interface{}{
				"nestedPropertyA": true,
			},
		},
		"nestedResources": []map[string]interface{}{{
			"configuration": map[string]interface{}{
				"configurationValue": true,
			},
		}},
	}), result)
}

func TestDefaults(t *testing.T) {
	// Produce maps with the following properties, and then validate them:
	//     - aaa string; no defaults, no inputs => empty
	//     - bbb string; no defaults, input "BBB" => "BBB"
	//     - ccc string; TF default "CCC", no inputs => "CCC"
	//     - cc2 string; TF default "CC2" (func), no inputs => "CC2"
	//     - ddd string; TF default "TFD", input "DDD" => "DDD"
	//     - dd2 string; TF default "TD2" (func), input "DDD" => "DDD"
	//     - eee string; PS default "EEE", no inputs => "EEE"
	//     - ee2 string; PS default "EE2" (func), no inputs => "EE2"
	//     - fff string; PS default "PSF", input "FFF" => "FFF"
	//     - ff2 string; PS default "PF2", input "FFF" => "FFF"
	//     - ggg string; TF default "TFG", PS default "PSG", no inputs => "PSG" (PS wins)
	//     - hhh string; TF default "TFH", PS default "PSH", input "HHH" => "HHH"
	//     - iii string; old default "OLI", TF default "TFI", PS default "PSI", no input => "OLD"
	//     - jjj string: old input "OLJ", no defaults, no input => no merged input
	//     - lll: old default "OLL", TF default "TFL", no input => "OLL"
	//     - mmm: old default "OLM", PS default "PSM", no input => "OLM"
	asset, err := resource.NewTextAsset("hello")
	assert.Nil(t, err)
	assets := make(AssetTable)
	tfs := map[string]*schema.Schema{
		"ccc": {Type: schema.TypeString, Default: "CCC"},
		"cc2": {Type: schema.TypeString, DefaultFunc: func() (interface{}, error) { return "CC2", nil }},
		"ddd": {Type: schema.TypeString, Default: "TFD"},
		"dd2": {Type: schema.TypeString, DefaultFunc: func() (interface{}, error) { return "TD2", nil }},
		"ggg": {Type: schema.TypeString, Default: "TFG"},
		"hhh": {Type: schema.TypeString, Default: "TFH"},
		"iii": {Type: schema.TypeString, Default: "TFI"},
		"jjj": {Type: schema.TypeString},
		"lll": {Type: schema.TypeString, Default: "TFL"},
		"mmm": {Type: schema.TypeString},
		"zzz": {Type: schema.TypeString},
	}
	ps := map[string]*SchemaInfo{
		"eee": {Default: &DefaultInfo{Value: "EEE"}},
		"ee2": {Default: &DefaultInfo{From: func(res *PulumiResource) (interface{}, error) { return "EE2", nil }}},
		"fff": {Default: &DefaultInfo{Value: "PSF"}},
		"ff2": {Default: &DefaultInfo{From: func(res *PulumiResource) (interface{}, error) { return "PF2", nil }}},
		"ggg": {Default: &DefaultInfo{Value: "PSG"}},
		"hhh": {Default: &DefaultInfo{Value: "PSH"}},
		"iii": {Default: &DefaultInfo{Value: "PSI"}},
		"mmm": {Default: &DefaultInfo{Value: "PSM"}},
		"zzz": {Asset: &AssetTranslation{Kind: FileAsset}},
	}
	olds := resource.PropertyMap{
		"iii": resource.NewStringProperty("OLI"),
		"jjj": resource.NewStringProperty("OLJ"),
		"lll": resource.NewStringProperty("OLL"),
		"mmm": resource.NewStringProperty("OLM"),
	}
	props := resource.PropertyMap{
		"bbb": resource.NewStringProperty("BBB"),
		"ddd": resource.NewStringProperty("DDD"),
		"dd2": resource.NewStringProperty("DDD"),
		"fff": resource.NewStringProperty("FFF"),
		"ff2": resource.NewStringProperty("FFF"),
		"hhh": resource.NewStringProperty("HHH"),
		"zzz": resource.NewAssetProperty(asset),
	}
	inputs, err := MakeTerraformInputs(nil, olds, props, tfs, ps, assets, true, false)
	assert.NoError(t, err)
	outputs := MakeTerraformOutputs(inputs, tfs, ps, assets, false)
	assert.Equal(t, resource.NewPropertyMapFromMap(map[string]interface{}{
		"bbb": "BBB",
		"ccc": "CCC",
		"cc2": "CC2",
		"ddd": "DDD",
		"dd2": "DDD",
		"eee": "EEE",
		"ee2": "EE2",
		"fff": "FFF",
		"ff2": "FFF",
		"ggg": "PSG",
		"hhh": "HHH",
		"iii": "OLI",
		"lll": "OLL",
		"mmm": "OLM",
		"zzz": asset,
	}), outputs)
}
