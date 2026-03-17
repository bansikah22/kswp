package scanner

import (
	"testing"
	"time"

	"github.com/bansikah22/kswp/test/mocks"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetCompletedJobs(t *testing.T) {
	client := mocks.NewMockClient()
	completedJobs, err := GetCompletedJobs(client.Clientset(), 24*time.Hour, "", metav1.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, completedJobs, 2)
	assert.Contains(t, []string{completedJobs[0].Name, completedJobs[1].Name}, "job-1")
	assert.Contains(t, []string{completedJobs[0].Name, completedJobs[1].Name}, "job-3")
}
