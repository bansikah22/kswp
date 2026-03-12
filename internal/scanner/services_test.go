package scanner

import (
	"testing"

	"github.com/bansikah22/kswp/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetOrphanServices(t *testing.T) {
	client := mocks.NewMockClient()
	orphanServices, err := GetOrphanServices(client.Clientset(), "")
	assert.NoError(t, err)
	assert.Len(t, orphanServices, 1)
	assert.Equal(t, "service-2", orphanServices[0].Name)
}
