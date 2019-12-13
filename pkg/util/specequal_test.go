package util

import (
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestSpecEqual(t *testing.T) {
	spec := v1.DeploymentSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "test",
					Image: "test/test",
				}},
			},
		},
	}
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	eq := SpecEqual(deployment, spec, true)
	require.False(t, eq)
	eq = SpecEqual(deployment, spec, true)
	spec.Template.Spec.Containers[0].Name = "new-test"
	eq = SpecEqual(deployment, spec, true)
	require.False(t, eq)
	eq = SpecEqual(deployment, spec, true)
	require.True(t, eq)
}
