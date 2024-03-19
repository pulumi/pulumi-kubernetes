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

//nolint:govet
package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// JSON format for tracking gRPC conversations. Normal methods have
// one entry for each req-resp conversation, streaming methods have
// one entry per each request or response over the stream.
type DebugInterceptorLogEntry struct {
	Method   string          `json:"method"`
	Request  json.RawMessage `json:"request,omitempty"`
	Response json.RawMessage `json:"response,omitempty"`
	Errors   []string        `json:"errors,omitempty"`
	Metadata interface{}     `json:"metadata,omitempty"`
}

// DebugInterceptorLog encapsulates the gRPC debug log file produced by the PULUMI_DEBUG_GRPC environment variable.
type DebugInterceptorLog struct {
	LogPath string
}

func NewDebugInterceptorLog(t *testing.T) (*DebugInterceptorLog, error) {
	f, err := os.CreateTemp("", "pulumi-grpc-debug-")
	if err != nil {
		return nil, fmt.Errorf("failed to create GRPC debug log file: %v", err)
	}
	defer contract.IgnoreClose(f)
	path, _ := filepath.Abs(f.Name())
	t.Logf("GRPC debug log file: %s", path)
	t.Cleanup(func() {
		if !t.Failed() {
			os.Remove(path)
		}
	})
	return &DebugInterceptorLog{LogPath: f.Name()}, nil
}

func (d *DebugInterceptorLog) Env() string {
	return fmt.Sprintf("PULUMI_DEBUG_GRPC=%s", d.LogPath)
}

func (d *DebugInterceptorLog) Close() error {
	return os.Remove(d.LogPath)
}

func (d *DebugInterceptorLog) Reset() error {
	return os.Truncate(d.LogPath, 0)
}

func (d *DebugInterceptorLog) ReadAll() (DebugInterceptorLogEntryList, error) {
	return ReadDebugInterceptorLogFile(d.LogPath)
}

type DebugInterceptorLogEntryList []DebugInterceptorLogEntry

// ReadDebugInterceptorLogFile parses the gRPC log file produced by the PULUMI_DEBUG_GRPC environment variable.
func ReadDebugInterceptorLogFile(name string) ([]DebugInterceptorLogEntry, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open GRPC debug log file %s: %v", name, err)
	}
	defer f.Close()

	var entries []DebugInterceptorLogEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		entry := &DebugInterceptorLogEntry{}
		err := json.Unmarshal([]byte(text), &entry)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GRPC debug log line (%s): %v", text, err)
		}
		entries = append(entries, *entry)
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("failed to read GRPC debug log file %s: %v", name, err)
	}
	return entries, nil
}

func FormatDebugInterceptorLog(value interface{}) (string, bool) {
	if m, ok := value.(DebugInterceptorLogEntry); ok {
		json, err := json.Marshal(m)
		if err != nil {
			return "", false
		}
		return string(json), true
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Struct:
		// obtain a pointer to the struct in order to call a method on the pointer receiver
		vp := reflect.New(val.Type())
		vp.Elem().Set(val)
		value = vp.Interface()
	}
	if m, ok := value.(proto.Message); ok {
		return protojson.Format(m), true
	}

	return "", false
}

// ListRegisterResource lists the RegisterResource calls in the log.
func (l DebugInterceptorLogEntryList) ListRegisterResource() RegisterResourceList {
	var results []RegisterResource
	for _, entry := range l {
		if r, ok := ParseRegisterResource(entry); ok {
			results = append(results, *r)
		}
	}
	return results
}

// RegisterResource is a decoded "/pulumirpc.ResourceMonitor/RegisterResource" RPC call.
type RegisterResource struct {
	Request  pulumirpc.RegisterResourceRequest
	Response pulumirpc.RegisterResourceResponse
	Errors   []string    `json:"errors,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

type RegisterResourceList []RegisterResource

// Named returns RegisterResource entries matching the given resource type, name, and parent.
func (l RegisterResourceList) Named(parent resource.URN, typ tokens.Type, name tokens.QName) RegisterResourceList {
	var results []RegisterResource
	for _, v := range l {
		if resource.URN(v.Request.Parent) == parent && tokens.Type(v.Request.Type) == typ && tokens.QName(v.Request.Name) == name {
			results = append(results, v)
		}
	}
	return results
}

// ParseRegisterResource parses a log entry as a RegisterResource call.
func ParseRegisterResource(entry DebugInterceptorLogEntry) (*RegisterResource, bool) {
	if entry.Method != "/pulumirpc.ResourceMonitor/RegisterResource" {
		return nil, false
	}
	request := &pulumirpc.RegisterResourceRequest{}
	err := protojson.Unmarshal(entry.Request, request)
	if err != nil {
		return nil, false
	}
	response := &pulumirpc.RegisterResourceResponse{}
	err = protojson.Unmarshal(entry.Response, response)
	if err != nil {
		return nil, false
	}
	result := RegisterResource{
		Request:  *request,
		Response: *response,
		Errors:   entry.Errors,
		Metadata: entry.Metadata,
	}
	return &result, true
}

// Invokes lists the Invoke calls in the log.
func (l DebugInterceptorLogEntryList) Invokes() InvokeList {
	var results []Invoke
	for _, entry := range l {
		if r, ok := ParseInvoke(entry); ok {
			results = append(results, *r)
		}
	}
	return results
}

// Invoke is a decoded "/pulumirpc.ResourceMonitor/Invoke" RPC call.
type Invoke struct {
	Request  pulumirpc.ResourceInvokeRequest
	Response pulumirpc.InvokeResponse
	Errors   []string    `json:"errors,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

type InvokeList []Invoke

// Lookup returns Invoke entries matching the given token (operation).
func (l InvokeList) Tok(tok string) InvokeList {
	var results []Invoke
	for _, v := range l {
		if v.Request.Tok == tok {
			results = append(results, v)
		}
	}
	return results
}

// Lookup returns Invoke entries matching the given provider and token (operation).
func (l InvokeList) ByProvider(providerUrn resource.URN) InvokeList {
	var results []Invoke
	for _, v := range l {
		if resource.URN(v.Request.Provider) == providerUrn {
			results = append(results, v)
		}
	}
	return results
}

// ParseInvoke parses a log entry as a Invoke call.
func ParseInvoke(entry DebugInterceptorLogEntry) (*Invoke, bool) {
	if entry.Method != "/pulumirpc.ResourceMonitor/Invoke" {
		return nil, false
	}
	request := &pulumirpc.ResourceInvokeRequest{}
	err := protojson.Unmarshal(entry.Request, request)
	if err != nil {
		return nil, false
	}
	response := &pulumirpc.InvokeResponse{}
	err = protojson.Unmarshal(entry.Response, response)
	if err != nil {
		return nil, false
	}
	result := Invoke{
		Request:  *request,
		Response: *response,
		Errors:   entry.Errors,
		Metadata: entry.Metadata,
	}
	return &result, true
}
