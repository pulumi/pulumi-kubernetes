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

package clusters

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	container "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type GKECluster struct {
	name            string
	previousContext string
	client          *container.ClusterManagerClient
	kubeconfig      api.Config
}

var _ Cluster = GKECluster{}

func NewGKECluster(namePrefix string, kubeconfig api.Config) (Cluster, error) {
	// It normally takes 7-10 minutes to create a cluster.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// Login to GCP.

	clusterName := namePrefix + randString()
	projectID := "pulumi-k8s-provider"
	location := "us-west1-a"

	serviceAccount := "pulumi-ci@pulumi-k8s-provider.iam.gserviceaccount.com"
	audience := fmt.Sprintf("https://iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/pulumi-ci/pulumi-ci", 637339343727) // Project number?

	// Default to local credentials.
	auth, err := google.DefaultTokenSource(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting default token source: %w", err)
	}
	if os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN") != "" {
		// If we're running as a GitHub action, grab an ID token so we can use
		// that as our token source.
		token, err := githubactions.GetIDToken(ctx, audience)
		if err != nil {
			return nil, fmt.Errorf("getting ID token: %w", err)
		}
		auth = oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
			TokenType:   "urn:ietf:params:oauth:token-type:jwt",
		})
	}

	token, err := auth.Token()
	if err != nil {
		return nil, fmt.Errorf("getting token: %w", err)
	}

	ts, err := impersonate.CredentialsTokenSource(ctx,
		impersonate.CredentialsConfig{
			Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
			TargetPrincipal: serviceAccount,
			Subject:         token.AccessToken,
		},
		option.WithTokenSource(auth),
		option.WithAudiences(audience),
	)
	if err != nil {
		return nil, fmt.Errorf("getting impersonate source: %w", err)
	}

	// Create the GKE client
	client, err := container.NewClusterManagerClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		log.Fatalf("creating GKE client: %v", err)
	}
	defer client.Close()

	// Define the cluster
	cluster := &containerpb.Cluster{
		Name:             clusterName,
		InitialNodeCount: 1,
		NodeConfig: &containerpb.NodeConfig{
			MachineType: "n1-standard-2",
		},
		Location: location,
	}

	// Create the cluster
	createClusterReq := &containerpb.CreateClusterRequest{
		Parent:  fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Cluster: cluster,
	}

	fmt.Println("creating GKE cluster...")
	_, err = client.CreateCluster(ctx, createClusterReq)
	if err != nil {
		return nil, fmt.Errorf("creating GKE cluster: %w", err)
	}

	fmt.Println("waiting for cluster to become ready..")

	for {
		req := &containerpb.GetClusterRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, location, clusterName),
		}

		cluster, err = client.GetCluster(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to get cluster details: %w", err)
		}

		fmt.Println("cluster is", cluster.Status.String())

		if cluster.Status == containerpb.Cluster_RUNNING {
			break
		}

		select {
		case <-time.After(10 * time.Second):
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Add the cluster to kubeconfig.

	token, err = ts.Token()
	if err != nil {
		return nil, fmt.Errorf("getting token for kubeconfig: %w", err)
	}

	if kubeconfig.Clusters == nil {
		kubeconfig.Clusters = map[string]*api.Cluster{}
	}
	kubeconfig.Clusters[clusterName] = &api.Cluster{
		Server:                   fmt.Sprintf("https://%s", cluster.Endpoint),
		CertificateAuthorityData: []byte(cluster.MasterAuth.ClusterCaCertificate),
	}
	if kubeconfig.Contexts == nil {
		kubeconfig.Contexts = map[string]*api.Context{}
	}
	kubeconfig.Contexts[clusterName] = &api.Context{
		Cluster:  clusterName,
		AuthInfo: clusterName,
	}
	if kubeconfig.AuthInfos == nil {
		kubeconfig.AuthInfos = map[string]*api.AuthInfo{}
	}
	kubeconfig.AuthInfos[clusterName] = &api.AuthInfo{
		Token: token.AccessToken, // Probably won't work...
		/*
			Exec: &api.ExecConfig{
				Command: "gke-gcloud-auth-plugin", // Not installed...
			},
		*/
	}

	previousContext := kubeconfig.CurrentContext
	kubeconfig.CurrentContext = clusterName

	if err := clientcmd.WriteToFile(kubeconfig, "~/.kube/config"); err != nil {
		return nil, fmt.Errorf("writing kubeconfig: %w", err)
	}

	return GKECluster{name: clusterName, previousContext: previousContext, client: client}, nil
}

func (c GKECluster) Name() string {
	return c.name
}

func (c GKECluster) Delete() error {
	_, err := c.client.DeleteCluster(context.Background(), &containerpb.DeleteClusterRequest{
		Name: c.name,
	})
	if err != nil {
		return fmt.Errorf("deleting cluster: %w", err)
	}

	// Remove the cluster's information and restore previous context.
	c.kubeconfig.CurrentContext = c.previousContext
	delete(c.kubeconfig.Clusters, c.name)
	delete(c.kubeconfig.AuthInfos, c.name)
	delete(c.kubeconfig.Contexts, c.name)

	if err := clientcmd.WriteToFile(c.kubeconfig, "~/.kube/config"); err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	return nil
}
