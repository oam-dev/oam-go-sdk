package types

import "k8s.io/apimachinery/pkg/runtime"

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

func (gvk GroupVersionKind) IsSame(ggvk GroupVersionKind) bool {
	return gvk.Group == ggvk.Group && gvk.Version == ggvk.Version && gvk.Kind == ggvk.Kind
}

type TemplateType string

func (t TemplateType) String() string {
	return string(t)
}

const (
	TemplateType_Type       = "type"
	TemplateType_GVInfo     = "groupversion_info"
	TemplateType_Controller = "controller"
	TemplateType_Main       = "main"
	TemplateType_Makefile   = "makefile"
)

type ResourceType string

const (
	ResourceType_Workload = "workload"
	ResourceType_Trait    = "trait"
	ResourceType_Exchange = "exchange"
)

type ExchangeGetter interface {
	runtime.Object
	GetExchange() []runtime.Object
}
