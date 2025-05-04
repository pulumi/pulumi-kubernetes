package condition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestEvalComplexConditions(t *testing.T) {
	tests := []struct {
		name      string
		condition interface{}
		obj       *unstructured.Unstructured
		expected  bool
		wantErr   bool
	}{
		{
			name: "simple string condition - true",
			condition: "jsonpath={.status.phase}=Running",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Running",
					},
				},
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "simple string condition - false",
			condition: "jsonpath={.status.phase}=Running",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Pending",
					},
				},
			},
			expected: false,
			wantErr:  false,
		},
		{
			name: "AND array - all true",
			condition: []interface{}{
				"jsonpath={.status.phase}=Running",
				"jsonpath={.status.ready}=true",
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Running",
						"ready": "true",
					},
				},
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "AND array - one false",
			condition: []interface{}{
				"jsonpath={.status.phase}=Running",
				"jsonpath={.status.ready}=true",
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Running",
						"ready": "false",
					},
				},
			},
			expected: false,
			wantErr:  false,
		},
		{
			name: "explicit AND - all true",
			condition: map[string]interface{}{
				"operator": "and",
				"conditions": []interface{}{
					"jsonpath={.status.phase}=Running",
					"jsonpath={.status.ready}=true",
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Running",
						"ready": "true",
					},
				},
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "explicit AND - one false",
			condition: map[string]interface{}{
				"operator": "and",
				"conditions": []interface{}{
					"jsonpath={.status.phase}=Running",
					"jsonpath={.status.ready}=true",
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Running",
						"ready": "false",
					},
				},
			},
			expected: false,
			wantErr:  false,
		},
		{
			name: "OR condition - one true",
			condition: map[string]interface{}{
				"operator": "or",
				"conditions": []interface{}{
					"jsonpath={.status.phase}=Running",
					"jsonpath={.status.phase}=Succeeded",
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Running",
					},
				},
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "OR condition - another true",
			condition: map[string]interface{}{
				"operator": "or",
				"conditions": []interface{}{
					"jsonpath={.status.phase}=Running",
					"jsonpath={.status.phase}=Succeeded",
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Succeeded",
					},
				},
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "OR condition - all false",
			condition: map[string]interface{}{
				"operator": "or",
				"conditions": []interface{}{
					"jsonpath={.status.phase}=Running",
					"jsonpath={.status.phase}=Succeeded",
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Pending",
					},
				},
			},
			expected: false,
			wantErr:  false,
		},
		{
			name: "nested AND-OR - true",
			condition: map[string]interface{}{
				"operator": "and",
				"conditions": []interface{}{
					"jsonpath={.status.initialized}=true",
					map[string]interface{}{
						"operator": "or",
						"conditions": []interface{}{
							"jsonpath={.status.phase}=Running",
							"jsonpath={.status.phase}=Succeeded",
						},
					},
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"initialized": "true",
						"phase": "Succeeded",
					},
				},
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "nested AND-OR - false (first condition)",
			condition: map[string]interface{}{
				"operator": "and",
				"conditions": []interface{}{
					"jsonpath={.status.initialized}=true",
					map[string]interface{}{
						"operator": "or",
						"conditions": []interface{}{
							"jsonpath={.status.phase}=Running",
							"jsonpath={.status.phase}=Succeeded",
						},
					},
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"initialized": "false",
						"phase": "Succeeded",
					},
				},
			},
			expected: false,
			wantErr:  false,
		},
		{
			name: "nested AND-OR - false (nested condition)",
			condition: map[string]interface{}{
				"operator": "and",
				"conditions": []interface{}{
					"jsonpath={.status.initialized}=true",
					map[string]interface{}{
						"operator": "or",
						"conditions": []interface{}{
							"jsonpath={.status.phase}=Running",
							"jsonpath={.status.phase}=Succeeded",
						},
					},
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"initialized": "true",
						"phase": "Pending",
					},
				},
			},
			expected: false,
			wantErr:  false,
		},
		{
			name: "nested OR-AND - true",
			condition: map[string]interface{}{
				"operator": "or",
				"conditions": []interface{}{
					"jsonpath={.status.phase}=Failed",
					map[string]interface{}{
						"operator": "and",
						"conditions": []interface{}{
							"jsonpath={.status.initialized}=true",
							"jsonpath={.status.ready}=true",
						},
					},
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"initialized": "true",
						"ready": "true",
						"phase": "Running",
					},
				},
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "invalid operator",
			condition: map[string]interface{}{
				"operator": "xor",
				"conditions": []interface{}{
					"jsonpath={.status.phase}=Running",
					"jsonpath={.status.phase}=Succeeded",
				},
			},
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			expected: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EvalCondition(tt.obj, tt.condition)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
