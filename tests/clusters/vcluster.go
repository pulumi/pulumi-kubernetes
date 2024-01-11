package clusters

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type vclusterList struct {
	Name string `json:"name"`
}

func NewVCluster(name string) (Cluster, error) {
	// Get the default kubeconfig path from ENV as this is where vcluster will write the kubeconfig.
	kubeconfigPath := os.Getenv("KUBECONFIG")

	name = normalizeName(name + "-" + randString())

	cmd := exec.Command("vcluster", "create", "--connect=false", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Errorf("failed to create cluster %s: %v\nstdout/stderr:\n%s", name, err, string(out))
	}

	return &VCluster{
		name:           name,
		kubeconfigPath: kubeconfigPath,
	}, nil
}

type VCluster struct {
	name           string
	kubeconfigPath string
}

func (c VCluster) KubeconfigPath() string {
	return c.kubeconfigPath
}

func (c VCluster) Name() string {
	return c.name
}

func (c VCluster) Exists() bool {
	cmd := exec.Command("vcluster", "list", "--output=json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	var clusters []vclusterList
	if err := json.Unmarshal(out, &clusters); err != nil {
		return false
	}

	for _, cluster := range clusters {
		if cluster.Name == c.name {
			return true
		}
	}

	return false
}

func (c VCluster) Connect() error {
	cmd := exec.Command("vcluster", "connect", c.name)
	err := cmd.Start()
	if err != nil {
		return errors.Errorf("failed to connect to cluster %s: %v", c.name, err)
	}

	return nil
}

func (c VCluster) disconnect() error {
	cmd := exec.Command("vcluster", "disconnect")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Errorf("failed to disconnect from cluster %s: %v\nstdout/stderr:\n%s", c.name, err, string(out))
	}

	return nil
}

func (c VCluster) Delete() error {
	err := c.disconnect()
	if err != nil {
		return err
	}

	cmd := exec.Command("vcluster", "delete", c.name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Errorf("failed to delete cluster %s: %v\nstdout/stderr:\n%s", c.name, err, string(out))
	}

	return nil
}
