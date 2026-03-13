package scanner

import (
	"testing"

	"github.com/bansikah22/kswp/test/mocks"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetOrphanServices(t *testing.T) {
	client := mocks.NewMockClient()
	orphanServices, err := GetOrphanServices(client.Clientset(), "", metav1.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, orphanServices, 1)
	assert.Equal(t, "service-2", orphanServices[0].Name)
}
