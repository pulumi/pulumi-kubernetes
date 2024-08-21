package condition

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCustom(t *testing.T) {
	stdout := logbuf{os.Stdout}

	tests := []struct {
		name          string
		expr          string
		obj           *unstructured.Unstructured
		wantSatisifed bool
	}{
		{
			name: "default True status not satisfied",
			expr: "condition=Custom",
			obj: &unstructured.Unstructured{
				Object: map[string]any{
					"status": map[string]any{
						"conditions": []any{
							map[string]any{
								"type":   "Custom",
								"status": "False",
							},
						},
					},
				},
			},
			wantSatisifed: false,
		},
		{
			name: "default True status satisfied",
			expr: "condition=Custom",
			obj: &unstructured.Unstructured{
				Object: map[string]any{
					"status": map[string]any{
						"conditions": []any{
							map[string]any{
								"type":   "Custom",
								"status": "True",
							},
						},
					},
				},
			},
			wantSatisifed: true,
		},
		{
			name: "default True status satisfied but condition's generation is too young",
			expr: "condition=Custom",
			obj: &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{
						"generation": int64(2),
					},
					"status": map[string]any{
						"conditions": []any{
							map[string]any{
								"type":               "Custom",
								"status":             "True",
								"observedGeneration": int64(1),
							},
						},
					},
				},
			},
			wantSatisifed: false,
		},
		{
			name: "default True status satisfied but status's generation is too young",
			expr: "condition=Custom",
			obj: &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{
						"generation": int64(2),
					},
					"status": map[string]any{
						"observedGeneration": int64(1),
						"conditions": []any{
							map[string]any{
								"type":   "Custom",
								"status": "True",
							},
						},
					},
				},
			},
			wantSatisifed: false,
		},
		{
			name: "custom status not satisfied",
			expr: "condition=Custom=foo",
			obj: &unstructured.Unstructured{
				Object: map[string]any{
					"status": map[string]any{
						"conditions": []any{
							map[string]any{
								"type":   "Custom",
								"status": "False",
							},
						},
					},
				},
			},
			wantSatisifed: false,
		},
		{
			name: "custom status satisfied",
			expr: "condition=Custom=Foo",
			obj: &unstructured.Unstructured{
				Object: map[string]any{
					"status": map[string]any{
						"conditions": []any{
							map[string]any{
								"type":   "Custom",
								"status": "foo",
							},
						},
					},
				},
			},
			wantSatisifed: true,
		},
		{
			name: "condition status missing",
			expr: "condition=Custom",
			obj: &unstructured.Unstructured{
				Object: map[string]any{
					"status": map[string]any{
						"conditions": []any{
							map[string]any{
								"type": "Custom",
							},
						},
					},
				},
			},
			wantSatisifed: false,
		},
		{
			name: "condition missing",
			expr: "condition=Custom",
			obj: &unstructured.Unstructured{
				Object: map[string]any{
					"status": map[string]any{
						"conditions": []any{
							map[string]any{
								"type":   "SomethingElse",
								"status": "foo",
							},
						},
					},
				},
			},
			wantSatisifed: false,
		},
		{
			name: "no conditions",
			expr: "condition=Custom",
			obj: &unstructured.Unstructured{
				Object: map[string]any{
					"status": map[string]any{},
				},
			},
			wantSatisifed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond, err := NewCustom(context.Background(), nil, stdout, tt.obj, tt.expr)
			require.NoError(t, err)

			actual, err := cond.Satisfied()
			require.NoError(t, err)

			assert.Equal(t, tt.wantSatisifed, actual)
		})
	}
}
