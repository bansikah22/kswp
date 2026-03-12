package kubernetes

import (
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
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &client{clientset: clientset}, nil
}
