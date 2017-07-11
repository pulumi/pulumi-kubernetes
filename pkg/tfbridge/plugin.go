// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"github.com/pulumi/lumi/pkg/diag"
	"github.com/pulumi/lumi/pkg/resource/provider"
)

// Plugin resolves the path to a Terraform plugin, loads it, and returns two connections to it: one is a standard
// plugin client that can be used to manage its lifetime and the other is a typed provider interface.
func Plugin(host *provider.HostClient, provBin string) (*goplugin.Client, terraform.ResourceProvider, error) {
	// Resolve the path to a plugin.
	plugins := discovery.ResolvePluginPaths([]string{provBin})
	if len(plugins) == 0 {
		return nil, nil, errors.Errorf("No Terraform plugin found at path '%v'", provBin)
	}
	// If multiple were returned (e.g., the path wasn't specific enough), we will choose the newest one.
	plug := plugins.Newest()

	// Now fire up the plugin process and connect to it with a client.  We have to go straight to the plugin interface
	// directly so that we can hook the stdout/stderr streams and redirect them to the appropriate Lumi RPC calls.
	client := goplugin.NewClient(&goplugin.ClientConfig{
		Cmd:             exec.Command(plug.Path),
		HandshakeConfig: plugin.Handshake,
		Managed:         true,
		Plugins:         plugin.PluginMap,
		Stderr: &logRedirector{
			writers: map[string]func(string) error{
				tfTracePrefix: func(msg string) error { return host.Log(diag.Debug, msg) },
				tfDebugPrefix: func(msg string) error { return host.Log(diag.Debug, msg) },
				tfInfoPrefix:  func(msg string) error { return host.Log(diag.Debug, msg) },
				tfWarnPrefix:  func(msg string) error { return host.Log(diag.Warning, msg) },
				tfErrorPrefix: func(msg string) error { return host.Log(diag.Error, msg) },
			},
			defaultWriter: func(msg string) error {
				fmt.Fprintf(os.Stderr, "%v\n", msg)
				return nil
			},
		},
	})
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, err
	}
	raw, err := rpcClient.Dispense(plugin.ProviderPluginName)
	if err != nil {
		return nil, nil, err
	}
	return client, raw.(terraform.ResourceProvider), nil
}

// logRedirector creates a new redirection writer that takes as input plugin stderr output, and routes it to the
// correct Lumi stream based on the standard Terraform logging output prefixes.
type logRedirector struct {
	writers       map[string]func(string) error // the writers for certain labels.
	defaultWriter func(string) error            // the writer for default outputs.
	buffer        []byte                        // a buffer that holds up to a line of output.
}

const (
	tfTracePrefix = "[TRACE]"
	tfDebugPrefix = "[DEBUG]"
	tfInfoPrefix  = "[INFO]"
	tfWarnPrefix  = "[WARN]"
	tfErrorPrefix = "[ERROR]"
)

func (lr *logRedirector) Write(p []byte) (n int, err error) {
	written := 0

	// If a line starts with [TRACE], [DEBUG], or [INFO], then we emit to a debug log entry.  If a line starts with
	// [WARN], we emit a warning.  If a line starts with [ERROR], on the other hand, we emit a normal stderr line.
	// All others simply get redirected to stdout as normal output.
	for len(p) > 0 {
		adv, tok, err := bufio.ScanLines(p, false)
		if err != nil {
			return written, err
		}

		// If adv == 0, there was no newline; buffer it all and move on.
		if adv == 0 {
			lr.buffer = append(lr.buffer, p...)
			written += len(p)
			break
		}

		// Otherwise, there was a newline; emit the buffer plus payload to the right place, and keep going if
		// there is more.
		lr.buffer = append(lr.buffer, tok...) // append the buffer.
		s := string(lr.buffer)

		// To do this we need to parse the label if there is one (e.g., [TRACE], et al).
		var label string
		if start := strings.IndexRune(s, '['); start != -1 {
			if end := strings.IndexRune(s[start:], ']'); end != -1 {
				label = s[start : start+end+1]
			}
		}
		if w, has := lr.writers[label]; has {
			if err := w(s); err != nil {
				return written, err
			}
		} else {
			if err := lr.defaultWriter(s); err != nil {
				return written, err
			}
		}

		// Now keep moving on provided there is more left in the buffer.
		lr.buffer = lr.buffer[:0] // clear out the buffer.
		p = p[adv:]               // advance beyond the extracted region.
		written += adv
	}

	return written, nil
}
