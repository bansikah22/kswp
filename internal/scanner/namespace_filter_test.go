package scanner

import (
	"testing"
)

func TestIsNamespaceExcluded(t *testing.T) {
	tests := []struct {
		name               string
		namespace          string
		excludedNamespaces []string
		expected           bool
	}{
		{
			name:               "empty excluded list",
			namespace:          "default",
			excludedNamespaces: []string{},
			expected:           false,
		},
		{
			name:               "namespace is excluded",
			namespace:          "kube-system",
			excludedNamespaces: []string{"kube-system", "kube-public"},
			expected:           true,
		},
		{
			name:               "namespace is not excluded",
			namespace:          "app-ns",
			excludedNamespaces: []string{"kube-system", "kube-public"},
			expected:           false,
		},
		{
			name:               "case insensitive match",
			namespace:          "KUBE-SYSTEM",
			excludedNamespaces: []string{"kube-system"},
			expected:           true,
		},
		{
			name:               "case insensitive excluded list",
			namespace:          "kube-system",
			excludedNamespaces: []string{"KUBE-SYSTEM"},
			expected:           true,
		},
		{
			name:               "partial match should not exclude",
			namespace:          "kube",
			excludedNamespaces: []string{"kube-system"},
			expected:           false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsNamespaceExcluded(test.namespace, test.excludedNamespaces)
			if result != test.expected {
				t.Errorf("IsNamespaceExcluded(%q, %v) = %v, expected %v",
					test.namespace, test.excludedNamespaces, result, test.expected)
			}
		})
	}
}

func TestFilterNamespaces(t *testing.T) {
	tests := []struct {
		name               string
		namespaces         []string
		excludedNamespaces []string
		expected           []string
	}{
		{
			name:               "no exclusions",
			namespaces:         []string{"default", "app", "system"},
			excludedNamespaces: []string{},
			expected:           []string{"default", "app", "system"},
		},
		{
			name:               "exclude some namespaces",
			namespaces:         []string{"default", "kube-system", "app", "kube-public"},
			excludedNamespaces: []string{"kube-system", "kube-public"},
			expected:           []string{"default", "app"},
		},
		{
			name:               "exclude all namespaces",
			namespaces:         []string{"kube-system", "kube-public"},
			excludedNamespaces: []string{"kube-system", "kube-public"},
			expected:           []string{},
		},
		{
			name:               "exclude nothing",
			namespaces:         []string{"default", "app"},
			excludedNamespaces: []string{"kube-system"},
			expected:           []string{"default", "app"},
		},
		{
			name:               "case insensitive filtering",
			namespaces:         []string{"default", "KUBE-SYSTEM"},
			excludedNamespaces: []string{"kube-system"},
			expected:           []string{"default"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FilterNamespaces(test.namespaces, test.excludedNamespaces)
			if !slicesEqual(result, test.expected) {
				t.Errorf("FilterNamespaces(%v, %v) = %v, expected %v",
					test.namespaces, test.excludedNamespaces, result, test.expected)
			}
		})
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
