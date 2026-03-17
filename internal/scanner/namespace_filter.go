package scanner

import (
	"strings"
)

const (
	KubeSystemNamespace    = "kube-system"
	KubePublicNamespace    = "kube-public"
	KubeNodeLeaseNamespace = "kube-node-lease"
	KubeApiserverNamespace = "kube-apiserver"
)

var DefaultExcludedNamespaces = []string{
	KubeSystemNamespace,
	KubePublicNamespace,
	KubeNodeLeaseNamespace,
}

func IsNamespaceExcluded(namespace string, excludedNamespaces []string) bool {
	for _, excluded := range excludedNamespaces {
		if strings.EqualFold(namespace, excluded) {
			return true
		}
	}
	return false
}

func FilterNamespaces(namespaces []string, excludedNamespaces []string) []string {
	if len(excludedNamespaces) == 0 {
		return namespaces
	}

	var filtered []string
	for _, ns := range namespaces {
		if !IsNamespaceExcluded(ns, excludedNamespaces) {
			filtered = append(filtered, ns)
		}
	}
	return filtered
}
