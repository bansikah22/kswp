package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/bansikah22/kswp/pkg/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ttlAnnotation = "cleaner/ttl"
)

func GetExpiredResources(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var resources []models.Resource

	configMaps, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list configmaps: %w", err)
	}

	for _, cm := range configMaps.Items {
		expired, reason := isExpired(cm.ObjectMeta)
		if expired {
			resources = append(resources, models.Resource{
				Name:      cm.Name,
				Namespace: cm.Namespace,
				Kind:      "ConfigMap",
				Reason:    reason,
			})
		}
	}

	secrets, err := clientset.CoreV1().Secrets(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	for _, s := range secrets.Items {
		expired, reason := isExpired(s.ObjectMeta)
		if expired {
			resources = append(resources, models.Resource{
				Name:      s.Name,
				Namespace: s.Namespace,
				Kind:      "Secret",
				Reason:    reason,
			})
		}
	}

	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	for _, s := range services.Items {
		expired, reason := isExpired(s.ObjectMeta)
		if expired {
			resources = append(resources, models.Resource{
				Name:      s.Name,
				Namespace: s.Namespace,
				Kind:      "Service",
				Reason:    reason,
			})
		}
	}

	replicaSets, err := clientset.AppsV1().ReplicaSets(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list replicasets: %w", err)
	}

	for _, rs := range replicaSets.Items {
		expired, reason := isExpired(rs.ObjectMeta)
		if expired {
			resources = append(resources, models.Resource{
				Name:      rs.Name,
				Namespace: rs.Namespace,
				Kind:      "ReplicaSet",
				Reason:    reason,
			})
		}
	}

	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	for _, j := range jobs.Items {
		expired, reason := isExpired(j.ObjectMeta)
		if expired {
			resources = append(resources, models.Resource{
				Name:      j.Name,
				Namespace: j.Namespace,
				Kind:      "Job",
				Reason:    reason,
			})
		}
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	for _, p := range pods.Items {
		expired, reason := isExpired(p.ObjectMeta)
		if expired {
			resources = append(resources, models.Resource{
				Name:      p.Name,
				Namespace: p.Namespace,
				Kind:      "Pod",
				Reason:    reason,
			})
		}
	}

	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list pvcs: %w", err)
	}

	for _, pvc := range pvcs.Items {
		expired, reason := isExpired(pvc.ObjectMeta)
		if expired {
			resources = append(resources, models.Resource{
				Name:      pvc.Name,
				Namespace: pvc.Namespace,
				Kind:      "PersistentVolumeClaim",
				Reason:    reason,
			})
		}
	}

	return resources, nil
}

func isExpired(meta metav1.ObjectMeta) (bool, string) {
	ttlValue, ok := meta.Annotations[ttlAnnotation]
	if !ok {
		return false, ""
	}

	ttl, err := time.ParseDuration(ttlValue)
	if err != nil {
		return false, ""
	}

	expirationTime := meta.CreationTimestamp.Add(ttl)
	if time.Now().After(expirationTime) {
		return true, fmt.Sprintf("Expired based on TTL: %s", ttlValue)
	}

	return false, ""
}
