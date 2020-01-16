package runtime

import (
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types/trait"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types/workload"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type dynamicWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              dynamicWorkloadSpec `json:"spec"`
}

type dynamicWorkloadSpec struct {
	Settings map[string]interface{} `json:"settings,omitempty"`
	Traits   []workload.Trait       `json:"traits,omitempty"`
}

type dynamicTrait struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              dynamicTraitSpec `json:"spec"`
}

type dynamicTraitSpec struct {
	Settings map[string]interface{}  `json:"settings,omitempty"`
	Workload trait.WorkloadReference `json:"workload"`
}

type dynamicTraitList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Items             []dynamicTrait `json:"items"`
}

// TraitExchanger for trait
type TraitExchanger interface {
	Resources(cli client.Client, self runtime.Object, exchangeListResource runtime.Object) error
	Create(cli client.Client, self runtime.Object, exchangeResource runtime.Object) error
	Delete(cli client.Client, self runtime.Object, exchangeResource runtime.Object) error
}

// WorkloadExchanger for workload
type WorkloadExchanger interface {
	Resources(cli client.Client, self runtime.Object, exchangeListResource runtime.Object) error
	AllExchangerReady(cli client.Client, self runtime.Object) (bool, error)
}
