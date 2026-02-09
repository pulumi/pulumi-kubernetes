package condition

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/jsonpath"
)

func TestJSONPath(t *testing.T) {
	stdout := logbuf{os.Stdout}

	tests := []struct {
		name      string
		uns       *unstructured.Unstructured
		expr      string
		wantReady bool
		wantErr   string
	}{
		{
			name:      "missing key",
			uns:       &unstructured.Unstructured{Object: map[string]any{}},
			expr:      "jsonpath={.foo}",
			wantReady: false,
		},
		{
			name: "key present",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": "bar",
			}},
			expr:      "jsonpath={.foo}",
			wantReady: true,
		},
		{
			name: "key present but empty",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": "",
			}},
			expr: "jsonpath={.foo}",
			// Ref:
			// https://github.com/kubernetes/kubectl/blob/c4be63c54b7188502c1a63bb884a0b05fac51ebd/pkg/cmd/wait/json.go#L72-L91
			wantReady: true,
		},
		{
			name: "key present but null",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": nil,
			}},
			expr: "jsonpath={.foo}",
			// Ref:
			// https://github.com/kubernetes/kubectl/blob/c4be63c54b7188502c1a63bb884a0b05fac51ebd/pkg/cmd/wait/json.go#L72-L91
			wantReady: true,
		},
		{
			name: "key present but observed generation is too young",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": "bar",
				"metadata": map[string]any{
					"generation": int64(2),
				},
				"status": map[string]any{
					"observedGeneration": int64(1),
				},
			}},
			expr:      "jsonpath={.foo}",
			wantReady: false,
		},
		{
			name: "key with matching value",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": "bar",
			}},
			expr:      "jsonpath={.foo}=bar",
			wantReady: true,
		},
		{
			name: "key with mismatched value",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": "baz",
			}},
			expr:      "jsonpath={.foo}=bar",
			wantReady: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsp, err := jsonpath.Parse(tt.expr)
			require.NoError(t, err)
			c, err := NewJSONPath(context.Background(), nil, stdout, tt.uns, jsp)
			require.NoError(t, err)

			actual, err := c.Satisfied()
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantReady, actual)
		})
	}
}
