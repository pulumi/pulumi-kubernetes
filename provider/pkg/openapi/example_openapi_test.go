// Copyright 2021, Pulumi Corporation.  All rights reserved.

package openapi

import (
	"fmt"
)

func ExamplePluck_pathFound() {
	obj := map[string]any{
		"a": map[string]any{
			"x": map[string]any{
				"foo": 1,
				"bar": 2,
			},
		},
	}

	raw, ok := Pluck(obj, "a", "x", "bar")
	fmt.Printf("found = %v\n", ok)
	fmt.Printf("a.x.bar = %v\n", raw)

	// Output:
	// found = true
	// a.x.bar = 2
}

func ExamplePluck_pathNotFound() {
	obj := map[string]any{
		"a": map[string]any{
			"x": map[string]any{
				"foo": 1,
				"bar": 2,
			},
		},
	}

	raw, ok := Pluck(obj, "a", "x", "baz")
	fmt.Printf("found = %v\n", ok)
	fmt.Printf("a.x.baz = %v\n", raw)

	// Output:
	// found = false
	// a.x.baz = <nil>
}
