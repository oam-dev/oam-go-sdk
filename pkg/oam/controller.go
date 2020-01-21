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

package oam

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/oam-dev/oam-go-sdk/pkg/config"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// Reconciler reconciles a runtime object in oam
type Reconciler struct {
	client.Client
	specType          SType
	Log               logr.Logger
	Scheme            *runtime.Scheme
	ControllerContext ControllerContext
}

// +kubebuilder:rbac:groups=*,resources=*,verbs=*
func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var name = r.specType
	ctx := context.Background()
	log := r.Log.WithValues(string(name), req.NamespacedName)
	actionCtx := &ActionContext{}

	var conf = name.RuntimeObj()

	opCode, err := r.getOpCode(ctx, req.NamespacedName, conf)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			// ignore not found error
			return ctrl.Result{}, nil
		}
		log.Error(err, "get operate code error")
		return ctrl.Result{}, err
	}

	// get operation code
	eType := CreateOrUpdate
	if opCode == config.DeleteOpCode {
		eType = Delete
	}

	// invoke handler
	handlers := getHandlers(name)
	// ApplicationConfiguration contains Components fileds, how applicationConfiguration
	// works with Components depends on implementor
	// for _, compConf := range conf.Spec.Components {
	// if handlers != nil {
	// for _, h := range handlers {
	// if err := h.Handle(actionCtx, &conf, &compConf, eType); err != nil {
	// log.Error(err, "handler handle error", "handler id", h.Id())
	// return ctrl.Result{}, err
	// }
	// }
	// }
	// }
	for _, h := range handlers {
		if err := h.Handle(actionCtx, conf, eType); err != nil {
			log.Error(err, "handler handle error", "handler id", h.Id())
			return ctrl.Result{}, err
		}
	}

	// do handler related actions
	if err := r.doActions(actionCtx, log); err != nil {
		log.Error(err, "do handler related actions error")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) doActions(actionCtx *ActionContext, log logr.Logger) error {
	actions := actionCtx.Gather()
	for _, action := range actions {
		if action.Provider != PTypeK8S {
			// todo: process this panic.
			panic("not support action provider:" + action.Provider)
		}
		robj := action.Plan.(runtime.Object)
		switch action.Command {
		case CmdTypeCreate:
			if err := r.Create(context.Background(), robj); err != nil {
				log.Error(err, "do create action error", "provider", "k8s", "plan", robj)
				return err
			}
		case CmdTypeUpdate:
			if err := r.Update(context.Background(), robj); err != nil {
				log.Error(err, "do update action error", "provider", "k8s", "plan", robj)
				return err
			}
		case CmdTypeDelete:
			if err := r.Delete(context.Background(), robj); err != nil {
				log.Error(err, "do delete action error", "provider", "k8s", "plan", robj)
				return err
			}
		}
	}
	return nil
}

func (r *Reconciler) getOpCode(ctx context.Context,
	name types.NamespacedName, conf runtime.Object) (opCode int, err error) {
	opCode = config.CreateOrUpdateOpCode
	if err := r.Get(ctx, name, conf); err != nil {
		// return with error
		return opCode, err
	}

	obj, err := meta.Accessor(conf)
	if err != nil {
		return opCode, err
	}

	// deletionTS is not zero, this object is deleted
	if !obj.GetDeletionTimestamp().IsZero() {
		opCode = config.DeleteOpCode
	}

	return opCode, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	obj := r.specType.RuntimeObj()
	bld := ctrl.NewControllerManagedBy(mgr).For(obj)

	owns := getOwns(r.specType)
	if owns != nil {
		for _, o := range owns {
			bld = bld.Owns(o)
		}
	}

	return bld.Complete(r)
}
