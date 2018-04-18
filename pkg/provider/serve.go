package provider

import (
	"os"

	"github.com/pulumi/pulumi/pkg/resource/provider"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
	lumirpc "github.com/pulumi/pulumi/sdk/proto/go"
	"k8s.io/client-go/tools/clientcmd"

	// Load auth plugins. Removing this will likely cause compilation error.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Serve launches the gRPC server for the Pulumi Kubernetes resource provider.
func Serve(providerName, version string) {
	// Start gRPC service.
	err := provider.Main(
		providerName, func(host *provider.HostClient) (lumirpc.ResourceProviderServer, error) {
			// Use client-go to resolve the final configuration values for the client. Typically these
			// values would would reside in the $KUBECONFIG file, but can also be altered in several
			// places, including in env variables, client-go default values, and (if we allowed it) CLI
			// flags.
			loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
			loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
			kubeconfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(
				loadingRules, &clientcmd.ConfigOverrides{}, os.Stdin)

			return kubeProvider(providerName, version, kubeconfig)
		})

	if err != nil {
		cmdutil.ExitError(err.Error())
	}
}
