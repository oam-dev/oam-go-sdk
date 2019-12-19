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
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplicationConfigurationSpec defines the desired state of ApplicationConfiguration
type ApplicationConfigurationSpec struct {
	//  +optional
	Variables []Variable `json:"variables,omitempty"`
	//  +optional
	Scopes     []ScopeBinding           `json:"scopes,omitempty"`
	Components []ComponentConfiguration `json:"components"`
}

/// A value that is substituted into a parameter.
type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ScopeBinding struct {
	Name string `json:"name"`
	Type string `json:"type"`

	// A properties object (for trait and scope configuration) is an object whose structure is determined by the trait or scope property schema. It may be a simple value, or it may be a complex object.
	// Properties are validated against the schema appropriate for the trait or scope.
	// +optional
	// Properties runtime.RawExtension `json:"properties,omitempty"`
	Properties runtime.RawExtension `json:"properties,omitempty"`
}

// ApplicationConfigurationStatus defines the observed state of ApplicationConfiguration
type ApplicationConfigurationStatus struct {
	// The phase of a application is a simple, high-level summary of where the whole  Application is in its lifecycle.
	// The conditions array contains more detail about the appConf's status.
	// There are five possible phase values:
	// Pending: The Application has been accepted by the Kubernetes system, but not get processed by EDAS.
	// Progressing: The Application has been processed by EDAS, related resources provision are progressing.
	// Ready: All related resources provision are ready, application is on serving.
	// Failed: Occur some failures in the process of creating Application, you can get detail infos from Conditions.
	// Unknown: For some reason the state of the Application could not be obtained, typically due to an
	// error in controller.
	// +optional
	Phase ApplicationPhase `json:"phase,omitempty"`

	// Module status array for all modules constitute this application.
	// Module is k8s build-in or CRD object, only show cswt level.
	// +optional
	Modules []ModuleStatus `json:"modules,omitempty"`

	// Represents the latest available observations of a application's current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ApplicationCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// ApplicationPhase is a label for the condition of a Application at the current time.
type ApplicationPhase string

// These are valid status of a Application
const (
	// ApplicationPending means the Application has been accepted by the Kubernetes system,
	// but not get processed by EDAS.
	ApplicationPending ApplicationPhase = "Pending"

	// ApplicationProgressing means the Application has been processed by EDAS,
	// related resources provision are progressing.
	ApplicationProgressing ApplicationPhase = "Progressing"

	// ApplicationReady means All related resources provision are ready, application is on serving.
	ApplicationReady ApplicationPhase = "Ready"

	// ApplicationFailed means some failures occurred in the process of creating Application,
	// you can get detail info from Conditions.
	ApplicationFailed ApplicationPhase = "Failed"
)

// ModuleStatus is a generic status holder for components
// +k8s:deepcopy-gen=true
type ModuleStatus struct {
	// NamespacedName of component
	NamespacedName string `json:"name,omitempty"`
	// Kind of component
	Kind string `json:"kind,omitempty"`
	// ComponentConfiguration groupVersion
	GroupVersion string `json:"groupVersion,omitempty"`
	// Status. Values: Progressing, Ready, Failed
	Status string `json:"status,omitempty"`
}

// +k8s:deepcopy-gen=true
type ApplicationCondition struct {
	// Type of Application condition.
	Type ApplicationConditionType `json:"type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	// - True means application in this condition type
	// - False means application not in this condition type
	// - Unknown means whether application in this condition type is unknown
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

type ApplicationConditionType string

// These are valid conditions of a application.
const (
	// Ready type shows application's readiness condition.
	Ready ApplicationConditionType = "Ready"

	// Cleanup type shows application's cleanup condition.
	Cleanup ApplicationConditionType = "Cleanup"

	// Error type shows application's error condition.
	Error ApplicationConditionType = "Error"
)

// +kubebuilder:object:root=true
// ApplicationConfiguration is the Schema for the operationalconfigurations API
type ApplicationConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ApplicationConfigurationSpec `json:"spec,omitempty"`
	// +optional
	Status ApplicationConfigurationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// ApplicationConfigurationList contains a list of ApplicationConfiguration
type ApplicationConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApplicationConfiguration{}, &ApplicationConfigurationList{})
}

type ComponentConfiguration struct {
	ComponentName string `json:"componentName"`
	InstanceName  string `json:"instanceName"`
	// extension field, workload reference name
	// TODO this should be removed as spec didn't have
	// +optional
	RefName string `json:"refName,omitempty"`
	// +optional
	ParameterValues []ParameterValue `json:"parameterValues,omitempty"`
	// +optional
	Traits []TraitBinding `json:"traits,omitempty"`
	// +optional
	ApplicationScopes []string `json:"applicationScopes,omitempty"`
}

type TraitBinding struct {
	Name string `json:"name"`
	// TODO this is extension field, should be added to spec or removed
	// +optional
	InstanceName string `json:"instanceName,omitempty"`
	// TODO this is extension field, trait resource reference name,should be removed
	// +optional
	RefName string `json:"refName,omitempty"`
	// TODO: change to Value
	// +optional
	Properties runtime.RawExtension `json:"properties,omitempty"`
}

// Check whether this componet configured specific trait.
func (c *ComponentConfiguration) ExistTrait(t string) bool {
	name, _, _ := c.ExtractTrait(t)
	return name != ""
}

// Get specific trait's Full name of this component and its parameterValues.
// If not exist, name is "" and parameterValues is nil.
// bool mark whether this trait has ref name.
func (c *ComponentConfiguration) ExtractTrait(t string) (string, bool, []ParameterValue) {
	var pvals []ParameterValue
	name := ""
	existing := c.Traits
	isRef := false
	for _, f := range existing {
		if f.Name == t {
			if f.RefName != "" {
				name = f.RefName
				isRef = true
			} else {
				name = c.InstanceName + "-" + t
			}
			pvals = parseParamValues(f.Properties.Raw)
			break
		}
	}
	return name, isRef, pvals
}

func parseParamValues(data []byte) []ParameterValue {
	pvals := make([]ParameterValue, 0)
	err := json.Unmarshal(data, &pvals)
	if err != nil {
		panic("json Unmarshal failed")
	}
	return pvals
}

func (c *ComponentConfiguration) GenTraitName(appConf *ApplicationConfiguration, t string) string {
	name, isRef, _ := c.ExtractTrait(t)
	if isRef {
		// if has reference name, use it.
		return name
	}
	return appConf.Name + "-" + name
}

func (c *ApplicationConfiguration) SetComponent(component *ComponentConfiguration) {
	newComponents := []ComponentConfiguration{}
	for _, v := range c.Spec.Components {
		if v.ComponentName == component.ComponentName {
			newComponents = append(newComponents, *component)

		} else {
			newComponents = append(newComponents, v)

		}
	}
	c.Spec.Components = newComponents

}

func (c *ComponentConfiguration) SetReplicas(s string) {
	newParameterValues := []ParameterValue{}
	for _, e := range c.ParameterValues {
		if e.Name == "replicas" {
			e.Value = s
		}
		newParameterValues = append(newParameterValues, e)

	}
	c.ParameterValues = newParameterValues
}
