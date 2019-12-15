package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ApplicationScopeSpec struct {
	Type                  string      `json:"type"`
	AllowComponentOverlap bool        `json:"allowComponentOverlap"`
	Parameters            []Parameter `json:"parameters"`
}

type ApplicationScopeStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApplicationScope struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ApplicationScopeSpec `json:"spec,omitempty"`
	// +optional
	Status ApplicationScopeStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApplicationScopeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationScope `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApplicationScope{}, &ApplicationScopeList{})
}
