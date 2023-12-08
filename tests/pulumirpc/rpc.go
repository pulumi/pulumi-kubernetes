//nolint:copylocks
package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

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

func NewDebugInterceptorLog() (*DebugInterceptorLog, error) {
	f, err := os.CreateTemp("", "pulumi-grpc-debug-")
	if err != nil {
		return nil, fmt.Errorf("failed to create GRPC debug log file: %v", err)
	}
	defer contract.IgnoreClose(f)
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
	if m, ok := value.(pulumirpc.RegisterResourceRequest); ok {
		return protojson.Format(&m), true
	}
	if m, ok := value.(pulumirpc.RegisterResourceResponse); ok {
		return protojson.Format(&m), true
	}
	if m, ok := value.(pulumirpc.Alias); ok {
		return protojson.Format(&m), true
	}
	return "", false
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
