package mocks

import (
	"time"

	"github.com/bansikah22/kswp/internal/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

type MockClient struct {
	clientset k8s.Interface
}

func (c *MockClient) Clientset() k8s.Interface {
	return c.clientset
}

func NewMockClient() kubernetes.Client {
	return &MockClient{
		clientset: fake.NewSimpleClientset(
			&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-1",
					Namespace: "default",
					Labels: map[string]string{
						"app": "app-1",
					},
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "config-volume",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "cm-1",
									},
								},
							},
						},
						{
							Name: "secret-volume",
							VolumeSource: v1.VolumeSource{
								Secret: &v1.SecretVolumeSource{
									SecretName: "secret-1",
								},
							},
						},
					},
				},
			},
			&v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cm-1",
					Namespace: "default",
				},
			},
			&v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cm-2",
					Namespace: "default",
				},
			},
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-1",
					Namespace: "default",
				},
			},
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-2",
					Namespace: "default",
				},
			},
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service-1",
					Namespace: "default",
				},
				Spec: v1.ServiceSpec{
					Selector: map[string]string{
						"app": "app-1",
					},
				},
			},
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service-2",
					Namespace: "default",
				},
				Spec: v1.ServiceSpec{
					Selector: map[string]string{
						"app": "app-2",
					},
				},
			},
			&v1.Endpoints{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service-1",
					Namespace: "default",
				},
				Subsets: []v1.EndpointSubset{
					{
						Addresses: []v1.EndpointAddress{
							{
								IP: "1.1.1.1",
							},
						},
					},
				},
			},
			&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment-1",
					Namespace: "default",
				},
			},
			&appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rs-1",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Deployment",
							Name: "deployment-1",
						},
					},
				},
				Spec: appsv1.ReplicaSetSpec{
					Replicas: func() *int32 { r := int32(0); return &r }(),
				},
			},
			&appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rs-2",
					Namespace: "default",
				},
				Spec: appsv1.ReplicaSetSpec{
					Replicas: func() *int32 { r := int32(0); return &r }(),
				},
			},
			&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job-1",
					Namespace: "default",
				},
				Status: batchv1.JobStatus{
					Succeeded:      1,
					CompletionTime: &metav1.Time{Time: time.Now().Add(-25 * time.Hour)},
				},
			},
			&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job-2",
					Namespace: "default",
				},
				Status: batchv1.JobStatus{
					Succeeded:      1,
					CompletionTime: &metav1.Time{Time: time.Now()},
				},
			},
		),
	}
}
