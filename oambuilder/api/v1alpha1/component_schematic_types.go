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

// +kubebuilder:validation:Optional
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type PortProtocol string

const (
	TCP PortProtocol = "TCP"
	UDP              = "UDP"
)

type WorkloadSetting struct {
	Name      string             `json:"name,omitempty"`
	Value     intstr.IntOrString `json:"value,omitempty"`
	Type      string             `json:"type,omitempty"`
	FromParam string             `json:"fromParam,omitempty"`
}

/// CPU describes a CPU resource allocation for a container.
///
/// The minimum number of logical cpus required for running this container.
type CPU struct {
	Required resource.Quantity `json:"required"`
}

/// Memory describes the memory allocation for a container.
///
/// The minimum amount of memory in MB required for running this container. The value should be a positive integer, greater than zero.
type Memory struct {
	Required resource.Quantity `json:"required"`
}

/// GPU describes a Container's need for a GPU.
///
/// The minimum number of gpus required for running this container.
type GPU struct {
	Required resource.Quantity `json:"required"`
}

/// Volume describes a path that is attached to a Container.
///
/// It specifies not only the location, but also the requirements.
type Volume struct {
	Name          string        `json:"name"`
	MountPath     string        `json:"mountPath"`
	AccessMode    AccessMode    `json:"accessMode"`
	SharingPolicy SharingPolicy `json:"sharingPolicy"`
	Disk          Disk          `json:"disk"`
}

/// AccessMode defines the access modes for file systems.
///
/// Currently, only read/write and read-only are supported.
type AccessMode string

const (
	RW AccessMode = "RW"
	RO AccessMode = "RO"
)

/// SharingPolicy defines whether a filesystem can be shared across containers.
///
/// An Exclusive filesystem can only be attached to one container.
type SharingPolicy string

const (
	Shared    SharingPolicy = "Shared"
	Exclusive SharingPolicy = "Exclusive"
)

// Disk describes the disk requirements for backing a Volume.
type Disk struct {
	Required  string `json:"required"`
	Ephemeral bool   `json:"ephemeral"`
}

type ExtendedResource struct {
	Name     string `json:"name"`
	Required string `json:"required"`
}

type Resources struct {
	Cpu      CPU                `json:"cpu"`
	Memory   Memory             `json:"memory"`
	Gpu      GPU                `json:"gpu"`
	Volumes  []Volume           `json:"volumes,omitempty"`
	Extended []ExtendedResource `json:"extended,omitempty"`
}

type Env struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	FromParam string `json:"fromParam"`
}

type Port struct {
	Name          string       `json:"name"`
	ContainerPort int32        `json:"port"`
	Protocol      PortProtocol `json:"protocol"`
}

type Exec struct {
	Command []string `json:"command,omitempty"`
}

type HttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HttpGet struct {
	Path        string       `json:"path"`
	Port        int32        `json:"port"`
	HttpHeaders []HttpHeader `json:"httpHeaders,omitempty"`
}

type TcpSocket struct {
	Port int32 `json:"port"`
}

type HealthProbe struct {
	Exec                Exec      `json:"exec"`
	HttpGet             HttpGet   `json:"httpGet"`
	TcpSocket           TcpSocket `json:"tcpSocket"`
	InitialDelaySeconds int32     `json:"initialDelaySeconds"`
	PeriodSeconds       int32     `json:"periodSeconds"`
	TimeoutSeconds      int32     `json:"timeoutSeconds"`
	SuccessThreshold    int32     `json:"successThreshold"`
	FailureThreshold    int32     `json:"failureThreshold"`
}

type ConfigFile struct {
	Path      string `json:"path"`
	Value     string `json:"value"`
	FromParam string `json:"fromParam"`
}

type Container struct {
	Name            string       `json:"name"`
	Image           string       `json:"image"`
	Resources       Resources    `json:"resources"`
	Cmd             []string     `json:"cmd,omitempty"`
	Args            []string     `json:"args,omitempty"`
	Env             []Env        `json:"env,omitempty"`
	Config          []ConfigFile `json:"config,omitempty"`
	Ports           []Port       `json:"ports,omitempty"`
	LivenessProbe   HealthProbe  `json:"livenessProbe"`
	ReadinessProbe  HealthProbe  `json:"readinessProbe"`
	ImagePullSecret string       `json:"imagePullSecret"`
}

type ParameterType string

const (
	Boolean ParameterType = "boolean"
	String  ParameterType = "string"
	Number  ParameterType = "number"
	Null    ParameterType = "null"
)

// parameter declaration
type Parameter struct {
	// The parameter's name. Must be unique per component.
	Name string `json:"name"`

	// A description of the parameter.
	// +optional
	Description string `json:"description,omitempty"`

	// The parameter's type. One of boolean, number, string, or null
	// as defined in the JSON specification and the JSON Schema Validation spec
	ParameterType ParameterType `json:"type"`

	// Whether a value must be provided for the parameter.
	// Default is false.
	// +optional
	Required bool `json:"required,omitempty"`

	// The parameter's default value.
	// type indicated by type field.
	// +optional
	Default string `json:"default,omitempty"`
}

// ComponentSpec defines the desired state of ComponentSchematic
type ComponentSpec struct {
	Parameters       []Parameter       `json:"parameters,omitempty"`
	WorkloadType     string            `json:"workloadType"`
	OsType           string            `json:"osType"`
	Arch             string            `json:"arch"`
	Containers       []Container       `json:"containers,omitempty"`
	WorkloadSettings []WorkloadSetting `json:"workloadSettings,omitempty"`
}

type ComponentStatus struct {
}

// +genclient

// ComponentSchematic is the Schema for the components API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ComponentSchematic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec,omitempty"`
	Status ComponentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ComponentSchematicList contains a list of ComponentSchematic
type ComponentSchematicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComponentSchematic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ComponentSchematic{}, &ComponentSchematicList{})
}
