package kubernetes

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client interface {
	Clientset() kubernetes.Interface
}

type client struct {
	clientset kubernetes.Interface
}

func (c *client) Clientset() kubernetes.Interface {
	return c.clientset
}

func NewClient() (Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		// Fallback to local kubeconfig if in-cluster config fails
		userHomeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return nil, fmt.Errorf("failed to get in-cluster config and user home directory: %v, %v", err, homeErr)
		}
		kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &client{clientset: clientset}, nil
}
