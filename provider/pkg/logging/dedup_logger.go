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
	statusLogger    *dedupLogger
	permanentLogger *dedupLogger
}

// NewLogger returns an initialized DedupLogger.
func NewLogger(ctx context.Context, host host.HostClient, urn resource.URN) *DedupLogger {
	status := &dedupLogger{
		messages: &TimeOrderedLogSet{},
		ctx:      ctx,
		urn:      urn,
	}
	permanent := &dedupLogger{
		messages: &TimeOrderedLogSet{},
		ctx:      ctx,
		urn:      urn,
	}
	// Host might be nil in testing.
	if host != nil {
		status.logF = host.LogStatus
		permanent.logF = host.Log
	}
	return &DedupLogger{statusLogger: status, permanentLogger: permanent}
}

// Log logs a permanent message. These are shown as status logs and displayed
// at the end of an interactive session.
func (l *DedupLogger) Log(sev diag.Severity, msg string) {
	l.permanentLogger.logMessage(sev, msg)
}

// LogStatus logs an ephemeral message. These are only shown temporally as
// status messages.
func (l *DedupLogger) LogStatus(sev diag.Severity, msg string) {
	l.statusLogger.logMessage(sev, msg)
}

type dedupLogger struct {
	messages *TimeOrderedLogSet
	index    int
	ctx      context.Context
	logF     func(context.Context, diag.Severity, resource.URN, string) error
	urn      resource.URN
	mux      sync.Mutex
}

// logMessage adds a message to the log set and flushes the queue to the host.
func (l *dedupLogger) logMessage(sev diag.Severity, msg string) {
	l.enqueueMessage(sev, msg)
	l.logNewMessages()
}

// enqueueMessage adds a message to the log set but does not log it to the host.
func (l *dedupLogger) enqueueMessage(severity diag.Severity, s string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.messages.Add(logging.Message{S: s, Severity: severity})
}

// getNewMessages returns the list of new messages since last calling GetNewMessages.
func (l *dedupLogger) getNewMessages() []logging.Message {
	l.mux.Lock()
	defer l.mux.Unlock()

	idx := l.index
	l.index = len(l.messages.Messages)
	return l.messages.Messages[idx:]
}

// logNewMessages logs any new messages to the host.
func (l *dedupLogger) logNewMessages() {
	if l.logF != nil {
		for _, msg := range l.getNewMessages() {
			_ = l.logF(l.ctx, msg.Severity, l.urn, msg.S)
		}
	}
}
