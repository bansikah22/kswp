package scanner

import (
	"context"
	"time"

	"github.com/bansikah22/kswp/pkg/models"
	v1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
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

func IsServiceOrphan(service v1.Service, endpointSlices []discoveryv1.EndpointSlice) (bool, string) {
	if service.Spec.Selector == nil {
		return false, "service has no selector"
	}
	for _, slice := range endpointSlices {
		if slice.Labels["kubernetes.io/service-name"] == service.Name {
			for _, endpoint := range slice.Endpoints {
				if endpoint.Conditions.Ready != nil && *endpoint.Conditions.Ready {
					return false, "service has endpoints"
				}
			}
		}
	}
	return true, "No pods match selector"
}

func GetEndpointSlices(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]discoveryv1.EndpointSlice, error) {
	endpointSlices, err := clientset.DiscoveryV1().EndpointSlices(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return endpointSlices.Items, nil
}

func GetOrphanServices(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var orphanServices []models.Resource
	services, err := GetServices(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	endpointSlices, err := GetEndpointSlices(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if ShouldExclude(service.ObjectMeta) {
			continue
		}
		orphan, reason := IsServiceOrphan(service, endpointSlices)
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
