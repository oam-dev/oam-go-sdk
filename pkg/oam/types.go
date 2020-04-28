package oam

import (
	"fmt"
	"sync"

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

var stypes = make(map[SType]runtime.Object)
var typeLock sync.Mutex

func init() {
	stypes[STypeScope] = new(v1alpha1.ApplicationScope)
	stypes[STypeTrait] = new(v1alpha1.Trait)
	stypes[STypeWorkloadType] = new(v1alpha1.WorkloadType)
	stypes[STypeComponent] = new(v1alpha1.ComponentSchematic)
	stypes[STypeApplicationConfiguration] = new(v1alpha1.ApplicationConfiguration)
}

func RegisterObject(tp SType, obj runtime.Object) {
	typeLock.Lock()
	defer typeLock.Unlock()
	stypes[tp] = obj
}

// RuntimeObj returns one of oam Objects matched the SType
func (s SType) RuntimeObj() runtime.Object {
	obj, err := s.GetRuntimeObj()
	if err == nil {
		return obj
	}
	panic("invalid spec type")
}

func (s SType) GetRuntimeObj() (runtime.Object, error) {
	typeLock.Lock()
	obj, ok := stypes[s]
	typeLock.Unlock()
	if ok {
		return obj.DeepCopyObject(), nil
	}
	return nil, fmt.Errorf("can't get spec type '%s', please register it", s)
}
