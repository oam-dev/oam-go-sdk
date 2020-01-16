package trait

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkloadReference struct {
	metav1.OwnerReference `json:",inline"`
	Namespace             string `json:"namespace"`
}

func (in *WorkloadReference) DeepCopyInto(out *WorkloadReference) {
	*out = *in
	if in.Controller != nil {
		in, out := &in.Controller, &out.Controller
		*out = new(bool)
		**out = **in
	}
	if in.BlockOwnerDeletion != nil {
		in, out := &in.BlockOwnerDeletion, &out.BlockOwnerDeletion
		*out = new(bool)
		**out = **in
	}
	return
}

func (in *WorkloadReference) DeepCopy() *WorkloadReference {
	if in == nil {
		return nil
	}
	out := new(WorkloadReference)
	in.DeepCopyInto(out)
	return out
}
