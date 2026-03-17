package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/bansikah22/kswp/pkg/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPods(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]v1.Pod, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func GetFailedPods(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var failedPods []models.Resource
	pods, err := GetPods(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	for _, pod := range pods {
		if ShouldExclude(pod.ObjectMeta) {
			continue
		}
		if pod.Status.Phase == v1.PodFailed {
			failedPods = append(failedPods, models.Resource{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				Kind:      "Pod",
				Reason:    "Pod has failed",
				Age:       time.Since(pod.CreationTimestamp.Time),
			})
		}
	}
	return failedPods, nil
}

func GetCompletedPods(clientset kubernetes.Interface, threshold time.Duration, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var completedPods []models.Resource
	pods, err := GetPods(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	for _, pod := range pods {
		if ShouldExclude(pod.ObjectMeta) {
			continue
		}
		if pod.Status.Phase == v1.PodSucceeded {
			if pod.Status.StartTime != nil {
				if pod.Status.StartTime.Time.Before(time.Now().Add(-threshold)) {
					completedPods = append(completedPods, models.Resource{
						Name:      pod.Name,
						Namespace: pod.Namespace,
						Kind:      "Pod",
						Reason:    fmt.Sprintf("Completed more than %s ago", threshold.String()),
						Age:       time.Since(pod.CreationTimestamp.Time),
					})
				}
			}
		}
	}
	return completedPods, nil
}
