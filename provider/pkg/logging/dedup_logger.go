// Copyright 2016-2022, Pulumi Corporation.
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
	"context"
	"sync"

	"github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// DedupLogger wraps a time-ordered log set to allow batched logging of unique messages.
// Operations on DedupLogger are safe to use concurrently.
type DedupLogger struct {
	messages *TimeOrderedLogSet
	index    int
	ctx      context.Context
	host     host.HostClient
	urn      resource.URN
	mux      sync.Mutex
}

// NewLogger returns an initialized DedupLogger.
func NewLogger(ctx context.Context, host host.HostClient, urn resource.URN) *DedupLogger {
	return &DedupLogger{
		messages: &TimeOrderedLogSet{},
		ctx:      ctx,
		host:     host,
		urn:      urn,
	}
}

// LogMessage adds a message to the log set and flushes the queue to the host.
func (l *DedupLogger) LogMessage(msg logging.Message) {
	l.EnqueueMessage(msg.Severity, msg.S)
	l.LogNewMessages()
}

// EnqueueMessage adds a message to the log set but does not log it to the host.
func (l *DedupLogger) EnqueueMessage(severity diag.Severity, s string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.messages.Add(logging.Message{S: s, Severity: severity})
}

// GetNewMessages returns the list of new messages since last calling GetNewMessages.
func (l *DedupLogger) GetNewMessages() []logging.Message {
	l.mux.Lock()
	defer l.mux.Unlock()

	idx := l.index
	l.index = len(l.messages.Messages)
	return l.messages.Messages[idx:]
}

// LogNewMessages logs any new messages to the host.
func (l *DedupLogger) LogNewMessages() {
	if l.host != nil {
		for _, msg := range l.GetNewMessages() {
			_ = l.host.LogStatus(l.ctx, msg.Severity, l.urn, msg.S)
		}
	}
}
