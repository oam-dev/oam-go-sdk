package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	policyv1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// ComponentConfiguration status
const (
	StatusReady       = "Ready"
	StatusProgressing = "Progressing"
	StatusFailed      = "Failed"
)

// Update component status with specific status and meta info.
func (s *ModuleStatus) Update(rsrc metav1.Object, status string) {
	ro := rsrc.(runtime.Object)
	gvk := ro.GetObjectKind().GroupVersionKind()
	s.NamespacedName = rsrc.GetNamespace() + string(types.Separator) + rsrc.GetName()
	s.GroupVersion = gvk.GroupVersion().String()
	s.Kind = gvk.GroupKind().Kind
	s.Status = status
}

// ResetComponentList - reset component list objects
func (m *ApplicationConfigurationStatus) resetComponentList() {
	m.Modules = []ModuleStatus{}
}

// Update App Status accord the components status
func (m *ApplicationConfigurationStatus) Update(rsrcs []metav1.Object, err error) {
	var ready = true
	m.resetComponentList()
	// compute components status
	for _, r := range rsrcs {
		os := ModuleStatus{}
		os.Update(r, StatusReady)
		switch r.(type) {
		case *appsv1.StatefulSet:
			os.Status = stsStatus(r.(*appsv1.StatefulSet))
		case *policyv1.PodDisruptionBudget:
			os.Status = pdbStatus(r.(*policyv1.PodDisruptionBudget))
		case *appsv1.Deployment:
			os.Status = deploymentStatus(r.(*appsv1.Deployment))
		case *appsv1.ReplicaSet:
			os.Status = replicasetStatus(r.(*appsv1.ReplicaSet))
		case *appsv1.DaemonSet:
			os.Status = daemonsetStatus(r.(*appsv1.DaemonSet))
		case *corev1.Pod:
			os.Status = podStatus(r.(*corev1.Pod))
		case *corev1.Service:
			os.Status = serviceStatus(r.(*corev1.Service))
		case *corev1.PersistentVolumeClaim:
			os.Status = pvcStatus(r.(*corev1.PersistentVolumeClaim))
		case *v1beta1.Ingress:
			os.Status = ingressStatus(r.(*v1beta1.Ingress))
		default:
			os.Status = StatusReady
		}
		m.Modules = append(m.Modules, os)
	}

	// aggregate
	for _, os := range m.Modules {
		if os.Status != StatusReady {
			ready = false
		}
	}
	if ready {
		m.Phase = ApplicationReady
		m.Ready("ComponentsReady", "all components ready")
	} else {
		m.Phase = ApplicationProgressing
		m.NotReady("ComponentsNotReady", "some components not ready")
	}
	if err != nil {
		m.SetConditionTrue(Error, "ErrorSeen", err.Error())
	}
}

// Resource specific logic -----------------------------------

// Statefulset
func stsStatus(rsrc *appsv1.StatefulSet) string {
	if rsrc.Status.ReadyReplicas == *rsrc.Spec.Replicas && rsrc.Status.CurrentReplicas == *rsrc.Spec.Replicas {
		return StatusReady
	}
	return StatusProgressing
}

// Deployment
func deploymentStatus(rsrc *appsv1.Deployment) string {
	status := StatusProgressing
	progress := true
	available := true
	for _, c := range rsrc.Status.Conditions {
		switch c.Type {
		case appsv1.DeploymentProgressing:
			// https://github.com/kubernetes/kubernetes/blob/a3ccea9d8743f2ff82e41b6c2af6dc2c41dc7b10/pkg/controller/deployment/progress.go#L52
			if c.Status != corev1.ConditionTrue || c.Reason != "NewReplicaSetAvailable" {
				progress = false
			}
		case appsv1.DeploymentAvailable:
			if c.Status == corev1.ConditionFalse {
				available = false
			}
		}
	}

	if progress && available {
		status = StatusReady
	}

	return status
}

// Replicaset
func replicasetStatus(rsrc *appsv1.ReplicaSet) string {
	status := StatusProgressing
	failure := false
	for _, c := range rsrc.Status.Conditions {
		switch c.Type {
		// https://github.com/kubernetes/kubernetes/blob/a3ccea9d8743f2ff82e41b6c2af6dc2c41dc7b10/pkg/controller/replicaset/replica_set_utils.go
		case appsv1.ReplicaSetReplicaFailure:
			if c.Status == corev1.ConditionTrue {
				failure = true
				break
			}
		}
	}

	if !failure && rsrc.Status.ReadyReplicas == rsrc.Status.Replicas && rsrc.Status.Replicas == rsrc.Status.AvailableReplicas {
		status = StatusReady
	}

	return status
}

// Daemonset
func daemonsetStatus(rsrc *appsv1.DaemonSet) string {
	status := StatusProgressing
	if rsrc.Status.DesiredNumberScheduled == rsrc.Status.NumberAvailable && rsrc.Status.DesiredNumberScheduled == rsrc.Status.NumberReady {
		status = StatusReady
	}
	return status
}

// PVC
func pvcStatus(rsrc *corev1.PersistentVolumeClaim) string {
	status := StatusProgressing
	if rsrc.Status.Phase == corev1.ClaimBound {
		status = StatusReady
	}
	return status
}

// Service
func serviceStatus(rsrc *corev1.Service) string {
	status := StatusReady
	if rsrc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		// For LoadBalancer, we need to wait ingress bind
		if len(rsrc.Status.LoadBalancer.Ingress) == 0 {
			// if no bind
			status = StatusProgressing
		}
	}
	return status
}

// Ingress
func ingressStatus(rsrc *v1beta1.Ingress) string {
	status := StatusReady
	if len(rsrc.Status.LoadBalancer.Ingress) == 0 {
		// if no bind
		status = StatusProgressing
	}
	return status
}

// Pod
func podStatus(rsrc *corev1.Pod) string {
	status := StatusProgressing
	for i := range rsrc.Status.Conditions {
		if rsrc.Status.Conditions[i].Type == corev1.PodReady &&
			rsrc.Status.Conditions[i].Status == corev1.ConditionTrue {
			status = StatusReady
			break
		}
	}
	return status
}

// PodDisruptionBudget
func pdbStatus(rsrc *policyv1.PodDisruptionBudget) string {
	if rsrc.Status.CurrentHealthy >= rsrc.Status.DesiredHealthy {
		return StatusReady
	}
	return StatusProgressing
}
