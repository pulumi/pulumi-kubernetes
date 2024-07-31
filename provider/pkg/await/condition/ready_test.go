package condition

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestReady(t *testing.T) {
	stdout := logbuf{os.Stdout}

	tests := []struct {
		name      string
		given     string
		wantReady bool
	}{
		{
			name: "PVC bound",
			given: `{
				"apiVersion": "v1",
				"kind": "PersistentVolumeClaim",
				"status": {
					"phase": "Bound"
				}
			}`,
			wantReady: true,
		},
		{
			name: "PVC unbound",
			given: `{
				"apiVersion": "v1",
				"kind": "PersistentVolumeClaim",
				"status": {
					"phase": "NotBound"
				}
			}`,
			wantReady: false,
		},
		{
			name: "missing status field is assumed to be ready",
			given: `{ 
				"apiVersion": "test.pulumi.com/v1",
				"kind": "ResourceWithoutStatus",
			}`,
			wantReady: true,
		},
		{
			name: "status without conditions is assumed to be ready",
			given: `{ 
				"apiVersion": "test.pulumi.com/v1",
				"kind": "ResourceWithoutConditions",
				"status": {
					"foo": "bar"
				}
			}`,
			wantReady: true,
		},
		{
			name: "status with empty conditions is assumed to be ready",
			given: `{ 
				"apiVersion": "test.pulumi.com/v1",
				"kind": "ResourceWithEmptyConditions",
				"status": {
					"conditions": []
				}
			}`,
			wantReady: true,
		},
		{
			name: "status without ready condition is assumed to be ready",
			given: `{
				"apiVersion": "test.pulumi.com/v1",
				"kind": "SomeThing",
				"status": {
					"conditions": [{
						"type": "SomethingOtherThanReady",
						"status": "False"
					}]
				}
			}`,
			wantReady: true,
		},
		{
			name: "unknown ready condition isn't ready",
			given: `{
				"apiVersion": "test.pulumi.com/v1",
				"kind": "SomeThing",
				"status": {
					"conditions": [{
						"type": "Ready",
						"status": "Unknown"
					}]
				}
			}`,
			wantReady: false,
		},
		{
			name: "non-standard ready status is assumed to be ready",
			given: `{
				"apiVersion": "test.pulumi.com/v1",
				"kind": "SomeThing",
				"status": {
					"conditions": [{
						"type": "Ready",
						"status": "NotTrueOrFalseOrUnknown"
					}]
				}
			}`,
			wantReady: true,
		},
		{
			name: "CRD progressing",
			given: `{
				apiVersion: apiextensions.k8s.io/v1,
				kind: CustomResourceDefinition,
				status: {
				  conditions: [
					{
						message: no conflicts found,
						reason: NoConflicts,
						status: "True",
						type: NamesAccepted,
					}
				  ]
				}
			}`,
			wantReady: false,
		},
		{
			name: "CRD ready",
			given: `{
				apiVersion: apiextensions.k8s.io/v1,
				kind: CustomResourceDefinition,
				status: {
				  conditions: [
					{
						message: no conflicts found,
						reason: NoConflicts,
						status: "True",
						type: NamesAccepted,
					},
					{
						message: the initial names have been accepted,
						reason: InitialNamesAccepted,
						status: "True",
						type: Established,
					}
				  ]
				}
			}`,
			wantReady: true,
		},
		{
			name: "cert-manager progressing",
			given: `{
				"apiVersion": "cert-manager.io/v1",
				"kind": "Issuer",
				"status": {
					"conditions": [{
						"type": "Ready",
						"status": "False"
					}]
				}
			}`,
			wantReady: false,
		},
		{
			name: "cert-manager ready",
			given: `{
				"apiVersion": "cert-manager.io/v1",
				"kind": "Issuer",
				"status": {
					"conditions": [{
						"type": "Ready",
						"status": "True"
					}]
				}
			}`,
			wantReady: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj map[string]any
			err := yaml.Unmarshal([]byte(tt.given), &obj)
			require.NoError(t, err)

			ready := NewReady(context.Background(), Static(nil), stdout, &unstructured.Unstructured{Object: obj})

			actual, err := ready.Satisfied()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantReady, actual)
		})
	}
}
