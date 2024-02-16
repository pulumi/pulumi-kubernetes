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

package provider

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
	"reflect"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	clientapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/yaml"
)

func TestHasComputedValue(t *testing.T) {
	tests := []struct {
		name             string
		obj              *unstructured.Unstructured
		hasComputedValue bool
	}{
		{
			name:             "nil object does not have a computed value",
			obj:              nil,
			hasComputedValue: false,
		},
		{
			name:             "Empty object does not have a computed value",
			obj:              &unstructured.Unstructured{},
			hasComputedValue: false,
		},
		{
			name:             "Object with no computed values does not have a computed value",
			obj:              &unstructured.Unstructured{Object: map[string]any{}},
			hasComputedValue: false,
		},
		{
			name: "Object with one concrete value does not have a computed value",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
			}},
			hasComputedValue: false,
		},
		{
			name: "Object with one computed value does have a computed value",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
				"field2": resource.Computed{},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with one nested computed value does have a computed value",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
				"field2": map[string]any{
					"field3": resource.Computed{},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with nested maps and no computed values",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
				"field2": map[string]any{
					"field3": "3",
				},
			}},
			hasComputedValue: false,
		},
		{
			name: "Object with doubly nested maps and 1 computed value",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
				"field2": map[string]any{
					"field3": "3",
					"field4": map[string]any{
						"field5": resource.Computed{},
					},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with nested slice of map[string]interface{} has a computed value",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
				"field2": []map[string]any{
					{"field3": resource.Computed{}},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with nested slice of interface{} has a computed value",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
				"field2": []any{
					resource.Computed{},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with nested slice of map[string]interface{} with nested slice of interface{} has a computed value",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
				"field2": []map[string]any{
					{"field3": []any{
						resource.Computed{},
					}},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Complex nested object with computed value",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
				"field2": []map[string]any{
					{"field3": []any{
						[]map[string]any{
							{"field4": []any{
								resource.Computed{},
							}},
						},
					}},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Complex nested object with no computed value",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"field1": 1,
				"field2": []map[string]any{
					{"field3": []any{
						[]map[string]any{
							{"field4": []any{
								"field5",
							}},
						},
					}},
				},
			}},
			hasComputedValue: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.hasComputedValue, hasComputedValue(test.obj), test.name)
	}
}

func TestAnnotateComputed(t *testing.T) {
	tests := []struct {
		name     string
		outValue resource.PropertyMap
		inValue  resource.PropertyMap
		expected resource.PropertyMap
	}{
		{
			name:     "out_notpresent",
			outValue: resource.PropertyMap{},
			inValue: resource.PropertyMap{
				"value": resource.NewStringProperty("foo"),
			},
			expected: resource.PropertyMap{},
		},
		{
			name: "out_computed",
			outValue: resource.PropertyMap{
				"value": resource.MakeComputed(resource.NewStringProperty("")),
			},
			inValue: resource.PropertyMap{
				"value": resource.MakeComputed(resource.NewStringProperty("")),
			},
			expected: resource.PropertyMap{
				"value": resource.MakeComputed(resource.NewStringProperty("")),
			},
		},
		{
			name: "out_object_null",
			outValue: resource.PropertyMap{
				"value": resource.NewObjectProperty(nil),
			},
			inValue: resource.PropertyMap{
				"value": resource.NewObjectProperty(resource.PropertyMap{}),
			},
			expected: resource.PropertyMap{
				"value": resource.NewObjectProperty(nil),
			},
		},
		{
			name: "in_value",
			outValue: resource.PropertyMap{
				"value": resource.NewStringProperty("foo"),
			},
			inValue: resource.PropertyMap{
				"value": resource.NewStringProperty("foo"),
			},
			expected: resource.PropertyMap{
				"value": resource.NewStringProperty("foo"),
			},
		},
		{
			name: "in_computed",
			outValue: resource.PropertyMap{
				"value": resource.NewStringProperty("foo"),
			},
			inValue: resource.PropertyMap{
				"value": resource.MakeComputed(resource.NewStringProperty("")),
			},
			expected: resource.PropertyMap{
				"value": resource.MakeComputed(resource.NewStringProperty("")),
			},
		},
		{
			name: "in_secret",
			outValue: resource.PropertyMap{
				"value": resource.NewStringProperty("foo"),
			},
			inValue: resource.PropertyMap{
				"value": resource.MakeSecret(resource.NewStringProperty("foo")),
			},
			expected: resource.PropertyMap{
				"value": resource.NewStringProperty("foo"),
			},
		},
		{
			name: "in_secret_computed",
			outValue: resource.PropertyMap{
				"value": resource.NewStringProperty("foo"),
			},
			inValue: resource.PropertyMap{
				"value": resource.MakeSecret(resource.MakeComputed(resource.NewStringProperty(""))),
			},
			expected: resource.PropertyMap{
				"value": resource.MakeComputed(resource.NewStringProperty("")),
			},
		},
		{
			name: "in_object",
			outValue: resource.PropertyMap{
				"value": resource.NewObjectProperty(resource.PropertyMap{
					"nested": resource.NewStringProperty("bar"),
				}),
			},
			inValue: resource.PropertyMap{
				"value": resource.NewObjectProperty(resource.PropertyMap{
					"nested": resource.MakeComputed(resource.NewStringProperty("")),
				}),
			},
			expected: resource.PropertyMap{
				"value": resource.NewObjectProperty(resource.PropertyMap{
					"nested": resource.MakeComputed(resource.NewStringProperty("")),
				}),
			},
		},
		{
			name: "in_object_null",
			outValue: resource.PropertyMap{
				"value": resource.NewObjectProperty(resource.PropertyMap{}),
			},
			inValue: resource.PropertyMap{
				"value": resource.NewObjectProperty(nil),
			},
			expected: resource.PropertyMap{
				"value": resource.NewObjectProperty(resource.PropertyMap{}),
			},
		},
		{
			name: "in_array",
			outValue: resource.PropertyMap{
				"value": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewStringProperty("foo"),
				}),
			},
			inValue: resource.PropertyMap{
				"value": resource.NewArrayProperty([]resource.PropertyValue{
					resource.MakeComputed(resource.NewStringProperty("")),
				}),
			},
			expected: resource.PropertyMap{
				"value": resource.NewArrayProperty([]resource.PropertyValue{
					resource.MakeComputed(resource.NewStringProperty("")),
				}),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			annotateComputed(tt.outValue, tt.inValue)
			assert.Equal(t, tt.expected, tt.outValue, "expected result")
		})
	}
}

