package scanner

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetAllNamespaces(clientset kubernetes.Interface) ([]string, error) {
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var names []string
	for _, ns := range namespaceList.Items {
		names = append(names, ns.Name)
	}
	return names, nil
}

func GetNamespacesToScan(clientset kubernetes.Interface, targetNamespace string, excludedNamespaces []string) ([]string, error) {
	if targetNamespace != "" {
		if IsNamespaceExcluded(targetNamespace, excludedNamespaces) {
			return []string{}, nil
		}
		return []string{targetNamespace}, nil
	}

	allNamespaces, err := GetAllNamespaces(clientset)
	if err != nil {
		return nil, err
	}

	return FilterNamespaces(allNamespaces, excludedNamespaces), nil
}
