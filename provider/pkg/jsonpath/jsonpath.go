package jsonpath

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"
)

// Parsed is a parsed JSONPath expression with an optional value.
type Parsed struct {
	Path   string // The JSONPath expression.
	Value  string // An optional value to test for equality against.
	parser *jsonpath.JSONPath
}

func (i *Parsed) String() string {
	s := i.Path
	if i.Value != "" {
		s += "=" + i.Value
	}
	return s
}

// Matches returns true if the JSONPath matches against the given object. If
// the JSONPath didn't include a value, then this will return true if the path
// exists. Otherwise, the path must exist and hold a value equal to the
// Instance's expected value.
func (i *Parsed) Matches(uns *unstructured.Unstructured) (MatchResult, error) {
	results, err := i.parser.FindResults(uns.Object)
	if err != nil {
		return MatchResult{}, fmt.Errorf("matching JSONPath: %w", err)
	}
	if len(results) == 0 || len(results[0]) == 0 {
		return MatchResult{Message: "Missing " + i.Path}, nil
	}
	value := results[0][0]
	switch value.Interface().(type) {
	case []any, map[string]any:
		return MatchResult{}, fmt.Errorf("%q has a non-primitive value (%v)", i.Path, value.String())
	}
	found := fmt.Sprint(value.Interface())
	if i.Value == "" {
		return MatchResult{Matched: true, Found: found}, nil
	}
	return MatchResult{Matched: i.Value == found, Found: found}, nil
}

// MatchResult contains information about a JSONPath match.
type MatchResult struct {
	Matched bool
	Found   string
	Message string
}

// Parse parses a single JSONPath expression. Only the strict syntax of
// `kubectl get -o jsonpath={...}` is accepted, because the "relaxed" syntax
// used by `wait` is somewhat buggy.
func Parse(expr string) (*Parsed, error) {
	if expr == "" {
		return nil, fmt.Errorf("expected a non-empty JSONPath expression")
	}
	if !strings.HasPrefix(expr, "jsonpath=") {
		return nil, fmt.Errorf("JSONPath expression must begin with a 'jsonpath=' prefix")
	}

	expr = strings.TrimPrefix(expr, "jsonpath=")

	// Split only on "=" and preserve "==".
	var path, value string
	placeholder := "�"
	expr = strings.Replace(expr, "==", placeholder, -1)
	parts := strings.Split(expr, "=")
	path = strings.Replace(parts[0], placeholder, "==", -1)
	if len(parts) > 2 {
		return nil, fmt.Errorf("format should be {.path}=value or {.path}, got %q", expr)
	}
	if len(parts) == 2 {
		value = strings.Replace(parts[1], placeholder, "==", -1)
		if value == "" {
			return nil, fmt.Errorf("%s= requires a value", path)
		}
	}

	if strings.HasPrefix(path, "'") && strings.HasSuffix(path, "'") {
		return nil, fmt.Errorf("%s should omit shell quotes", path)
	}

	parser := jsonpath.New("pulumi").AllowMissingKeys(true)
	if err := parser.Parse(path); err != nil {
		return nil, fmt.Errorf("parsing %q: %w", path, err)
	}

	return &Parsed{Path: path, Value: value, parser: parser}, nil
}
