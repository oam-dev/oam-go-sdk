package oam

import (
	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Action struct {
	// plugin do action, e.g: k8s, helm, ...
	Provider PType
	// action command, e.g: Create, Update
	Command CmdType
	// action content, for k8s plugin, this is k8s object, for helm plugin, this is helm chart address.
	Plan interface{}
}

type PType string

const (
	PTypeK8S PType = "k8s"
)

type EType string

const (
	CreateOrUpdate EType = "CreateOrUpdate"
	Delete         EType = "Delete"
)

type CmdType string

const (
	CmdTypeUpdate CmdType = "Update"
	CmdTypeCreate CmdType = "Create"
	CmdTypeDelete CmdType = "Delete"
)

// Hook triggered by application configuration modify event and execute before or after handlers.

// Side effects should be occurred in hook and state should be stored in ctx, the framework will pass
// ctx along this reconcile.

// ac is modified ApplicationConfiguration.

type Hook interface {
	Identity
	Exec(ctx *ActionContext, ac runtime.Object, EventType EType) error
}

// Handler triggered by components, traits, scopes modify event, actions should be generate and add to ctx.
// For actions need to be processed early, use ctx.AddPre; for actions need to be processed late, use ctx.AddPost.
// For normal actions, just  use ctx.Add, OAM framework will do preActions -> actions -> postActions for you.

type Handler interface {
	Identity
	Handle(ctx *ActionContext, ac runtime.Object, EventType EType) error
}

type Identity interface {
	Id() string
}

// SType for spec type
type SType string

const (
	STypeComponent                = "component"
	STypeScope                    = "scope"
	STypeWorkloadType             = "workloadType"
	STypeApplicationConfiguration = "applicationConfiguration"
	STypeTrait                    = "trait"
)

// RuntimeObj returns one of oam Objects matched the SType
func (s SType) RuntimeObj() runtime.Object {
	switch s {
	case STypeScope:
		return new(v1alpha1.ApplicationScope)
	case STypeTrait:
		return new(v1alpha1.Trait)
	case STypeWorkloadType:
		return new(v1alpha1.WorkloadType)
	case STypeComponent:
		return new(v1alpha1.ComponentSchematic)
	case STypeApplicationConfiguration:
		return new(v1alpha1.ApplicationConfiguration)
	default:
		panic("invalide spec type")
	}
}
