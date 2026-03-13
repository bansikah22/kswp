package scanner

import (
	"testing"

	"github.com/bansikah22/kswp/test/mocks"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetOldReplicaSets(t *testing.T) {
	client := mocks.NewMockClient()
	oldReplicaSets, err := GetOldReplicaSets(client.Clientset(), "", metav1.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, oldReplicaSets, 1)
	assert.Equal(t, "rs-2", oldReplicaSets[0].Name)
}
