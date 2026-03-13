package scanner

import (
	"testing"

	"github.com/bansikah22/kswp/test/mocks"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetUnusedSecrets(t *testing.T) {
	client := mocks.NewMockClient()
	unusedSecrets, err := GetUnusedSecrets(client.Clientset(), "", metav1.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, unusedSecrets, 1)
	assert.Equal(t, "secret-2", unusedSecrets[0].Name)
}
