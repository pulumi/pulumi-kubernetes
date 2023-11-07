// Copyright 2016-2023, Pulumi Corporation.
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

package crd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/httputil"
)

// ReadFromLocalOrRemote reads the contents of a file from the local filesystem or from a remote URL.
func ReadFromLocalOrRemote(pathOrURL string, headers map[string]string) (io.ReadCloser, error) {
	if strings.HasPrefix(pathOrURL, "https://") {
		client := &http.Client{Timeout: 30 * time.Second}
		req, err := http.NewRequest("GET", pathOrURL, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create HTTP request for %q: %w", pathOrURL, err)
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}

		resp, err := httputil.DoWithRetry(req, client)
		if err != nil {
			return nil, fmt.Errorf("failed to make HTTP request to %q: %w", pathOrURL, err)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP request to %q failed with status %d", pathOrURL, resp.StatusCode)
		}
		return resp.Body, nil
	}
	file, err := os.Open(pathOrURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", pathOrURL, err)
	}
	return file, nil
}
