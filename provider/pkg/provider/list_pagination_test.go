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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestEffectiveLimit(t *testing.T) {
	cases := []struct {
		name      string
		pageSize  int64
		remaining *int64
		want      int64
	}{
		{"no caps anywhere — K8s returns everything", 0, nil, 0},
		{"only page size — no session cap", 25, nil, 25},
		{"only remaining — caller did not set page size", 0, ptr.To(int64(50)), 50},
		{"page smaller than remaining — page caps the call", 25, ptr.To(int64(50)), 25},
		{"remaining smaller than page — final partial page", 25, ptr.To(int64(10)), 10},
		{"equal", 25, ptr.To(int64(25)), 25},
		{"remaining = 1, about to finish", 25, ptr.To(int64(1)), 1},
		{"page = 1, remaining = 100", 1, ptr.To(int64(100)), 1},
		{"both 1", 1, ptr.To(int64(1)), 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := effectiveLimit(tc.pageSize, tc.remaining)
			assert.Equal(t, tc.want, got)
			assert.GreaterOrEqual(t, got, int64(0), "result must never be negative")
			if tc.pageSize > 0 {
				assert.LessOrEqual(t, got, tc.pageSize, "result must not exceed page_size when set")
			}
			if tc.remaining != nil {
				assert.LessOrEqual(t, got, *tc.remaining, "result must not exceed remaining when set")
			}
		})
	}
}

func TestEffectiveLimit_PanicsOnExhaustedRemaining(t *testing.T) {
	assert.Panics(t, func() {
		effectiveLimit(10, ptr.To(int64(0)))
	}, "remaining=*0 (cap exhausted) is a programming error; caller should have stopped")
}

func TestEffectiveLimit_PanicsOnNegativeRemaining(t *testing.T) {
	assert.Panics(t, func() {
		effectiveLimit(10, ptr.To(int64(-1)))
	}, "negative remaining is invalid")
}

func TestEffectiveLimit_PanicsOnNegativePageSize(t *testing.T) {
	assert.Panics(t, func() {
		effectiveLimit(-1, nil)
	}, "negative pageSize is invalid — caller should have rejected the request upstream")
}

func TestIsZero(t *testing.T) {
	cases := []struct {
		name string
		cont continuationState
		want bool
	}{
		{"empty struct", continuationState{}, true},
		{
			"more K8s pages available, no session cap",
			continuationState{K8sContinue: "abc"}, false,
		},
		{
			"more K8s pages, session budget unspent",
			continuationState{K8sContinue: "abc", Remaining: ptr.To(int64(25))}, false,
		},
		{
			"session cap exhausted",
			continuationState{K8sContinue: "abc", Remaining: ptr.To(int64(0))}, true,
		},
		{
			"no more K8s pages, session budget unspent",
			continuationState{Remaining: ptr.To(int64(5))}, true,
		},
		{
			"no more K8s pages and cap exhausted",
			continuationState{Remaining: ptr.To(int64(0))}, true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.cont.isZero())
		})
	}
}

func TestContinuationRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		in   continuationState
	}{
		{"k8s cursor only — no cap", continuationState{K8sContinue: "abc-cursor"}},
		{"k8s cursor + remaining set", continuationState{K8sContinue: "cursor-xyz", Remaining: ptr.To(int64(100))}},
		{
			"k8s cursor with chars outside url-safe base64",
			continuationState{K8sContinue: "abc/def+xyz=", Remaining: ptr.To(int64(5))},
		},
		{"large remaining", continuationState{K8sContinue: "x", Remaining: ptr.To(int64(1_000_000))}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tok, err := encodeContinuation(tc.in)
			require.NoError(t, err)
			require.NotEmpty(t, tok, "non-zero state must produce a non-empty token")
			out, err := decodeContinuation(tok)
			require.NoError(t, err)
			assert.Equal(t, tc.in.K8sContinue, out.K8sContinue)
			if tc.in.Remaining == nil {
				assert.Nil(t, out.Remaining, "nil remaining must round-trip as nil (no cap)")
			} else {
				require.NotNil(t, out.Remaining)
				assert.Equal(t, *tc.in.Remaining, *out.Remaining)
			}
		})
	}
}

func TestEncodeContinuation_ZeroStateIsEmpty(t *testing.T) {
	tok, err := encodeContinuation(continuationState{})
	require.NoError(t, err)
	assert.Empty(t, tok, "zero state must encode to empty string so no continuation message is emitted")
}

func TestEncodeContinuation_CapExhaustedIsEmpty(t *testing.T) {
	tok, err := encodeContinuation(continuationState{K8sContinue: "abc", Remaining: ptr.To(int64(0))})
	require.NoError(t, err)
	assert.Empty(t, tok, "cap-exhausted state must encode to empty string even when K8s has more pages")
}

func TestEncodeContinuation_K8sDoneIsEmpty(t *testing.T) {
	tok, err := encodeContinuation(continuationState{Remaining: ptr.To(int64(50))})
	require.NoError(t, err)
	assert.Empty(t, tok, "no K8s cursor means K8s is done — nothing to resume")
}

func TestDecodeContinuation_EmptyIsZero(t *testing.T) {
	state, err := decodeContinuation("")
	require.NoError(t, err)
	assert.Equal(t, continuationState{}, state, "empty token must decode to zero state — used on first call")
}

func TestDecodeContinuation_RemainingZeroRoundtrips(t *testing.T) {
	// A continuation token carrying Remaining=0 should not normally be emitted, but if
	// callers ever decode such a token (e.g. hand-rolled in a test), the pointer must
	// not be confused with nil.
	encoded := base64.RawURLEncoding.EncodeToString([]byte(`{"k8sContinue":"abc","remaining":0}`))
	state, err := decodeContinuation(encoded)
	require.NoError(t, err)
	require.NotNil(t, state.Remaining)
	assert.Equal(t, int64(0), *state.Remaining)
}

func TestDecodeContinuation_MalformedRejected(t *testing.T) {
	notJSON := base64.RawURLEncoding.EncodeToString([]byte("not-json"))
	wrongShape := base64.RawURLEncoding.EncodeToString([]byte(`["array","not","object"]`))

	cases := []struct {
		name  string
		token string
	}{
		{"not base64", "not-base64-!@#$"},
		{"valid base64, not JSON", notJSON},
		{"valid base64, wrong JSON shape", wrongShape},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := decodeContinuation(tc.token)
			assert.Error(t, err, "malformed token must surface an error so callers can return InvalidArgument")
		})
	}
}
