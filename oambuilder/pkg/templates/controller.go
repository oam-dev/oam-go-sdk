package templates

import (
	"strings"
)

const (
	controllerImportMark        = "+CONTROLLER_IMPORT"
	controllerTypeMark          = "+CONTROLLER_TYPE"
	controllerExchangeUsageMark = "+CONTROLLER_EXCHANGE_USAGE"
	controllerMgrSetupMark      = "+CONTROLLER_MGR_SETUP"
)

type templates string

func (t templates) replaceAll(str1, str2 string) templates {
	return templates(strings.Replace(t.String(), str1, str2, -1))
}

func (t templates) String() string {
	return string(t)
}

var controllerTemplate templates = `{{ .Boilerplate }}

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// apierr "k8s.io/apimachinery/pkg/api/errors"
	+CONTROLLER_IMPORT

	oamruntime "github.com/oam-dev/oam-go-sdk/oambuilder/pkg/runtime"
	api "{{ .Resource.GoPackage }}/{{ .Resource.Version }}"
)

// {{ .Resource.Kind }}Reconciler reconciles a {{ .Resource.Kind }} object
type {{ .Resource.Kind }}Reconciler struct {
	client.Client
	Exchanger oamruntime.+CONTROLLER_TYPEExchanger
	Log logr.Logger
	Scheme *runtime.Scheme
}

func (r *{{ .Resource.Kind }}Reconciler) exchangers() []runtime.Object {
	return []runtime.Object{
		// replace returns with your exchanges crd type
		// &v1.ExampleExchange{}
	}
}

// +kubebuilder:rbac:groups={{.Resource.GroupDomain}},resources={{ .Resource.Plural }},verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups={{.Resource.GroupDomain}},resources={{ .Resource.Plural }}/status,verbs=get;update;patch

func (r *{{ .Resource.Kind }}Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("{{ .Resource.Kind | lower }}", req.NamespacedName)

	// instance, err := r.instance(req)
	// if err != nil {
	// 	if apierr.IsNotFound(err) {
	// 		return ctrl.Result{}, nil
	// 	}
	// 	return ctrl.Result{}, errors.Wrap(err, "r.instance")
	// }
	// exList := &v1.ExampleExchangeList{}
	// _ := r.Exchanger.Resources(r.Client, instance, exList)
	+CONTROLLER_EXCHANGE_USAGE

	// your logic here

	return ctrl.Result{}, nil
}

func (r *{{ .Resource.Kind }}Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := r.initExchange(); err != nil {
		return err
	}
	 b := ctrl.NewControllerManagedBy(mgr).
		For(&api.{{ .Resource.Kind }}{})

	for _, ex := range r.exchangers() {
		b = b.+CONTROLLER_MGR_SETUP
	}
	return b.Complete(r)
}

func (r *{{ .Resource.Kind }}Reconciler) initExchange() error {
	var opt []oamruntime.ExchangerOption
	for _, ex := range r.exchangers() {
		opt = append(opt, oamruntime.WithExchanger(ex))
	}
	var err error
	r.Exchanger, err = oamruntime.New+CONTROLLER_TYPEExchanger(opt...)
	return err
}

func (r *{{ .Resource.Kind }}Reconciler) instance(req ctrl.Request) (*api.{{ .Resource.Kind }}, error) {
	i := &api.{{ .Resource.Kind }}{}
	return i, r.Get(context.TODO(), req.NamespacedName, i)
}
`

var (
	traitControllerTemplate = controllerTemplate.replaceAll(controllerImportMark, "").
				replaceAll(controllerTypeMark, "Trait").
				replaceAll(controllerMgrSetupMark, "Owns(ex)").
				replaceAll(controllerExchangeUsageMark, `
	// create exchange example
	// var ex = &v1.ExampleExchange{}
	// do something set ex.Spec
	// err := r.Exchange.Create(r.Client, instance, ex)
	// handle err
	// delete exchange example
	// for _, ex := range exList.Items {
	//    err :=  r.Exchange.Delete(r.Client, instance, ex)
	//   handle err
	// }
	// update exchange example
	// err := r.Update(context.TODO(), ex.Items[0])
	// handle err
				`)

	workloadControllerTemplate = controllerTemplate.
					replaceAll(controllerImportMark, `"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/controller-runtime/pkg/handler"`).
		replaceAll(controllerTypeMark, "Workload").
		replaceAll(controllerMgrSetupMark, `Watches(&source.Kind{
			Type: ex,
		}, &handler.EnqueueRequestForOwner{
			OwnerType:    &api.{{ .Resource.Kind }}{},
			IsController: false,
		})
	`).replaceAll(controllerExchangeUsageMark, `
	// if ready, err := r.Exchanger.AllExchangerReady(r.Client, instance); err != nil {
	// 	return ctrl.Result{}, errors.Wrap(err, "r.Exchanger.AllExchangerReady")
	// } else if !ready {
	// 	return ctrl.Result{}, nil
	// }`)
)

func TraitControllerTemplate() string {
	return traitControllerTemplate.String()
}

func WorkloadControllerTemplate() string {
	return workloadControllerTemplate.String()
}
