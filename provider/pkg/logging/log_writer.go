// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"bytes"
	"io"
	"sync"
)

// logWriter is an io.Writer that writes to a logging function, buffering as necessary.

type logF func(format string, args ...interface{})

type LogWriter struct {
	l      logF
	prefix string

	// Holds buffered text for the next write or flush
	// if we haven't yet seen a newline.
	buff bytes.Buffer
	mu   sync.Mutex // guards buff
}

var _ io.Writer = (*LogWriter)(nil)

// NewLogWriter builds and returns an io.Writer that
// writes messages to the given logging function.
// It ensures that each line is logged separately.
//
// Any trailing buffered text that does not end with a newline
// is flushed when the writer flushes.
//
// The returned writer is safe for concurrent use.
func NewLogWriter(l logF) *LogWriter {
	return NewLogWriterPrefixed(l, "")
}

// NewLogWriterPrefixed is a variant of LogWriter
// that prepends the given prefix to each line.
func NewLogWriterPrefixed(l logF, prefix string) *LogWriter {
	w := LogWriter{l: l, prefix: prefix}
	return &w
}

func (w *LogWriter) Write(bs []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// log adds a newline so we should not write bs as-is.
	// Instead, we'll call log one line at a time.
	//
	// To handle the case when Write is called with a partial line,
	// we use a buffer.
	total := len(bs)
	for len(bs) > 0 {
		idx := bytes.IndexByte(bs, '\n')
		if idx < 0 {
			// No newline. Buffer it for later.
			w.buff.Write(bs)
			break
		}

		var line []byte
		line, bs = bs[:idx], bs[idx+1:]

		if w.buff.Len() == 0 {
			// Nothing buffered from a prior partial write.
			// This is the majority case.
			w.l("%s%s", w.prefix, line)
			continue
		}

		// There's a prior partial write. Join and flush.
		w.buff.Write(line)
		w.l("%s%s", w.prefix, w.buff.String())
		w.buff.Reset()
	}
	return total, nil
}

// flush flushes buffered text, even if it doesn't end with a newline.
func (w *LogWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buff.Len() > 0 {
		w.l("%s%s", w.prefix, w.buff.String())
		w.buff.Reset()
	}
}

func (w *LogWriter) Close() error {
	w.Flush()
	return nil
}
