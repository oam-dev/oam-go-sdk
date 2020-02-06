package v1alpha1

import (
	"testing"

	"github.com/oam-dev/oam-go-sdk/apis/flags"

	v1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"

	"github.com/oam-dev/oam-go-sdk/apis/handlers"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNoModules(t *testing.T) {
	as := new(ApplicationConfigurationStatus)
	as.Update([]metav1.Object{}, nil)
	assert.Equal(t, flags.StatusProgressing, string(as.Phase))
}

func crdStatus(r metav1.Object) string {
	rsrc, ok := r.(*v1.Job)
	if !ok {
		return flags.StatusUnknown
	}
	if rsrc.Status.Failed == 0 && rsrc.Status.Succeeded > 0 {
		return flags.StatusReady
	}
	return flags.StatusProgressing
}

func TestRegisterStatusHandler(t *testing.T) {
	jb := new(v1.Job)
	jb.Status.Failed = 0
	jb.Status.Succeeded = 1
	handlers.RegisterStatusHandler(jb.GetObjectKind().GroupVersionKind(), crdStatus)
	as := new(ApplicationConfigurationStatus)
	as.Update([]metav1.Object{jb, new(v1beta1.CronJob)}, nil)
	assert.Equal(t, flags.StatusReady, as.Modules[0].Status)
	assert.Equal(t, flags.StatusUnknown, as.Modules[1].Status)
	assert.Equal(t, flags.StatusProgressing, string(as.Phase))
}
