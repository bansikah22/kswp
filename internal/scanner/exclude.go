package scanner

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ExcludeAnnotation = "kswp.io/exclude"
)

// ShouldExclude returns true if a resource is marked to be excluded from scanning.
// Resources are excluded if they have the kswp.io/exclude annotation set to "true".
func ShouldExclude(meta metav1.ObjectMeta) bool {
	if meta.Annotations == nil {
		return false
	}

	value, ok := meta.Annotations[ExcludeAnnotation]
	return ok && value == "true"
}
