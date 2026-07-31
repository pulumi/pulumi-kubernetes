package jsonpath

import (
	"fmt"
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
		return MatchResult{Message: i.Path + " selected nothing"}, nil
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

	parts := splitOnValue(raw)

	var value string
	path := parts[0]
	if len(parts) > 2 {
		return nil, fmt.Errorf("format should be {.path}=value or {.path}, got %q", raw)
	}
	if len(parts) == 2 {
		value = parts[1]
		if value == "" {
			return nil, fmt.Errorf("%s= requires a value", path)
		}
	}

	if strings.HasPrefix(path, "'") && strings.HasSuffix(path, "'") {
		return nil, fmt.Errorf("%s should omit shell quotes", path)
	}

	// Comparing one path against another isn't supported; without this the value is treated as a literal and never matches.
	if strings.HasPrefix(value, "{") {
		return nil, fmt.Errorf("%s=%s compares against another JSONPath, which is not supported; the value must be a literal", path, value)
	}

	pathWithoutBrackets, ok := unbracket(path)
	if !ok {
		return nil, fmt.Errorf("%s should be wrapped in brackets { ... }", path)
	}
	if !strings.HasPrefix(pathWithoutBrackets, "$") {
		pathWithoutBrackets = "$" + pathWithoutBrackets
	}

	parser := jsonpath.NewParser()
	expr, err := parser.Parse(pathWithoutBrackets)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %w", path, err)
	}

	return &Parsed{expr: expr, Path: path, Value: value}, nil
}

// splitOnValue splits raw into the bracketed JSONPath expression and an optional value, breaking only on "=" found
// outside the brackets. Comparison operators inside the expression ("==", "!=", "<=", ">=") are left intact, while a
// value containing a bare "=" still yields more than two elements so it can be rejected as ambiguous.
func splitOnValue(raw string) []string {
	var parts []string
	var element strings.Builder
	depth := 0
	var quote byte // The active string-literal delimiter, when inside the brackets.

	for i := 0; i < len(raw); i++ {
		c := raw[i]
		switch {
		case quote != 0:
			// Inside a string literal, where only the closing delimiter is structural.
			if c == '\\' && i < len(raw)-1 {
				element.WriteByte(c)
				i++
				c = raw[i]
				break
			}
			if c == quote {
				quote = 0
			}
		case depth > 0 && (c == '\'' || c == '"'):
			quote = c
		case c == '{':
			depth++
		case c == '}':
			if depth > 0 {
				depth--
			}
		case c == '=' && depth == 0:
			parts = append(parts, element.String())
			element.Reset()
			continue
		}
		element.WriteByte(c)
	}

	return append(parts, element.String())
}

// unbracket returns the contents of the first balanced { ... } group in path, ignoring braces inside string literals.
func unbracket(path string) (string, bool) {
	start := strings.IndexByte(path, '{')
	if start < 0 {
		return "", false
	}

	depth := 0
	var quote byte

	for i := start; i < len(path); i++ {
		c := path[i]
		switch {
		case quote != 0:
			if c == '\\' && i < len(path)-1 {
				i++
			} else if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"':
			quote = c
		case c == '{':
			depth++
		case c == '}':
			depth--
			if depth == 0 {
				return strings.TrimSpace(path[start+1 : i]), true
			}
		}
	}

	return "", false
}
