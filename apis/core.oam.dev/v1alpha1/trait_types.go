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

// TraitSpec defines the desired state of Trait
type TraitSpec struct {
	Group   string `json:"type"`
	Version string `json:"version"`

	Names Names `json:"names"`

	// The list of workload types that this trait applies to. "*" means any workload type
	// Default is ['*']
	// +optional
	AppliesTo []string `json:"appliesTo,omitempty"`

	// description of json schema
	// +optional
	Properties string `json:"properties,omitempty"`
}

// TraitStatus defines the observed state of Trait
type TraitStatus struct {
}

// +genclient
// +kubebuilder:object:root=true
// Trait is the Schema for the traits API
type Trait struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TraitSpec   `json:"spec,omitempty"`
	Status TraitStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TraitList contains a list of Trait
type TraitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Trait `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Trait{}, &TraitList{})
}
