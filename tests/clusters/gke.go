package clusters

// import (
// 	"context"
// 	"embed"
// 	"time"
// )

// //go:embed gke-stack/*
// var pulumiProgram embed.FS

// type GKE struct {
// }

// func NewGKECluster(stackName string) (Cluster, error) {
// 	ctx := context.WithTimeout(context.Background(), 15*time.Minute) // It normally takes 7-10 minutes to create a cluster.

// 	return GKE{}, nil
// }

// func (c GKE) Name() string {
// 	return "gke-cluster"
// }

// func (c GKE) Connect() error {
// 	return nil
// }

// func (c GKE) Delete() error {
// 	return nil
// }

// func (c GKE) KubeconfigPath() string {
// 	return ""
// }
