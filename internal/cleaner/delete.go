package cleaner

import (
	"context"
	"fmt"

	"github.com/bansikah22/kswp/pkg/models"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeleteResource(clientset kubernetes.Interface, resource models.Resource) error {
	var err error
	switch resource.Kind {
	case "ConfigMap":
		err = clientset.CoreV1().ConfigMaps(resource.Namespace).Delete(context.TODO(), resource.Name, v1.DeleteOptions{})
	case "Secret":
		err = clientset.CoreV1().Secrets(resource.Namespace).Delete(context.TODO(), resource.Name, v1.DeleteOptions{})
	case "Service":
		err = clientset.CoreV1().Services(resource.Namespace).Delete(context.TODO(), resource.Name, v1.DeleteOptions{})
	case "ReplicaSet":
		err = clientset.AppsV1().ReplicaSets(resource.Namespace).Delete(context.TODO(), resource.Name, v1.DeleteOptions{})
	case "Job":
		err = clientset.BatchV1().Jobs(resource.Namespace).Delete(context.TODO(), resource.Name, v1.DeleteOptions{})
	case "PersistentVolumeClaim":
		err = clientset.CoreV1().PersistentVolumeClaims(resource.Namespace).Delete(context.TODO(), resource.Name, v1.DeleteOptions{})
	default:
		return fmt.Errorf("unknown resource kind: %s", resource.Kind)
	}

	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("%s %s/%s not found, skipping\n", resource.Kind, resource.Namespace, resource.Name)
			return nil
		}
		return err
	}

	fmt.Printf("%s %s/%s deleted\n", resource.Kind, resource.Namespace, resource.Name)
	return nil
}
