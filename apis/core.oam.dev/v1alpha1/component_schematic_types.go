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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type PortProtocol string

const (
	TCP PortProtocol = "TCP"
	UDP              = "UDP"
)

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
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	// +optional
	AccessMode AccessMode `json:"accessMode,omitempty"`
	// +optional
	SharingPolicy SharingPolicy `json:"sharingPolicy,omitempty"`
	// +optional
	Disk *Disk `json:"disk,omitempty"`
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
	Ephemeral bool   `json:"ephemeral,omitempty"`
}

// ExtendedResource give extension ability
type ExtendedResource struct {
	Name     string `json:"name"`
	Required string `json:"required"`
}

/// Resources defines the resources required by a container.
type Resources struct {
	Cpu    CPU    `json:"cpu"`
	Memory Memory `json:"memory"`
	// +optional
	Gpu GPU `json:"gpu,omitempty"`
	// +optional
	Volumes []Volume `json:"volumes,omitempty"`
	// +optional
	Extended []ExtendedResource `json:"extended,omitempty"`
}

// Env describes an environment variable for a container.
type Env struct {
	Name string `json:"name"`
	// +optional
	Value string `json:"value,omitempty"`
	// +optional
	FromParam string `json:"fromParam,omitempty"`
}

// Port describes a port on a Container.
type Port struct {
	Name          string `json:"name"`
	ContainerPort int32  `json:"containerPort"`
	// +optional
	Protocol PortProtocol `json:"protocol,omitempty"`
}

type Exec struct {
	Command []string `json:"command"`
}

type HttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HttpGet struct {
	Path string `json:"path"`
	Port int32  `json:"port"`
	// +optional
	HttpHeaders []HttpHeader `json:"httpHeaders"`
}

/// TcpSocket defines a socket used for health probing.
type TcpSocket struct {
	Port int32 `json:"port"`
}

// HealthProbe describes a probe used to check on the health of a Container.
type HealthProbe struct {
	// +optional
	Exec *Exec `json:"exec,omitempty"`
	// +optional
	HttpGet *HttpGet `json:"httpGet,omitempty"`
	// +optional
	TcpSocket *TcpSocket `json:"tcpSocket,omitempty"`
	// +optional
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`
	// +optional
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`
	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`
	// +optional
	SuccessThreshold int32 `json:"successThreshold,omitempty"`
	// +optional
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
}

// ConfigFile describes locations to write configuration as files accessible within the container
type ConfigFile struct {
	Path string `json:"path"`
	// +optional
	Value string `json:"value,omitempty"`
	// +optional
	FromParam string `json:"fromParam,omitempty"`
}

// Container describes the container configuration for a Component.
type Container struct {
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Resources Resources `json:"resources"`
	// +optional
	Cmd []string `json:"cmd,omitempty"`
	// +optional
	Args []string `json:"args,omitempty"`
	// +optional
	Env []Env `json:"env,omitempty"`
	// +optional
	Config []ConfigFile `json:"config,omitempty"`
	// +optional
	Ports []Port `json:"ports,omitempty"`
	// +optional
	LivenessProbe *HealthProbe `json:"livenessProbe,omitempty"`
	// +optional
	ReadinessProbe *HealthProbe `json:"readinessProbe,omitempty"`
	// +optional
	ImagePullSecret string `json:"imagePullSecret,omitempty"`
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

type Expose struct {
	Name string `json:"name"`
}

type Consume struct {
	Name string `json:"name"`
	// +optional
	As string `json:"as"`
}

// ComponentSpec defines the desired state of ComponentSchematic
type ComponentSpec struct {
	// +optional
	Parameters   []Parameter `json:"parameters,omitempty"`
	WorkloadType string      `json:"workloadType"`
	// +optional
	OsType string `json:"osType,omitempty"`
	// +optional
	Arch string `json:"arch,omitempty"`
	// +optional
	Containers []Container `json:"containers,omitempty"`

	// +optional
	Expose []Expose `json:"expose,omitempty"`
	// +optional
	Consume []Consume `json:"consume,omitempty"`
	// +optional
	WorkloadSettings runtime.RawExtension `json:"workloadSettings,omitempty"`
}

type ComponentStatus struct {
}

// +genclient

// ComponentSchematic is the Schema for the components API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ComponentSchematic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ComponentSpec `json:"spec"`
	// +optional
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
