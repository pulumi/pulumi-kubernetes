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

// logF is an abstract logging function that accepts a format string and arguments.
// The function is expected to write a newline after each message.
type logF func(format string, args ...interface{})

// logWriter is an io.Writer that writes to a logging function, buffering as necessary.
type LogWriter struct {
	l      logF
	prefix string

	// Holds buffered text for the next write or flush
	// if we haven't yet seen a newline.
	buff bytes.Buffer
	mu   sync.Mutex // guards buff
}

type Option interface {
	apply(*LogWriter)
}

type prefixOption string

func (p prefixOption) apply(l *LogWriter) {
	l.prefix = string(p)
}

// WithPrefix prepends the given prefix to each line.
func WithPrefix(prefix string) Option {
	return prefixOption(prefix)
}

// NewLogWriter builds and returns an io.Writer that
// writes messages to the given logging function.
// It ensures that each line is logged separately.
//
// Any trailing buffered text that does not end with a newline
// is flushed when the writer flushes.
//
// The returned writer is safe for concurrent use.
func NewLogWriter(l logF, opts ...Option) *LogWriter {
	w := &LogWriter{l: l}
	for _, o := range opts {
		o.apply(w)
	}
	return w
}

var _ io.Writer = (*LogWriter)(nil)

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
