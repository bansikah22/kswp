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

func GetSecrets(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]v1.Secret, error) {
	secrets, err := clientset.CoreV1().Secrets(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return secrets.Items, nil
}

func IsSecretUsed(secret v1.Secret, pods []v1.Pod) (bool, string) {
	for _, pod := range pods {
		for _, volume := range pod.Spec.Volumes {
			if volume.Secret != nil && volume.Secret.SecretName == secret.Name {
				return true, fmt.Sprintf("used by pod %s in volume %s", pod.Name, volume.Name)
			}
		}
		for _, container := range pod.Spec.Containers {
			for _, envFrom := range container.EnvFrom {
				if envFrom.SecretRef != nil && envFrom.SecretRef.Name == secret.Name {
					return true, fmt.Sprintf("used by pod %s in container %s", pod.Name, container.Name)
				}
			}
			for _, env := range container.Env {
				if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil && env.ValueFrom.SecretKeyRef.Name == secret.Name {
					return true, fmt.Sprintf("used by pod %s in container %s", pod.Name, container.Name)
				}
			}
		}
	}
	return false, "Not used by any pod"
}

func GetUnusedSecrets(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var unusedSecrets []models.Resource
	pods, err := GetPods(clientset, namespace, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	secrets, err := GetSecrets(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	for _, secret := range secrets {
		if ShouldExclude(secret.ObjectMeta) {
			continue
		}
		used, reason := IsSecretUsed(secret, pods)
		if !used {
			unusedSecrets = append(unusedSecrets, models.Resource{
				Name:      secret.Name,
				Namespace: secret.Namespace,
				Kind:      "Secret",
				Reason:    reason,
				Age:       time.Since(secret.CreationTimestamp.Time),
			})
		}
	}
	return unusedSecrets, nil
}
