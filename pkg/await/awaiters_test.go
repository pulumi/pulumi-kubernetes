package await

import (
	"reflect"
	"strings"
	"testing"
)

func TestPluck(t *testing.T) {
	obj2 := map[string]interface{}{
		"string2": "value2",
	}
	obj1 := map[string]interface{}{
		"obj2": obj2,
	}
	apiObject := map[string]interface{}{
		"string1": "value1",
		"obj1":    obj1,
	}

	var tests = []struct {
		path           []string
		expected       interface{}
		expectedType   reflect.Type
		expectedExists bool
		expectedError  bool
	}{
		{
			path:           []string{"string1"},
			expectedType:   reflect.TypeOf(""),
			expected:       "value1",
			expectedExists: true,
		},
		{
			path:          []string{"string1"},
			expectedType:  reflect.TypeOf(0),
			expectedError: true,
		},
		{
			path:          []string{"string1", "invalid"},
			expectedError: true,
		},
		{
			path:           []string{"string2"},
			expectedExists: false,
		},
		{
			path:           []string{"obj1"},
			expectedType:   reflect.TypeOf(map[string]interface{}{}),
			expected:       obj1,
			expectedExists: true,
		},
		{
			path:          []string{"obj1"},
			expectedType:  reflect.TypeOf(0),
			expectedError: true,
		},
		{
			path:           []string{"obj1", "obj2"},
			expectedType:   reflect.TypeOf(map[string]interface{}{}),
			expected:       obj2,
			expectedExists: true,
		},
		{
			path:          []string{"obj1", "obj2"},
			expectedType:  reflect.TypeOf(0),
			expectedError: true,
		},
		{
			path:           []string{"obj1", "string2"},
			expectedExists: false,
		},
		{
			path:           []string{"obj1", "obj2", "string2"},
			expectedType:   reflect.TypeOf(""),
			expected:       "value2",
			expectedExists: true,
		},
		{
			path:          []string{"obj1", "obj2", "string2"},
			expectedType:  reflect.TypeOf(map[string]string{}),
			expectedError: true,
		},
		{
			path:          []string{"obj1", "obj2", "string2", "invalid"},
			expectedError: true,
		},
	}

	for _, test := range tests {
		result, exists, err := pluckT(apiObject, test.expectedType, test.path...)
		pathStr := strings.Join(test.path, ".")

		if test.expectedExists == false || test.expectedError == true {
			if test.expectedExists == false && exists != false {
				t.Errorf("Path '%s' not expected to exist", pathStr)
			}

			if test.expectedError == true && err == nil {
				t.Errorf("Path '%s' expected to return error", pathStr)
			}

			// Successful test.
			continue
		}

		if err != nil {
			t.Errorf("Path '%s' expected to exist, but got error: %v", pathStr, err)
		}

		if !exists {
			t.Errorf("Path '%s' expected to exist, but does not", pathStr)
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Path '%s' expected value '%v' but got '%v'", pathStr, test.expected, result)
		}
	}
}
