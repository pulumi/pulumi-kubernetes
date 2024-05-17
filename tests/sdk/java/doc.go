// Copyright 2024, Pulumi Corporation.
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

// Package test contains a number of end-to-end tests written in YAML despite
// being located under sdk/java. This is necessary because our CI jobs don't
// currently execute a test step for YAML. We don't have any Java SDK tests at
// the moment so we might as well give this shard something to do.
//
// See https://github.com/pulumi/ci-mgmt/issues/676#issuecomment-2085892601 for
// a discussion of how to better distribute these tests in CI.
package test
