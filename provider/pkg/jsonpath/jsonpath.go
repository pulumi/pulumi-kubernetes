package jsonpath

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/theory/jsonpath"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Parsed is a parsed JSONPath expression with an optional value.
type Parsed struct {
	expr *jsonpath.Path

	Path  string // The user's JSONPath expression, for display purposes.
	Value string // An optional value to test for equality against.
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
	results := slices.Collect(i.expr.Select(uns.Object).All())
	if len(results) == 0 {
		return MatchResult{Message: "Missing " + i.Path}, nil
	}
	value := results[0]
	switch v := value.(type) {
	case []any, map[string]any:
		if i.Value == "" {
			// We don't care about complex types if we're matching anything.
			return MatchResult{Matched: true}, nil
		}
		return MatchResult{}, fmt.Errorf("%q has a non-primitive value (%v)", i.Path, v)
	}
	found := fmt.Sprint(value)
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

var _bracketExpr = regexp.MustCompile(`\{\s*(.*?)\s*\}`)

// Parse parses a single JSONPath expression. Only the strict syntax of
// `kubectl get -o jsonpath={...}` is accepted, because the "relaxed" syntax
// used by `wait` is somewhat buggy.
func Parse(raw string) (*Parsed, error) {
	if raw == "" {
		return nil, fmt.Errorf("expected a non-empty JSONPath expression")
	}
	if !strings.HasPrefix(raw, "jsonpath=") {
		return nil, fmt.Errorf("JSONPath expression must begin with a 'jsonpath=' prefix")
	}
	raw = strings.TrimPrefix(raw, "jsonpath=")

	// Split only on "=" and preserve "==".
	placeholder := "�"
	normalized := strings.Replace(raw, "==", placeholder, -1)
	parts := strings.Split(normalized, "=")

	var value string
	path := strings.Replace(parts[0], placeholder, "==", -1)
	if len(parts) > 2 {
		return nil, fmt.Errorf("format should be {.path}=value or {.path}, got %q", raw)
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

	matches := _bracketExpr.FindStringSubmatch(path)
	if len(matches) != 2 {
		return nil, fmt.Errorf("%s should be wrapped in brackets { ... }", path)
	}
	pathWithoutBrackets := matches[1]
	if !strings.HasPrefix("$", pathWithoutBrackets) {
		pathWithoutBrackets = "$" + pathWithoutBrackets
	}

	parser := jsonpath.NewParser()
	expr, err := parser.Parse(pathWithoutBrackets)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %w", path, err)
	}

	return &Parsed{expr: expr, Path: path, Value: value}, nil
}
