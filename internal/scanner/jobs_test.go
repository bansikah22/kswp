package scanner

import (
	"testing"
	"time"

	"github.com/bansikah22/kswp/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetCompletedJobs(t *testing.T) {
	client := mocks.NewMockClient()
	completedJobs, err := GetCompletedJobs(client.Clientset(), 24*time.Hour, "")
	assert.NoError(t, err)
	assert.Len(t, completedJobs, 1)
	assert.Equal(t, "job-1", completedJobs[0].Name)
}
