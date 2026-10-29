// Copyright 2026, Pulumi Corporation.
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

package await

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	kubetesting "k8s.io/client-go/testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/informers"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	fakehost "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host/fake"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
)

func TestDeletionUIDPrecondition(t *testing.T) {
	for _, tt := range []struct {
		name        string
		uid         types.UID
		deleteError error
	}{
		{"recorded UID", "original", nil},
		{"missing UID", "", nil},
		{"already deleted", "original", apierrors.NewNotFound(schema.GroupResource{Resource: "pods"}, "foo")},
		{
			"UID conflict", "original",
			apierrors.NewConflict(schema.GroupResource{Resource: "pods"}, "foo", errors.New("UID precondition failed")),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			obj := validPodUnstructured.DeepCopy()
			obj.SetUID(tt.uid)
			client, _, _, transport := fake.NewSimpleDynamicClient(fake.WithObjects(obj))
			deletes := 0
			transport.PrependReactor("delete", "pods", func(action kubetesting.Action) (bool, runtime.Object, error) {
				deletes++
				options := action.(kubetesting.DeleteAction).GetDeleteOptions()
				if tt.uid == "" {
					require.Nil(t, options.Preconditions)
				} else {
					require.NotNil(t, options.Preconditions)
					require.NotNil(t, options.Preconditions.UID)
					require.Equal(t, tt.uid, *options.Preconditions.UID)
				}
				return tt.deleteError != nil, nil, tt.deleteError
			})
			urn := resource.NewURN(tokens.QName("test"), tokens.PackageName("test"), "", "kubernetes:core/v1:Pod", "foo")
			host := &fakehost.HostClient{}
			err := Deletion(DeleteConfig{
				ProviderConfig: ProviderConfig{
					Context: t.Context(), ClientSet: client, URN: urn, Host: host,
					DedupLogger: logging.NewLogger(t.Context(), host, urn),
					Factories:   informers.NewFactories(t.Context()),
				},
				Inputs: obj.DeepCopy(), Outputs: obj, Name: obj.GetName(), Timeout: 2,
			})
			require.Equal(t, 1, deletes)
			if apierrors.IsConflict(tt.deleteError) {
				require.ErrorIs(t, err, tt.deleteError)
			} else {
				require.NoError(t, err)
			}
			actions := transport.Actions()
			require.NotEmpty(t, actions)
			require.Equal(t, "delete", actions[0].GetVerb())
			if tt.deleteError != nil {
				require.Len(t, actions, 1)
			}
		})
	}
}