func TestFqName(t *testing.T) {
	obj := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "tests/v1alpha1",
			"kind":       "Test",
			"metadata": map[string]any{
				"name": "myname",
			},
		},
	}

	if n := fqName(obj.GetNamespace(), obj.GetName()); n != "myname" {
		t.Errorf("Got %q for %v", n, obj)
	}

	obj.SetNamespace("mynamespace")
	if n := fqName(obj.GetNamespace(), obj.GetName()); n != "mynamespace/myname" {
		t.Errorf("Got %q for %v", n, obj)
	}
}

var (
	validKubeconfig = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeE1EUXlOekUxTkRjd05Wb1hEVE14TURReU5URTFORGN3TlZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTmVoCkJNOUowSVkrZFI5UmZVSjI5SlRxYjF2U3QwZUsxNDN1aVBxZElJR3hiWFFvVmV6ZDhRUXloSUFsUG91Z0VWS0gKUjRoTFRreEZJS01XQ1F0dGNCdkVaRnZBRmtyeVBzWU81RWgxRjZHdzJNbDYvNWtvU1psM1hTMDVyN1hnTUdWTQp5cVJRaDMvVWJFcVZkWkRlRWlBSnh6N3JQSUMxc1FUSlVqVTZUY2JaRFVYVkdGMVZMck9uRkJlUmg1NkwwN2RiCjJTeGd3dFhmNVNTMEFlYnJrT0REYzUwUUdYc250UkZONzE5YnlhblZCc3VzWm5mZnZIRWs1bnE1NUFMdGE0bjcKNkZDR2pRNHhYY2hsYTVvMWlreityN2pMenJ5NlNsdHJQWU5ML2VYNHgvRU0xUFFuVktlUWloRTJoNzRyakhLcApibDRwNjZPSjhseWRGa0VKQWVNQ0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZJWHlrY0tnZGI0SlhEM0tSelNKSG4rdlRCeXlNQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFBd1M1WjZvV056WnRhNE1EeWczNWJmcjVRaTlIdG5WN2c5Y2VOdmNUSStxd0d4VUhZUApnZzJSb1Q0ZU5MMkVQV1lBUmRrZVNtU3JYMGtFL1QydFRueGJFWENNdEI2TjhPRnZiQ3VnVzZPK1pwSDNwNHR0ClVFQ0UxT3ZiWHd5MkMvdmkyaXJuOWtEd3I3SkFVQ2FGRVppcmVPTWNDbGp6ZURNTDBDOUpqQlJOUmRqWHVscmIKSlRwL0RiWVJ0OFVJNW0zaVFIa2luakRHVkVhVHIzamVCTTZQakl1L25sakNlK1ovV0wyb3pFbzgydzN2cHpONQp2MGRvaHFkVmxPMzJnZDYrQlFORjhmUDI5bzBkT0NBalYvNHdCYmNjdUh1YnZCZ3U0cnFIc0hvZzM3MUVUdWwvCjlJbHFrZ2FmemVydVBzNms1UGFaUE1iK25nbzNZRG5ndkhuSAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://kubernetes.docker.internal:6443
  name: docker-desktop
contexts:
- context:
    cluster: docker-desktop
    user: docker-desktop
  name: docker-desktop
current-context: docker-desktop
kind: Config
preferences: {}
users:
- name: docker-desktop
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURGVENDQWYyZ0F3SUJBZ0lJWnBUZjVmbTZDQW93RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TVRBME1qY3hOVFEzTURWYUZ3MHlNakExTVRJeU16SXdNVEphTURZeApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sc3dHUVlEVlFRREV4SmtiMk5yWlhJdFptOXlMV1JsCmMydDBiM0F3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRRENMaDZ3MWV5WGpOQ3AKSzI5ZFJJQ3o3eHd3K1ZPVXVYYlh2R2NJTEFxaElUdlR3WTJqUmVaTFFXK3B5Wm9XUUdWZm5EYVZ5TGxmUUVaOQpXQm9IcEkvWGVvVWl4Uy9mWmVPN1RTeXA1bFpLcExzaXBMSE1RazN1NHp2d1RqelJITFJ4Q3k3b2RWTUVxVWFyCkkveUxiVUMxL1RkaGc3WkVZTFVrbEE4bWVhWFpHMGx5ZjA4UEdTdTVLUUJuTFVlbXk5OHNqV2U3YTBvdlRZd3kKTUhveUhyS0VGV0xCTmYrTm5TMTY3ZFBONzhTNCtENThobGxZTmZEZDVHbXJYYUFBYzVxeHhCSW5VcmhkSDJQawo2YUZkZXduQjFRQlV6OWVlVUJEVlFoQmwwbXNoMmRUSWF3cGlnbENrTnR1RlhoTEhGMitaRjFCSzN5VnZaYURsCmsyOTNnanlwQWdNQkFBR2pTREJHTUE0R0ExVWREd0VCL3dRRUF3SUZvREFUQmdOVkhTVUVEREFLQmdnckJnRUYKQlFjREFqQWZCZ05WSFNNRUdEQVdnQlNGOHBIQ29IVytDVnc5eWtjMGlSNS9yMHdjc2pBTkJna3Foa2lHOXcwQgpBUXNGQUFPQ0FRRUFnT1dxR2Q0TnlCRzFDOUovb1NVTmxzdkxSWXp4eEluZ0VsT09MUmlNN2t2dTduRXV6SHBYCkViODh3di9SSU1qWlFlbDFOTmdLWFJvb0hOSmpXcVppOG5aMEIxangySnNmaldrZWlPUE1aTjZqNzhzdDBqWmsKZDErSW5Oem1raEo4ck92UjVCd2xFcDNUcUtTN3J6dzF4MnphRkxUVWtZblh6Wnp4TkU3VGZuZVJVSG4zVyt4SwpXMFFaS3RkUlcvV0M5M3AvckcvZXp2Z0o5dCthenZwa2V1bklTUm5lbFpGQzgrZTR3ZXdoZm15TmRtUFVySThkCnQyMzhxeVhaNzZMTERKTFFDSTRieFRSVVpJM2NDdFY4bzU0UThnVHAvNklaRXVkV3dPbEdXa0FackdaMXNCN2QKQ1RRbjRVTVBXV0JmTzBPcFdZS3hGcVg4U1FpQndQaWhDdz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBd2k0ZXNOWHNsNHpRcVN0dlhVU0FzKzhjTVBsVGxMbDIxN3huQ0N3S29TRTcwOEdOCm8wWG1TMEZ2cWNtYUZrQmxYNXcybGNpNVgwQkdmVmdhQjZTUDEzcUZJc1V2MzJYanUwMHNxZVpXU3FTN0lxU3gKekVKTjd1TTc4RTQ4MFJ5MGNRc3U2SFZUQktsR3F5UDhpMjFBdGYwM1lZTzJSR0MxSkpRUEpubWwyUnRKY245UApEeGtydVNrQVp5MUhwc3ZmTEkxbnUydEtMMDJNTWpCNk1oNnloQlZpd1RYL2paMHRldTNUemUvRXVQZytmSVpaCldEWHczZVJwcTEyZ0FIT2FzY1FTSjFLNFhSOWo1T21oWFhzSndkVUFWTS9YbmxBUTFVSVFaZEpySWRuVXlHc0sKWW9KUXBEYmJoVjRTeHhkdm1SZFFTdDhsYjJXZzVaTnZkNEk4cVFJREFRQUJBb0lCQVFDeitSY05BMW1EcFVvSQpZVytZWEZPRmNnc0pBUzJNWE5GZlp3bC9zNEl1a2FUbndTOUxzdytkbElxd0xXQ1pXeG9hSWFrZDdxcVJNL3VoClZUVGEvSlV0UEN1RmJJblFYcGxTRWxkaEtWRzFZVFRwQ1FpWnJxS1kxUmZLeEZqdDM5TUdLejFReXQwbEp0ZU8KNjQyNGxJd3pvUHZoYjdoUmEraTRmRm9HYVIxa09KN1dGcFNwM2pUa0pZckFpQWViL2IxUlZ5Rk9sNm9IcEozcQo3dmxoaHZibklJcXdrMHp4VU1ya1ArTDR2azhLdjhEcVZZMTg1M1B5UWJ3cm1EUnBkNWYvTmRwZ0lrMWNjUURZCk1OUUxPd3NaRThsdTJZMW9PcTVpRmhZZEFmM2o2Ykt2Vi92WFAxdXhtdlZSMEZ6eWU0L2JuaXBUcWdNYUI1ZnQKQWJ5MVJsL2hBb0dCQU5vRnBLVExmTGYvQUFqNGw4TmlvdjJ2OXhVWmlrQVhNL0NPREpDVzVEU2hZakNZV0tDNApYNm8vdlJ6bHF3NWNaMDEvKzBYZ2lRa0tBQUdacXlzRnRsYjgveDdkYThnVWVDRVhGcDUzNlRiZXdHaXJlSlRoCkNCSnhlQ0x1cjdLVmp6RnpHdFZlTzY4VzRHRjg0ZmlQaUk0ZVJETlExN0pYVVR2cGtZTGNCQ2RsQW9HQkFPUUIKU05hdS9GYVdHVG5DbkVKa2pqRm9JdWFGbnVLaHphd1FIUHVFSFV6TFMzT1Fqem54MVRmMGU5aWxkRHBoekJ0SQpoNUgrbzFvUmhNYlN5Z2g1SGQ5aE1nekM3cjcrNmdPU2hMOEdnNjJwNU13YzhSVUhnZWhOWmkxSEJaeUh0VGFFCmg3LzA2YjBOV3lyMDRVcGNSZXJIME44TWdSWXI2emZ4K25MblpGWDFBb0dBSE1kLzYwejlJcUNqbFl1VEpQU0IKUlhHVDhSSVZBTTdQU1dMRzM5TTdQb05MSGRVT1pmRFFsLzJmN2crWEcrY3dyN2RFS1A0eHVLSzhTM25JY1g1bwppbVVOSERyb1Bsb05YWGpad0lOZG9xT1d6SHBPQ1lFRytzQkZ0bjdCYkpaM2QzU1ZSek1RTjlXU091d3NQQTVlClhUdzdqbmFPY25rNlBPbGhEdUFTSUUwQ2dZRUFuNGpHam5DaDMzUG04cU5ZOHB1cFlxaWF3dkY3MnRlY01XaVUKM3VmeUdHbW13WlhFb2FhMHFoSkhGYSt2UTZwcVJpelpyeTJjM3NpalB2citvaThjMTlBS1ZTT1FLZFB6cWN3NwpWZTRZOU1xTGJNWlRhWU4zUWpQbDZvaG5STDh2N0pXTzVxRlhheENOV2VFK1FlbU9nbGlOcllQeVRyRXNSRmpzCkJMb2pXb0VDZ1lCaWMrWjJvSzNTTmpzL3J5ZFU1Lzg3T3NVbExHamxKbDI1NE0xaGl3RmVsd3pUWjNXWjFuZlkKcS80Mm5GR3VRQUQ3RFFwSTBCSnFWVTJCQlZySmlSeFhROVlXUStCb3Q5VU4yRVJQQmhFeityU0Y0MnhybnZobApsZTU2NHVmK3VBdCt2K2ZjZUtYVnVDRDN1ZGdxL2d5ejNCaHN5VkJxZkFoNy9oNndOTmhIb3c9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`

	certAuthData = []byte(`-----BEGIN CERTIFICATE-----
MIIC5zCCAc+gAwIBAgIBADANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwprdWJl
cm5ldGVzMB4XDTIxMDQyNzE1NDcwNVoXDTMxMDQyNTE1NDcwNVowFTETMBEGA1UE
AxMKa3ViZXJuZXRlczCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANeh
BM9J0IY+dR9RfUJ29JTqb1vSt0eK143uiPqdIIGxbXQoVezd8QQyhIAlPougEVKH
R4hLTkxFIKMWCQttcBvEZFvAFkryPsYO5Eh1F6Gw2Ml6/5koSZl3XS05r7XgMGVM
yqRQh3/UbEqVdZDeEiAJxz7rPIC1sQTJUjU6TcbZDUXVGF1VLrOnFBeRh56L07db
2SxgwtXf5SS0AebrkODDc50QGXsntRFN719byanVBsusZnffvHEk5nq55ALta4n7
6FCGjQ4xXchla5o1ikz+r7jLzry6SltrPYNL/eX4x/EM1PQnVKeQihE2h74rjHKp
bl4p66OJ8lydFkEJAeMCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMA8GA1UdEwEB
/wQFMAMBAf8wHQYDVR0OBBYEFIXykcKgdb4JXD3KRzSJHn+vTByyMA0GCSqGSIb3
DQEBCwUAA4IBAQAwS5Z6oWNzZta4MDyg35bfr5Qi9HtnV7g9ceNvcTI+qwGxUHYP
gg2RoT4eNL2EPWYARdkeSmSrX0kE/T2tTnxbEXCMtB6N8OFvbCugW6O+ZpH3p4tt
UECE1OvbXwy2C/vi2irn9kDwr7JAUCaFEZireOMcCljzeDML0C9JjBRNRdjXulrb
JTp/DbYRt8UI5m3iQHkinjDGVEaTr3jeBM6PjIu/nljCe+Z/WL2ozEo82w3vpzN5
v0dohqdVlO32gd6+BQNF8fP29o0dOCAjV/4wBbccuHubvBgu4rqHsHog371ETul/
9IlqkgafzeruPs6k5PaZPMb+ngo3YDngvHnH
-----END CERTIFICATE-----
`)
)

func writeKubeConfig(t *testing.T, dir string, config *clientcmdapi.Config) string {
	f, err := os.CreateTemp(dir, "kubeconfig-")
	if err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
	if err = clientcmd.WriteToFile(*config, f.Name()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err = os.Remove(f.Name()); err != nil {
			t.Fatal(err)
		}
	})
	return f.Name()
}

func Test_parseKubeconfigString(t *testing.T) {
	// write out a valid kubeconfig into the home directory
	usr, err := user.Current()
	assert.NoError(t, err)
	validConfig, _ := clientcmd.Load([]byte(validKubeconfig))
	validConfigFile := writeKubeConfig(t, usr.HomeDir, validConfig)

	tests := []struct {
		name           string
		pathOrContents string
		want           *clientapi.Config
		wantErr        string
	}{
		{
			name:           "empty",
			pathOrContents: "",
			want:           &clientapi.Config{},
		},
		{
			name:           "invalid",
			pathOrContents: "invalid",
			wantErr:        "json parse error",
		},
		{
			name:           "string",
			pathOrContents: validKubeconfig,
			want:           validConfig,
		},
		{
			name:           "home_dir",
			pathOrContents: "~",
			wantErr:        "is a directory",
		},
		{
			name:           "home_file",
			pathOrContents: path.Join("~", path.Base(validConfigFile)),
			want:           validConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseKubeconfigString(tt.pathOrContents)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else if assert.NoError(t, err) {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parseKubeconfigString() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_parseKubeconfigPropertyValue(t *testing.T) {
	validConfig, _ := clientcmd.Load([]byte(validKubeconfig))

	// convert the kubeconfig (YAML) to a JSON representation similar to:
	// https://github.com/pulumi/pulumi-eks/blob/160a7531e7377d6e8e41d3dbc62d2cee8703e630/nodejs/eks/cluster.ts#L199-L261
	validConfigJSON, _ := yaml.YAMLToJSON([]byte(validKubeconfig))
	validConfigObject := make(map[string]interface{})
	_ = json.Unmarshal(validConfigJSON, &validConfigObject)

	tests := []struct {
		name       string
		kubeconfig resource.PropertyValue
		want       *clientapi.Config
		wantErr    string
	}{
		{
			name:       "null",
			kubeconfig: resource.NewNullProperty(),
			want:       &clientapi.Config{},
		},
		{
			name:       "invalid",
			kubeconfig: resource.NewStringProperty("invalid"),
			wantErr:    "json parse error",
		},
		{
			name:       "empty",
			kubeconfig: resource.NewStringProperty(""),
			want:       &clientapi.Config{},
		},
		{
			name:       "string",
			kubeconfig: resource.NewStringProperty(validKubeconfig),
			want:       validConfig,
		},
		{
			name:       "object",
			kubeconfig: resource.NewObjectProperty(resource.NewPropertyMapFromMap(validConfigObject)),
			want:       validConfig,
		},
		{
			name:       "unsupported_resource_type",
			kubeconfig: resource.NewBoolProperty(true),
			wantErr:    "unexpected kubeconfig format",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseKubeconfigPropertyValue(tt.kubeconfig)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else if assert.NoError(t, err) {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parseKubeconfigPropertyValue() = %v, want %v", got, tt.want)
				}
			}
		})
	}

}

func Test_getActiveClusterFromConfig(t *testing.T) {
	validConfig, _ := clientcmd.Load([]byte(validKubeconfig))

	type args struct {
		config    *clientapi.Config
		overrides resource.PropertyMap
	}
	tests := []struct {
		name string
		args args
		want *clientapi.Cluster
		ok   bool
	}{
		{
			name: "nil",
			args: args{
				config:    nil,
				overrides: map[resource.PropertyKey]resource.PropertyValue{},
			},
			want: &clientapi.Cluster{},
			ok:   false,
		},
		{
			name: "empty",
			args: args{
				config:    &clientapi.Config{},
				overrides: map[resource.PropertyKey]resource.PropertyValue{},
			},
			want: &clientapi.Cluster{},
			ok:   false,
		},
		{
			name: "valid",
			args: args{
				config:    validConfig,
				overrides: map[resource.PropertyKey]resource.PropertyValue{},
			},
			want: &clientapi.Cluster{
				Server:                   "https://kubernetes.docker.internal:6443",
				CertificateAuthorityData: certAuthData,
				Extensions:               map[string]runtime.Object{},
			},
			ok: true,
		},
		{
			name: "invalid_context_override",
			args: args{
				config: validConfig,
				overrides: map[resource.PropertyKey]resource.PropertyValue{
					resource.PropertyKey("context"): {V: "foo"},
				},
			},
			want: &clientapi.Cluster{},
			ok:   false,
		},
		{
			name: "invalid_cluster_override",
			args: args{
				config: validConfig,
				overrides: map[resource.PropertyKey]resource.PropertyValue{
					resource.PropertyKey("cluster"): {V: "foo"},
				},
			},
			want: &clientapi.Cluster{},
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := getActiveClusterFromConfig(tt.args.config, tt.args.overrides)
			if ok != tt.ok {
				t.Errorf("getActiveClusterFromConfig() = %v, ok %v", ok, tt.ok)
			}
			if ok && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getActiveClusterFromConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPruneMap(t *testing.T) {
	oldLiveJSON := []byte(`{
  "__fieldManager": "pulumi-kubernetes",
  "__initialApiVersion": "apps/v1",
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "annotations": {
      "deployment.kubernetes.io/revision": "1",
      "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"apps/v1\",\"kind\":\"Deployment\",\"metadata\":{\"labels\":{\"app.kubernetes.io/managed-by\":\"pulumi\"},\"name\":\"scale-test\",\"namespace\":\"default\"},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"app\":\"nginx\"}},\"template\":{\"metadata\":{\"labels\":{\"app\":\"nginx\"}},\"spec\":{\"containers\":[{\"image\":\"nginx:1.13\",\"name\":\"nginx\",\"ports\":[{\"containerPort\":80}],\"securityContext\":{\"capabilities\":{\"add\":[\"NET_ADMIN\",\"SYS_TIME\"]}}}]}}}}\n"
    },
    "creationTimestamp": "2023-06-02T15:33:42Z",
    "generation": 1,
    "labels": {
      "app.kubernetes.io/managed-by": "pulumi"
    },
    "managedFields": [
      {
        "apiVersion": "apps/v1",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:metadata": {
            "f:annotations": {
              ".": {},
              "f:kubectl.kubernetes.io/last-applied-configuration": {}
            },
            "f:labels": {
              ".": {},
              "f:app.kubernetes.io/managed-by": {}
            }
          },
          "f:spec": {
            "f:progressDeadlineSeconds": {},
            "f:replicas": {},
            "f:revisionHistoryLimit": {},
            "f:selector": {},
            "f:strategy": {
              "f:rollingUpdate": {
                ".": {},
                "f:maxSurge": {},
                "f:maxUnavailable": {}
              },
              "f:type": {}
            },
            "f:template": {
              "f:metadata": {
                "f:labels": {
                  ".": {},
                  "f:app": {}
                }
              },
              "f:spec": {
                "f:containers": {
                  "k:{\"name\":\"nginx\"}": {
                    ".": {},
                    "f:image": {},
                    "f:imagePullPolicy": {},
                    "f:name": {},
                    "f:ports": {
                      ".": {},
                      "k:{\"containerPort\":80,\"protocol\":\"TCP\"}": {
                        ".": {},
                        "f:containerPort": {},
                        "f:protocol": {}
                      }
                    },
                    "f:resources": {},
                    "f:securityContext": {
                      ".": {},
                      "f:capabilities": {
                        ".": {},
                        "f:add": {}
                      }
                    },
                    "f:terminationMessagePath": {},
                    "f:terminationMessagePolicy": {}
                  }
                },
                "f:dnsPolicy": {},
                "f:restartPolicy": {},
                "f:schedulerName": {},
                "f:securityContext": {},
                "f:terminationGracePeriodSeconds": {}
              }
            }
          }
        },
        "manager": "pulumi-kubernetes",
        "operation": "Update",
        "time": "2023-06-02T15:33:42Z"
      },
      {
        "apiVersion": "apps/v1",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:metadata": {
            "f:annotations": {
              "f:deployment.kubernetes.io/revision": {}
            }
          },
          "f:status": {
            "f:availableReplicas": {},
            "f:conditions": {
              ".": {},
              "k:{\"type\":\"Available\"}": {
                ".": {},
                "f:lastTransitionTime": {},
                "f:lastUpdateTime": {},
                "f:message": {},
                "f:reason": {},
                "f:status": {},
                "f:type": {}
              },
              "k:{\"type\":\"Progressing\"}": {
                ".": {},
                "f:lastTransitionTime": {},
                "f:lastUpdateTime": {},
                "f:message": {},
                "f:reason": {},
                "f:status": {},
                "f:type": {}
              }
            },
            "f:observedGeneration": {},
            "f:readyReplicas": {},
            "f:replicas": {},
            "f:updatedReplicas": {}
          }
        },
        "manager": "kube-controller-manager",
        "operation": "Update",
        "subresource": "status",
        "time": "2023-06-02T15:33:43Z"
      }
    ],
    "name": "scale-test",
    "namespace": "default",
    "resourceVersion": "6639864",
    "uid": "12990b2a-5476-4f4c-be63-cb5fcafe0cf2"
  },
  "spec": {
    "progressDeadlineSeconds": 600,
    "replicas": 1,
    "revisionHistoryLimit": 10,
    "selector": {
      "matchLabels": {
        "app": "nginx"
      }
    },
    "strategy": {
      "rollingUpdate": {
        "maxSurge": "25%",
        "maxUnavailable": "25%"
      },
      "type": "RollingUpdate"
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "nginx"
        }
      },
      "spec": {
        "containers": [
          {
            "image": "nginx:1.13",
            "imagePullPolicy": "IfNotPresent",
            "name": "nginx",
            "ports": [
              {
                "containerPort": 80,
                "protocol": "TCP"
              }
            ],
            "resources": {},
            "securityContext": {
              "capabilities": {
                "add": [
                  "NET_ADMIN",
                  "SYS_TIME"
                ]
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File"
          }
        ],
        "dnsPolicy": "ClusterFirst",
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "terminationGracePeriodSeconds": 30
      }
    }
  },
  "status": {
    "availableReplicas": 1,
    "conditions": [
      {
        "lastTransitionTime": "2023-06-02T15:33:43Z",
        "lastUpdateTime": "2023-06-02T15:33:43Z",
        "message": "Deployment has minimum availability.",
        "reason": "MinimumReplicasAvailable",
        "status": "True",
        "type": "Available"
      },
      {
        "lastTransitionTime": "2023-06-02T15:33:42Z",
        "lastUpdateTime": "2023-06-02T15:33:43Z",
        "message": "ReplicaSet \"scale-test-544b74d7f9\" has successfully progressed.",
        "reason": "NewReplicaSetAvailable",
        "status": "True",
        "type": "Progressing"
      }
    ],
    "observedGeneration": 1,
    "readyReplicas": 1,
    "replicas": 1,
    "updatedReplicas": 1
  }
}`)

	oldInputsJSON := []byte(`{
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "labels": {
      "app.kubernetes.io/managed-by": "pulumi"
    },
    "name": "scale-test",
    "namespace": "default"
  },
  "spec": {
    "replicas": 1,
    "selector": {
      "matchLabels": {
        "app": "nginx"
      }
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "nginx"
        }
      },
      "spec": {
        "containers": [
          {
            "image": "nginx:1.13",
            "name": "nginx",
            "ports": [
              {
                "containerPort": 80
              }
            ],
            "securityContext": {
              "capabilities": {
                "add": [
                  "NET_ADMIN",
                  "SYS_TIME"
                ]
              }
            }
          }
        ]
      }
    }
  }
}`)

	var err error
	var source, target map[string]any
	err = json.Unmarshal(oldInputsJSON, &target)
	assert.NoErrorf(t, err, "failed to unmarshal oldInputsJSON")
	err = json.Unmarshal(oldLiveJSON, &source)
	assert.NoErrorf(t, err, "failed to unmarshal oldLiveJSON")

	type args struct {
		source map[string]any
		target map[string]any
	}
	tests := []struct {
		name        string
		description string
		args        args
		want        map[string]any
	}{
		{
			name:        "empty target",
			description: "empty target map should result in empty result map",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": "b",
				},
				target: map[string]any{},
			},
			want: map[string]any{},
		},
		{
			name:        "empty source",
			description: "empty source map should result in empty result map",
			args: args{
				source: map[string]any{},
				target: map[string]any{
					"a": "a",
					"b": "b",
				},
			},
			want: map[string]any{},
		},
		{
			name:        "matching keys with different values",
			description: "a map where target has matching keys and different values",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": "b",
				},
				target: map[string]any{
					"a": "A",
					"b": "B",
				},
			},
			want: map[string]any{
				"a": "a",
				"b": "b",
			},
		},
		{
			name:        "matching keys with nil source value",
			description: "a map where target has matching keys and source has a nil value",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": nil,
				},
				target: map[string]any{
					"a": "A",
					"b": "B",
				},
			},
			want: map[string]any{
				"a": "a",
				"b": nil,
			},
		},
		{
			name:        "matching keys with nil target value",
			description: "a map where target has matching keys and target has a nil value",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": "b",
				},
				target: map[string]any{
					"a": "A",
					"b": nil,
				},
			},
			want: map[string]any{
				"a": "a",
				"b": "b",
			},
		},
		{
			name:        "matching keys but different value types",
			description: "a map where target has matching keys and different value types",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": "b",
				},
				target: map[string]any{
					"a": "A",
					"b": 2,
				},
			},
			want: map[string]any{
				"a": "a",
				"b": "b", // key is present in target, so keep it even though the value type doesn't match
			},
		},
		{
			name:        "simple map",
			description: "a map where target matches",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": "b",
				},
				target: map[string]any{
					"a": "a",
					"b": "b",
				},
			},
			want: map[string]any{
				"a": "a",
				"b": "b",
			},
		},
		{
			name:        "simple map subset",
			description: "a map where target is a subset",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": "b", // not present in target, so will be ignored
				},
				target: map[string]any{
					"a": "A",
				},
			},
			want: map[string]any{
				"a": "a",
			},
		},
		{
			name:        "simple map superset",
			description: "a map where target is a superset",
			args: args{
				source: map[string]any{
					"a": "a",
				},
				target: map[string]any{
					"a": "A",
					"b": "B", // the extra key will be ignored if not present in source
				},
			},
			want: map[string]any{
				"a": "a",
			},
		},
		{
			name:        "nested map",
			description: "a map with a nested map where target matches",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": map[string]any{
						"c": "c",
					},
				},
				target: map[string]any{
					"a": "a",
					"b": map[string]any{
						"c": "c",
					},
				},
			},
			want: map[string]any{
				"a": "a",
				"b": map[string]any{
					"c": "c",
				},
			},
		},
		{
			name:        "nested map subset",
			description: "a map with a nested map where target is a subset",
			args: args{
				source: map[string]any{
					"a": "a", // not present in target, so will be ignored
					"b": map[string]any{
						"c": "c",
					},
				},
				target: map[string]any{
					"b": map[string]any{
						"c": "C",
					},
				},
			},
			want: map[string]any{
				"b": map[string]any{
					"c": "c",
				},
			},
		},
		{
			name:        "nested map superset",
			description: "a map with a nested map where target is a superset",
			args: args{
				source: map[string]any{
					"a": "a", // not present in target, so will be ignored
					"b": map[string]any{
						"c": "c",
					},
				},
				target: map[string]any{
					"b": map[string]any{
						"c": "C",
					},
					"d": "D", // the extra key will be ignored if not present in source
				},
			},
			want: map[string]any{
				"b": map[string]any{
					"c": "c",
				},
			},
		},
		{
			name:        "nested map with nil",
			description: "a map with a nested map with nil where target matches",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": map[string]any{
						"c": nil,
					},
				},
				target: map[string]any{
					"a": "a",
					"b": map[string]any{
						"c": nil,
					},
				},
			},
			want: map[string]any{
				"a": "a",
				"b": map[string]any{
					"c": nil,
				},
			},
		},
		{
			name:        "nested empty map",
			description: "a map with an empty nested map where target matches",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": map[string]any{},
				},
				target: map[string]any{
					"a": "a",
					"b": map[string]any{},
				},
			},
			want: map[string]any{
				"a": "a",
				"b": map[string]any{},
			},
		},
		{
			name:        "nested value slice",
			description: "a map with a nested slice of simple values where target matches",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": []any{"c", "d"},
				},
				target: map[string]any{
					"a": "a",
					"b": []any{"c", "d"},
				},
			},
			want: map[string]any{
				"a": "a",
				"b": []any{"c", "d"},
			},
		},
		{
			name:        "nested value slice subset",
			description: "a map with a nested slice of simple values where target is a subset",
			args: args{
				source: map[string]any{
					"a": "a", // not present in target, so will be ignored
					"b": []any{"c", "d"},
				},
				target: map[string]any{
					"b": []any{"c"},
				},
			},
			want: map[string]any{
				"b": []any{"c"}, // items compared by index, so only the first item in source will be kept
			},
		},
		{
			name:        "nested value slice superset",
			description: "a map with a nested slice of simple values where target is a superset",
			args: args{
				source: map[string]any{
					"a": "a", // not present in target, so will be ignored
					"b": []any{"c", "d"},
				},
				target: map[string]any{
					"b": []any{"c", "d", "e"},
				},
			},
			want: map[string]any{
				"b": []any{"c", "d"},
			},
		},
		{
			name:        "nested empty slice",
			description: "a map with an empty nested slice of simple values where target matches",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": []any{},
				},
				target: map[string]any{
					"a": "a",
					"b": []any{},
				},
			},
			want: map[string]any{
				"a": "a",
				"b": []any{},
			},
		},
		{
			name:        "nested map slice",
			description: "a map with a nested slice of map values where target matches",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": []any{
						map[string]any{
							"c": "c",
							"d": "d",
						},
					},
				},
				target: map[string]any{
					"a": "a",
					"b": []any{
						map[string]any{
							"c": "c",
							"d": "d",
						},
					},
				},
			},
			want: map[string]any{
				"a": "a",
				"b": []any{
					map[string]any{
						"c": "c",
						"d": "d",
					},
				},
			},
		},
		{
			name:        "nested map slice subset",
			description: "a map with a nested slice of map values where target is a subset",
			args: args{
				source: map[string]any{
					"a": "a", // not present in target, so will be ignored
					"b": []any{
						map[string]any{
							"c": "c",
							"d": "d", // not present in target, so will be ignored
						},
					},
				},
				target: map[string]any{
					"b": []any{
						map[string]any{
							"c": "c",
						},
					},
				},
			},
			want: map[string]any{
				"b": []any{
					map[string]any{
						"c": "c",
					},
				},
			},
		},
		{
			name:        "nested map slice superset",
			description: "a map with a nested slice of map values where target is a superset",
			args: args{
				source: map[string]any{
					"a": "a", // not present in target, so will be ignored
					"b": []any{
						map[string]any{
							"c": "c",
							"d": "d",
						},
					},
				},
				target: map[string]any{
					"b": []any{
						map[string]any{
							"c": "c",
							"d": "d",
							"e": "e",
						},
					},
				},
			},
			want: map[string]any{
				"b": []any{
					map[string]any{
						"c": "c",
						"d": "d",
					},
				},
			},
		},
		{
			name:        "real data",
			description: "a complex source and target using real data from the provider",
			args: args{
				source: source,
				target: target,
			},
			want: target,
		},
		{
			name:        "empty nil map",
			description: "nil map should result in nil result",
			args: args{
				source: nil,
				target: nil,
			},
			want: nil,
		},
		{
			name:        "empty nil source map",
			description: "nil source map should result in nil result",
			args: args{
				source: nil,
				target: map[string]any{
					"a": "a",
					"b": "b",
				},
			},
			want: nil,
		},
		{
			name:        "empty nil target map",
			description: "nil target map should result in nil result",
			args: args{
				source: map[string]any{
					"a": "a",
					"b": "b",
				},
				target: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pruneMap(tt.args.source, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nDescription: %s\nExpected:    %+v\nActual:      %+v", tt.description, tt.want, got)
			}
		})
	}
}

func TestPruneSlice(t *testing.T) {
	type args struct {
		source []any
		target []any
	}
	tests := []struct {
		name        string
		description string
		args        args
		want        []any
	}{
		{
			name:        "empty target",
			description: "empty target slice should result in empty result slice",
			args: args{
				source: []any{"a", "b"},
				target: []any{},
			},
			want: []any{},
		},
		{
			name:        "empty source",
			description: "empty source slice should result in empty result slice",
			args: args{
				source: []any{},
				target: []any{"a", "b"},
			},
			want: []any{},
		},
		{
			name:        "nil target",
			description: "nil target slice should result in nil result",
			args: args{
				source: []any{"a", "b"},
				target: nil,
			},
			want: nil,
		},
		{
			name:        "nil source",
			description: "nil source slice should result in nil result",
			args: args{
				source: nil,
				target: []any{"a", "b"},
			},
			want: nil,
		},
		{
			name:        "matching number of elements with different values",
			description: "a slice where target has matching number of elements with different values",
			args: args{
				source: []any{"a", "b"},
				target: []any{"c", "d"},
			},
			want: []any{"a", "b"},
		},
		{
			name:        "matching number of elements but different types",
			description: "a slice where target has matching number of elements with different types",
			args: args{
				source: []any{"a", "b"},
				target: []any{1, 2},
			},
			want: []any{"a", "b"},
		},
		{
			name:        "simple slice",
			description: "a slice where target matches",
			args: args{
				source: []any{"a", "b"},
				target: []any{"a", "b"},
			},
			want: []any{"a", "b"},
		},
		{
			name:        "simple slice subset",
			description: "a slice where target is a subset",
			args: args{
				source: []any{"a", "b"},
				target: []any{"a"},
			},
			want: []any{"a"},
		},
		{
			name:        "simple slice superset",
			description: "a slice where target is a superset",
			args: args{
				source: []any{"a"},
				target: []any{"a", "b"},
			},
			want: []any{"a"},
		},
		{
			name:        "map slice",
			description: "a slice of map values where target matches",
			args: args{
				source: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
				},
				target: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
				},
			},
			want: []any{
				map[string]any{
					"a": "a",
					"b": "b",
				},
			},
		},
		{
			name:        "map slice subset",
			description: "a slice of map values where target is a subset",
			args: args{
				source: []any{
					map[string]any{
						"a": "a",
						"b": "b", // not present in target, so will be dropped
					},
					map[string]any{
						"c": "c",
						"d": "d",
					},
				},
				target: []any{
					map[string]any{
						"a": "a",
					},
				},
			},
			want: []any{
				map[string]any{
					"a": "a",
				},
			},
		},
		{
			name:        "map slice superset",
			description: "a slice of map values where target is a superset",
			args: args{
				source: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
				},
				target: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
					map[string]any{
						"c": "c",
						"d": "d",
					},
				},
			},
			want: []any{
				map[string]any{
					"a": "a",
					"b": "b",
				},
			},
		},
		{
			name:        "map slice with nil value",
			description: "a slice of map values that include a nil value",
			args: args{
				source: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
					nil,
				},
				target: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
					nil,
				},
			},
			want: []any{
				map[string]any{
					"a": "a",
					"b": "b",
				},
				nil,
			},
		},
		{
			name:        "map slice with empty non-nil value",
			description: "a slice of map values that include an empty non-nil value",
			args: args{
				source: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
					map[string]any{},
				},
				target: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
					map[string]any{},
				},
			},
			want: []any{
				map[string]any{
					"a": "a",
					"b": "b",
				},
				map[string]any{},
			},
		},
		{
			name:        "map slice with empty non-nil value in target",
			description: "a slice of map values that include an empty non-nil value only in target",
			args: args{
				source: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
					map[string]any{
						"a": "a",
						"b": "b",
					},
				},
				target: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
					map[string]any{},
				},
			},
			want: []any{
				map[string]any{
					"a": "a",
					"b": "b",
				},
				map[string]any{},
			},
		},
		{
			name:        "map slice with empty non-nil value in source",
			description: "a slice of map values that include an empty non-nil value only in source",
			args: args{
				source: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
					map[string]any{},
				},
				target: []any{
					map[string]any{
						"a": "a",
						"b": "b",
					},
					map[string]any{
						"a": "a",
						"b": "b",
					},
				},
			},
			want: []any{
				map[string]any{
					"a": "a",
					"b": "b",
				},
				map[string]any{},
			},
		},
		{
			name:        "nil slice",
			description: "nil slice should return nil",
			args: args{
				source: nil,
				target: nil,
			},
			want: nil,
		},
		{
			name:        "non-nil empty slice",
			description: "non-nil empty slice should return empty slice",
			args: args{
				source: []any{},
				target: []any{},
			},
			want: []any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pruneSlice(tt.args.source, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nDescription: %s\nExpected:    %+v\nActual:      %+v", tt.description, tt.want, got)
			}
		})
	}
}
