/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type Names struct {
	Kind string `json:"kind"`
	// +optional
	Singular string `json:"singular,omitempty"`
	// +optional
	Plural string `json:"plural,omitempty"`
}

// WorkloadTypeSpec defines the desired state of WorkloadType
type WorkloadTypeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The group that this workload type belongs to (e.g. core.hydra.io)
	//Group string `json:"group"`

	Names Names `json:"names"`

	// Workload type, GVK
	Group   string `json:"group"`
	Version string `json:"version"`

	// The workload type's settings options.
	// +optional
	Settings string `json:"settings,omitempty"`
}

// WorkloadTypeStatus defines the observed state of WorkloadType
type WorkloadTypeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// WorkloadType is the Schema for the workloadtypes API
type WorkloadType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkloadTypeSpec   `json:"spec,omitempty"`
	Status WorkloadTypeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WorkloadTypeList contains a list of WorkloadType
type WorkloadTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkloadType `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WorkloadType{}, &WorkloadTypeList{})
}
