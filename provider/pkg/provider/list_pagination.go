// Copyright 2016-2026, Pulumi Corporation.
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

package provider

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// continuationState carries provider-side pagination state across paginated List calls.
type continuationState struct {
	K8sContinue string `json:"k8sContinue,omitempty"` // K8s's own pagination cursor
	Remaining   *int64 `json:"remaining,omitempty"`   // items still allowed under the request's limit; nil = no cap, *0 = cap exhausted
}

// isZero reports whether the pagination state carries nothing worth emitting.
func (c continuationState) isZero() bool {
	if c.K8sContinue == "" {
		return true
	}
	if c.Remaining != nil && *c.Remaining == 0 {
		return true
	}
	return false
}

// effectiveLimit returns the value to set on metav1.ListOptions.Limit for a single
// paginated K8s call, given the request's page_size and the remaining items allowed
// under the session cap. A return of 0 means "no cap" — K8s returns all matching items.
//
// Programming errors panic: pageSize must be non-negative (caller validates the
// request), and *remaining must not be zero or negative (caller must stop before
// invoking this when the cap is exhausted).
func effectiveLimit(pageSize int64, remaining *int64) int64 {
	if pageSize < 0 {
		panic("effectiveLimit: pageSize is negative — caller should have validated the request")
	}
	if remaining != nil && *remaining <= 0 {
		panic("effectiveLimit: remaining points to a non-positive value — caller should have stopped before invoking this")
	}
	if remaining == nil {
		return pageSize
	}
	if pageSize == 0 {
		return *remaining
	}
	return min(pageSize, *remaining)
}

// encodeContinuation serializes the pagination state to an opaque token suitable for
// emission as a ListResponse continuation_token. Zero state encodes to an empty string
// so callers can use the empty value to mean "do not emit a continuation message."
func encodeContinuation(c continuationState) (string, error) {
	if c.isZero() {
		return "", nil
	}
	b, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("encode continuation: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// decodeContinuation parses a token produced by encodeContinuation back into pagination state.
// An empty token yields zero state with no error (first call in a session). A malformed
// token returns an error so the caller can surface InvalidArgument rather than silently
// recovering, since stale or hand-rolled tokens indicate protocol misuse.
func decodeContinuation(token string) (continuationState, error) {
	if token == "" {
		return continuationState{}, nil
	}
	b, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return continuationState{}, fmt.Errorf("decode continuation: %w", err)
	}
	var decoded continuationState
	if err := json.Unmarshal(b, &decoded); err != nil {
		return continuationState{}, fmt.Errorf("decode continuation: %w", err)
	}
	return decoded, nil
}
