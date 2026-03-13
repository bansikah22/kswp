package scanner

import (
	"testing"

	"github.com/bansikah22/kswp/test/mocks"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetUnusedConfigMaps(t *testing.T) {
	client := mocks.NewMockClient()
	unusedConfigMaps, err := GetUnusedConfigMaps(client.Clientset(), "", metav1.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, unusedConfigMaps, 1)
	assert.Equal(t, "cm-2", unusedConfigMaps[0].Name)
}
