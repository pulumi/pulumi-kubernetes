// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLogDirector ensures that logging redirects to the right place.
func TestLogRedirector(t *testing.T) {
	lines := []string{
		"no prefix #1\n",
		"[TRACE] trace line #1\n",
		"[TRACE] trace line #2\n",
		"no prefix #2\n",
		"[DEBUG] debug line #1\n",
		"[DEBUG] debug line #2\n",
		"[INFO] info line #1\n",
		"no prefix #3\n",
		"[INFO] info line #2\n",
		"[WARN] warning line #1\n",
		"[WARN] warning line #2\n",
		"[ERROR] error line #1\n",
		"[ERROR] error line #2\n",
		"no prefix #4\n",
		"[TRACE] trace line #3\n",
		"[DEBUG] debug line #3\n",
		"[INFO] info line #3\n",
		"[WARN] warning line #3\n",
		"[ERROR] error line #3\n",
		"no prefix #5\n",
	}

	var traces []string
	var debugs []string
	var infos []string
	var warnings []string
	var errors []string

	ld := &LogRedirector{
		enabled: true,
		writers: map[string]func(string) error{
			tfTracePrefix: func(msg string) error {
				traces = append(traces, msg)
				return nil
			},
			tfDebugPrefix: func(msg string) error {
				debugs = append(debugs, msg)
				return nil
			},
			tfInfoPrefix: func(msg string) error {
				infos = append(infos, msg)
				return nil
			},
			tfWarnPrefix: func(msg string) error {
				warnings = append(warnings, msg)
				return nil
			},
			tfErrorPrefix: func(msg string) error {
				errors = append(errors, msg)
				return nil
			},
		},
	}

	// For each line, spit 16 byte increments into the redirector.
	for _, line := range lines {
		for len(line) > 0 {
			sz := 16
			if sz > len(line) {
				sz = len(line)
			}
			n, err := ld.Write([]byte(line[:sz]))
			assert.Nil(t, err)
			assert.Equal(t, n, sz)
			line = line[sz:]
		}
	}

	assert.Equal(t, 3, len(traces))
	assert.Equal(t, 3+5, len(debugs)) // debugs get defaults
	assert.Equal(t, 3, len(infos))
	assert.Equal(t, 3, len(warnings))
	assert.Equal(t, 3, len(errors))
}
