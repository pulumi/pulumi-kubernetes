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

package test

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// These are needed for the Java test suite to pass if we provide our custom flags.
	_ = flag.Bool("create-cluster", false, "Create a cluster for the tests, default is false and to use an existing cluster")
	_ = flag.String("cluster-type", "kind", "The type of cluster to create for the tests, default is kind")
	_ = flag.Int("num-clusters", 1, "The number of clusters to create for the tests, default is 1")
	flag.Parse()

	os.Exit(m.Run())
}
