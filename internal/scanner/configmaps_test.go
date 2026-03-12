package scanner

import (
	"testing"

	"github.com/bansikah22/kswp/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetUnusedConfigMaps(t *testing.T) {
	client := mocks.NewMockClient()
	unusedConfigMaps, err := GetUnusedConfigMaps(client.Clientset(), "")
	assert.NoError(t, err)
	assert.Len(t, unusedConfigMaps, 1)
	assert.Equal(t, "cm-2", unusedConfigMaps[0].Name)
}
