package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		wantPath  string
		wantValue string
		wantErr   string
	}{
		{
			name:    "empty expression",
			expr:    "",
			wantErr: "non-empty",
		},
		{
			name:    "missing prefix",
			expr:    "{.foo}",
			wantErr: "jsonpath=",
		},
		{
			name:    "quoted key with value",
			expr:    "jsonpath='{.status.phase}'=Running",
			wantErr: "omit shell quotes",
		},
		{
			name:    "missing value",
			expr:    "jsonpath={.metadata.name}=",
			wantErr: "{.metadata.name}= requires a value",
		},
		{
			name:    "invalid expression with repeated =",
			expr:    "jsonpath={.metadata.name}='test=wrong'",
			wantErr: "format should be {.path}=value or {.path}",
		},
		{
			name:    "complex expressions are not supported",
			expr:    "jsonpath={.status.conditions[?(@.type==\"Failed\"||@.type==\"Complete\")].status}=True",
			wantErr: "unrecognized character",
		},
		{
			name:     "key with any value",
			expr:     "jsonpath={.foo}",
			wantPath: "{.foo}",
		},
		{
			name:      "key with value",
			expr:      "jsonpath={.foo}=bar",
			wantPath:  "{.foo}",
			wantValue: "bar",
		},
		{
			name:      "preserve ==",
			expr:      `jsonpath={.status.containerStatuses[?(@.name=="foobar")].ready}=True`,
			wantPath:  `{.status.containerStatuses[?(@.name=="foobar")].ready}`,
			wantValue: "True",
		},
		{
			name:     "padded brackets",
			expr:     "jsonpath={ .webhooks[].clientConfig.caBundle }",
			wantPath: `{ .webhooks[].clientConfig.caBundle }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance, err := Parse(tt.expr)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantPath, instance.Path)
			assert.Equal(t, tt.wantValue, instance.Value)
		})
	}
}

func TestMatches(t *testing.T) {
	tests := []struct {
		name string
		expr string
		uns  *unstructured.Unstructured
		want MatchResult
		// wantInfo  string
		wantErr string
	}{
		{
			name: "no match",
			expr: "jsonpath={.foo}",
			uns:  &unstructured.Unstructured{Object: map[string]any{}},
			want: MatchResult{Matched: false, Message: "Missing {.foo}"},
		},
		{
			name: "key exists",
			expr: "jsonpath={ .foo }",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": nil,
			}},
			want: MatchResult{Matched: true, Found: "<nil>"},
		},
		{
			name: "object exists",
			expr: "jsonpath={ .foo }",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			}},
			want: MatchResult{Matched: true},
		},
		{
			name: "key exists with non-primitive value",
			expr: "jsonpath={.foo}",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": []string{"boo"},
			}},
			want: MatchResult{Matched: true, Found: "[boo]"},
		},
		{
			name: "value matches",
			expr: "jsonpath={.foo}=bar",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": "bar",
			}},
			want: MatchResult{Matched: true, Found: "bar"},
		},
		{
			name: "value mismatch",
			expr: "jsonpath={.foo}=bar",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": "baz",
			}},
			want: MatchResult{Matched: false, Found: "baz"},
		},
		{
			name: "value match against some array element",
			expr: "jsonpath={.foo[*].bar}=baz",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": []any{
					map[string]any{
						"ignored": "true",
					},
					map[string]any{
						"bar": "baz",
					},
					map[string]any{
						"something else": "true",
					},
				},
			}},
			want: MatchResult{Matched: true, Found: "baz"},
		},
		{
			name: "value match against specific array element",
			expr: "jsonpath={.foo[1].bar}=baz",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": []any{
					map[string]any{
						"bar": "not-baz",
					},
					map[string]any{
						"bar": "baz",
					},
				},
			}},
			want: MatchResult{Matched: true, Found: "baz"},
		},
		{
			name: "value mismatch against specific array element",
			expr: "jsonpath={.foo[0].bar}=baz",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": []any{
					map[string]any{
						"bar": "not-baz",
					},
					map[string]any{
						"bar": "baz",
					},
				},
			}},
			want: MatchResult{Matched: false, Found: "not-baz"},
		},
		{
			name: "value match against non-primitive value",
			expr: "jsonpath={.foo}=bar",
			uns: &unstructured.Unstructured{Object: map[string]any{
				"foo": []any{"bar"},
			}},
			wantErr: "has a non-primitive value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := Parse(tt.expr)
			require.NoError(t, err)

			actual, err := i.Matches(tt.uns)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		given Parsed
		want  string
	}{
		{
			given: Parsed{Path: "{.foo}"},
			want:  "{.foo}",
		},
		{
			given: Parsed{Path: "{.foo}", Value: "1"},
			want:  "{.foo}=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.given.String())
		})
	}
}
