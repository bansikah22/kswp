package scanner

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestShouldExclude(t *testing.T) {
	tests := []struct {
		name        string
		annotations map[string]string
		expected    bool
	}{
		{
			name:        "no annotations",
			annotations: nil,
			expected:    false,
		},
		{
			name:        "exclude annotation set to true",
			annotations: map[string]string{ExcludeAnnotation: "true"},
			expected:    true,
		},
		{
			name:        "exclude annotation set to false",
			annotations: map[string]string{ExcludeAnnotation: "false"},
			expected:    false,
		},
		{
			name:        "exclude annotation with other value",
			annotations: map[string]string{ExcludeAnnotation: "yes"},
			expected:    false,
		},
		{
			name:        "other annotations but no exclude",
			annotations: map[string]string{"app": "nginx", "version": "1.0"},
			expected:    false,
		},
		{
			name:        "exclude annotation and other annotations",
			annotations: map[string]string{ExcludeAnnotation: "true", "app": "nginx"},
			expected:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			meta := metav1.ObjectMeta{
				Annotations: test.annotations,
			}
			result := ShouldExclude(meta)
			if result != test.expected {
				t.Errorf("ShouldExclude(%v) = %v, expected %v", test.annotations, result, test.expected)
			}
		})
	}
}
