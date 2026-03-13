package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/bansikah22/kswp/pkg/models"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetJobs(clientset kubernetes.Interface, namespace string, listOptions metav1.ListOptions) ([]batchv1.Job, error) {
	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return jobs.Items, nil
}

func IsJobCompleted(job batchv1.Job, threshold time.Duration) (bool, string) {
	if job.Status.Succeeded > 0 {
		if job.Status.CompletionTime != nil {
			if job.Status.CompletionTime.Time.Before(time.Now().Add(-threshold)) {
				return true, fmt.Sprintf("Completed more than %s ago", threshold.String())
			}
			return false, fmt.Sprintf("Completed less than %s ago", threshold.String())
		}
	}
	return false, "Job has not completed"
}

func GetCompletedJobs(clientset kubernetes.Interface, threshold time.Duration, namespace string, listOptions metav1.ListOptions) ([]models.Resource, error) {
	var completedJobs []models.Resource
	jobs, err := GetJobs(clientset, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	for _, job := range jobs {
		completed, reason := IsJobCompleted(job, threshold)
		if completed {
			completedJobs = append(completedJobs, models.Resource{
				Name:      job.Name,
				Namespace: job.Namespace,
				Kind:      "Job",
				Reason:    reason,
				Age:       time.Since(job.CreationTimestamp.Time),
			})
		}
	}
	return completedJobs, nil
}
