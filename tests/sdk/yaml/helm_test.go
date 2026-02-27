package test

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/stretchr/testify/require"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

func TestHelmWithPrivateOCIRegistry(t *testing.T) {
	// Create a private ECR registry which can accept OCI artifacts.

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		t.Skip(err)
	}

	svc := ecr.New(sess)
	name := strings.ToLower(t.Name()) + fmt.Sprint(rand.Intn(1000)) //nolint:gosec // Doesn't need to be secure.

	params := &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(name),
	}
	resp, err := svc.CreateRepository(params)
	require.NoError(t, err)

	t.Cleanup(func() {
		// Make sure we remove this repo afterwards.
		_, _ = svc.DeleteRepository(&ecr.DeleteRepositoryInput{
			Force:          aws.Bool(true),
			RegistryId:     resp.Repository.RegistryId,
			RepositoryName: resp.Repository.RepositoryName,
		})
	})

	// Grab authToken for the registry.
	authToken, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	require.NoError(t, err)

	b64token := authToken.AuthorizationData[0].AuthorizationToken
	token, err := base64.StdEncoding.DecodeString(*b64token)
	require.NoError(t, err)

	parts := strings.SplitN(string(token), ":", 2)
	username := parts[0]
	password := parts[1]

	// Setup an OCI client for the registry.
	reg, err := remote.NewRegistry(*resp.Repository.RegistryId + ".dkr.ecr.us-west-2.amazonaws.com")
	require.NoError(t, err)

	// Login to the private OCI registry.
	creds := auth.Credential{
		Username: username,
		Password: password,
	}
	reg.Client = &auth.Client{
		Credential: auth.StaticCredential(reg.Reference.Registry, creds),
	}
	require.NoError(t, reg.Ping(t.Context()))

	// Fetch a remote nginx chart into memory.
	ref := "20.0.3"
	memStore := memory.New()
	sourceRepo, err := remote.NewRepository("registry-1.docker.io/bitnamicharts/nginx")
	require.NoError(t, err)

	// Add DockerHub auth token for higher RPS limit
	if dockerPassword := os.Getenv("DOCKER_HUB_PASSWORD"); dockerPassword != "" {
		dockerCreds := auth.Credential{
			Username: "pulumibot",
			Password: dockerPassword,
		}
		sourceRepo.Client = &auth.Client{
			Credential: auth.StaticCredential(sourceRepo.Reference.Registry, dockerCreds),
		}
	}

	_, err = oras.Copy(t.Context(), sourceRepo, ref, memStore, ref, oras.DefaultCopyOptions)
	require.NoError(t, err)

	// Push the chart to our private registry/repo.
	repo, err := reg.Repository(t.Context(), name)
	require.NoError(t, err)
	_, err = oras.Copy(t.Context(), memStore, ref, repo, ref, oras.DefaultCopyOptions)
	require.NoError(t, err)

	// Now run a Pulumi program which fetches our private chart.
	test := pulumitest.NewPulumiTest(t, "testdata/oci",
		opttest.SkipInstall(),
	)

	test.SetConfig(t, "chart", "oci://"+*resp.Repository.RepositoryUri)
	test.SetConfig(t, "version", ref)
	test.SetConfig(t, "username", username)
	test.SetConfig(t, "password", password)

	test.Preview(t)
}
