package scanner

import (
	"context"
	"time"

	"github.com/bansikah22/kswp/pkg/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetServices(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]v1.Service, error) {
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return services.Items, nil
}

func IsServiceOrphan(service v1.Service, endpoints []v1.Endpoints) (bool, string) {
	if service.Spec.Selector == nil {
		return false, "service has no selector"
	}
	for _, endpoint := range endpoints {
		if endpoint.Name == service.Name && len(endpoint.Subsets) > 0 {
			return false, "service has endpoints"
		}
	}
	return true, "No pods match selector"
}

func GetEndpoints(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]v1.Endpoints, error) {
	endpoints, err := clientset.CoreV1().Endpoints(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return endpoints.Items, nil
}

func GetOrphanServices(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var orphanServices []models.Resource
	services, err := GetServices(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	endpoints, err := GetEndpoints(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		orphan, reason := IsServiceOrphan(service, endpoints)
		if orphan {
			orphanServices = append(orphanServices, models.Resource{
				Name:      service.Name,
				Namespace: service.Namespace,
				Kind:      "Service",
				Reason:    reason,
				Age:       time.Since(service.CreationTimestamp.Time),
			})
		}
	}
	return orphanServices, nil
}
