package scanner

import (
	"context"
	"fmt"

	"github.com/bansikah22/kswp/pkg/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetUnusedPersistentVolumeClaims(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var resources []models.Resource

	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list pvcs: %w", err)
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	for _, pvc := range pvcs.Items {
		if pvc.Status.Phase == v1.ClaimPending || pvc.Status.Phase == v1.ClaimLost {
			resources = append(resources, models.Resource{
				Name:      pvc.Name,
				Namespace: pvc.Namespace,
				Kind:      "PersistentVolumeClaim",
				Reason:    fmt.Sprintf("PVC is in a %s state", pvc.Status.Phase),
			})
			continue
		}

		if pvc.Status.Phase == v1.ClaimBound {
			isUsed := false
			for _, pod := range pods.Items {
				for _, volume := range pod.Spec.Volumes {
					if volume.PersistentVolumeClaim != nil && volume.PersistentVolumeClaim.ClaimName == pvc.Name {
						isUsed = true
						break
					}
				}
				if isUsed {
					break
				}
			}

			if !isUsed {
				resources = append(resources, models.Resource{
					Name:      pvc.Name,
					Namespace: pvc.Namespace,
					Kind:      "PersistentVolumeClaim",
					Reason:    "PVC is bound but not used by any pod",
				})
			}
		}
	}

	return resources, nil
}
