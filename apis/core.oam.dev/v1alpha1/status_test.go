package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func crdStatus(r metav1.Object) string {
	rsrc, ok := r.(*ApplicationConfiguration)
	if !ok {
		return StatusUnknown
	}
	if rsrc.Status.Phase == "Ready" {
		return StatusReady
	}
	return StatusProgressing
}

func TestRegisterStatusHandler(t *testing.T) {
	sts := new(ApplicationConfiguration)
	sts.Status.Phase = "Ready"
	RegisterStatusHandler(sts.GetObjectKind().GroupVersionKind(), crdStatus)
	as := new(ApplicationConfigurationStatus)
	as.Update([]metav1.Object{sts, new(ComponentSchematic)}, nil)
	assert.Equal(t, StatusReady, as.Modules[0].Status)
	assert.Equal(t, StatusUnknown, as.Modules[1].Status)
	assert.Equal(t, StatusProgressing, string(as.Phase))
}
