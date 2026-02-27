// Copyright 2025, Pulumi Corporation.
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

package helm

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// newTLSConfig is a simplification of Helm's internal tlsutil. It's
// responsible for creating a simple tls.Config, given some on-disk key paths.
//
// https://github.com/helm/helm/blob/01adbab466b6133936cac0c56a99274715f7c085/internal/tlsutil/tls.go#L17
func newTLSConfig(certFile, keyFile, caFile string, insecureSkipVerify bool) (*tls.Config, error) {
	config := &tls.Config{InsecureSkipVerify: insecureSkipVerify} //nolint:gosec // Intentional insecureSkipVerify.

	if certFile != "" && keyFile != "" {
		certPEMBlock, err := os.ReadFile(certFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read cert file: %q: %w", certFile, err)
		}

		keyPEMBlock, err := os.ReadFile(keyFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read key file: %q: %w", keyFile, err)
		}

		cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			return nil, fmt.Errorf("unable to load cert from key pair: %w", err)
		}

		config.Certificates = []tls.Certificate{cert}
	}

	if caFile != "" {
		caPEMBlock, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("can't read CA file: %q: %w", caFile, err)
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(caPEMBlock) {
			return nil, fmt.Errorf("failed to append certificates from pem block")
		}

		config.RootCAs = cp
	}

	return config, nil
}
