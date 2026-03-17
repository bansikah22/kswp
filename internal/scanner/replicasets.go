package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/bansikah22/kswp/pkg/models"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetReplicaSets(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]appsv1.ReplicaSet, error) {
	replicasets, err := clientset.AppsV1().ReplicaSets(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return replicasets.Items, nil
}

func GetDeployments(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]appsv1.Deployment, error) {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return deployments.Items, nil
}

func IsReplicaSetOld(rs appsv1.ReplicaSet, deployments []appsv1.Deployment) (bool, string) {
	if rs.Spec.Replicas != nil && *rs.Spec.Replicas > 0 {
		return false, fmt.Sprintf("has %d replicas", *rs.Spec.Replicas)
	}
	for _, owner := range rs.OwnerReferences {
		if owner.Kind == "Deployment" {
			for _, deployment := range deployments {
				if deployment.Name == owner.Name {
					return false, fmt.Sprintf("owned by deployment %s", deployment.Name)
				}
			}
		}
	}
	return true, "Not used by any active deployment"
}

func GetOldReplicaSets(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var oldReplicaSets []models.Resource
	replicasets, err := GetReplicaSets(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	deployments, err := GetDeployments(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	for _, rs := range replicasets {
		if ShouldExclude(rs.ObjectMeta) {
			continue
		}
		old, reason := IsReplicaSetOld(rs, deployments)
		if old {
			oldReplicaSets = append(oldReplicaSets, models.Resource{
				Name:      rs.Name,
				Namespace: rs.Namespace,
				Kind:      "ReplicaSet",
				Reason:    reason,
				Age:       time.Since(rs.CreationTimestamp.Time),
			})
		}
	}
	return oldReplicaSets, nil
}
