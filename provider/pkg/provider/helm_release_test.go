package provider

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_MergeMaps(t *testing.T) {
	m := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"d": []interface{}{
					"1", "2",
				},
			},
		},
	}

	override := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"d": []interface{}{
					"3", "4",
				},
			},
		},
	}

	for _, test := range []struct {
		name     string
		dest     map[string]interface{}
		src      map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "Precedence",
			allowNil: false,
			dest:     m,
			src:      override,
			expected: override, // Expect the override to take precedence
		},
		{
			name:     "Merge maps",
			allowNil: false,
			dest:     m,
			src: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": []interface{}{
							"3", "4",
						},
						"f": true,
					},
				},
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"d": []interface{}{
							"1", "2",
						},
						"c": []interface{}{
							"3", "4",
						},
						"f": true,
					},
				},
			},
		},
		{
			name:     "Dest Has Nil Values- disallow nil",
			allowNil: false,
			dest:     m,
			src: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": interface{}(nil),
						"e": (*interface{})(nil),
					},
				},
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"d": []interface{}{
							"1", "2",
						},
					},
				},
			},
		},
		{
			name:     "Dest Has Nil Values- allow nil",
			allowNil: true,
			dest:     m,
			src: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": interface{}(nil),
						"e": (*interface{})(nil),
					},
				},
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": interface{}(nil),
						"e": (*interface{})(nil),
						"d": []interface{}{
							"1", "2",
						},
					},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			merged, err := mergeMaps(test.dest, test.src, test.allowNil)
			require.NoError(t, err)
			assert.Equal(t, test.expected, merged)
		})
	}
}
