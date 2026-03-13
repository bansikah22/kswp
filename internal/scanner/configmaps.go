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

func GetConfigMaps(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]v1.ConfigMap, error) {
	configmaps, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return configmaps.Items, nil
}

func IsConfigMapUsed(configMap v1.ConfigMap, pods []v1.Pod) (bool, string) {
	for _, pod := range pods {
		for _, volume := range pod.Spec.Volumes {
			if volume.ConfigMap != nil && volume.ConfigMap.Name == configMap.Name {
				return true, fmt.Sprintf("used by pod %s in volume %s", pod.Name, volume.Name)
			}
		}
		for _, container := range pod.Spec.Containers {
			for _, envFrom := range container.EnvFrom {
				if envFrom.ConfigMapRef != nil && envFrom.ConfigMapRef.Name == configMap.Name {
					return true, fmt.Sprintf("used by pod %s in container %s", pod.Name, container.Name)
				}
			}
			for _, env := range container.Env {
				if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil && env.ValueFrom.ConfigMapKeyRef.Name == configMap.Name {
					return true, fmt.Sprintf("used by pod %s in container %s", pod.Name, container.Name)
				}
			}
		}
	}
	return false, "Not used by any pod"
}

func GetUnusedConfigMaps(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var unusedConfigMaps []models.Resource
	pods, err := GetPods(clientset, namespace, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	configmaps, err := GetConfigMaps(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	for _, cm := range configmaps {
		used, reason := IsConfigMapUsed(cm, pods)
		if !used {
			unusedConfigMaps = append(unusedConfigMaps, models.Resource{
				Name:      cm.Name,
				Namespace: cm.Namespace,
				Kind:      "ConfigMap",
				Reason:    reason,
				Age:       time.Since(cm.CreationTimestamp.Time),
			})
		}
	}
	return unusedConfigMaps, nil
}
