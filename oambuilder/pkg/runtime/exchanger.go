package runtime

import (
	"context"
	"encoding/json"

	"github.com/go-logr/logr"
	"github.com/juju/errors"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ExchangerOption func(e *Exchanger) error

type Exchanger struct {
	exchangers map[schema.GroupVersionKind]runtime.Object
	log        logr.Logger
}

// Resources sets list resource which realated to self resource
func (e *Exchanger) Resources(cli client.Client, self runtime.Object, list runtime.Object) error {
	var gvk = list.GetObjectKind().GroupVersionKind()
	if _, ok := e.exchangers[gvk]; !ok {
		return errors.Errorf("unknown list resource %v", gvk)
	}
	getter := self.(types.WorkloadUIDGetter)
	objMeta, err := meta.Accessor(self)
	if err != nil {
		return err
	}
	labels := objMeta.GetLabels()
	return cli.List(context.TODO(),
		list,
		client.InNamespace(objMeta.GetNamespace()),
		&client.MatchingLabels{
			types.LABEL_OAM_UUID:     labels[types.LABEL_OAM_UUID],
			types.LABEL_OAM_WORKLOAD: string(getter.WorkloadUID()),
		})
}

// CreateExchangeResource automatically sets owner reference,
// finalizer, label of workload  into exchange resource,
// then create the exchange resource
// note: this function MUST be called by tait
func (e *Exchanger) Create(cli client.Client, traitResource runtime.Object, exchangeResource runtime.Object) error {
	var t dynamicTrait
	err := parse(traitResource, &t)
	if err != nil {
		return err
	}
	exObjMeta, err := meta.Accessor(exchangeResource)
	if err != nil {
		return err
	}
	var has bool
	for _, o := range exObjMeta.GetOwnerReferences() {
		if o.Kind == t.Spec.Workload.Kind &&
			o.APIVersion == t.Spec.Workload.APIVersion &&
			o.Name == t.Spec.Workload.Name {
			e.log.Info("create exchange resource %v find workload %v owner reference setted", exObjMeta.GetName(), t.Spec.Workload.Name)
			has = true
			break
		}
	}
	if !has {
		exObjMeta.SetOwnerReferences(append(exObjMeta.GetOwnerReferences(), t.Spec.Workload.OwnerReference))
	}
	exObjMeta.SetFinalizers(addFinalizer(exObjMeta.GetFinalizers(), t.Spec.Workload.Name))
	exObjMeta.SetLabels(addLabel(exObjMeta.GetLabels(), types.LABEL_OAM_WORKLOAD, string(t.Spec.Workload.UID)))
	has = false

	return cli.Create(context.TODO(), exchangeResource)
}

// Delete exchange resource
func (e *Exchanger) Delete(cli client.Client, traitResource runtime.Object, exchangeResource runtime.Object) error {
	var t dynamicTrait
	err := parse(traitResource, &t)
	if err != nil {
		return err
	}
	exObjMeta, err := meta.Accessor(exchangeResource)
	if err != nil {
		return err
	}
	exObjMeta.SetFinalizers(delFinalizer(exObjMeta.GetFinalizers(), t.Spec.Workload.Name))
	return cli.Delete(context.TODO(), exchangeResource)
}

func addLabel(r map[string]string, k, v string) map[string]string {
	if r == nil {
		return map[string]string{k: v}
	}
	r[k] = v
	return r
}

func delFinalizer(f []string, str string) []string {
	for i, v := range f {
		if v == str {
			return append(f[:i], f[i+1:]...)
		}
	}
	return f
}

func addFinalizer(f []string, str string) []string {
	for _, v := range f {
		if v == str {
			return f
		}
	}
	return append(f, str)
}

// AllExchangerReady returns all related exchange crds, by given workload resource, are created or not
func (e *Exchanger) AllExchangerReady(cli client.Client, self runtime.Object) (bool, error) {
	var d dynamicWorkload
	err := parse(self, &d)
	if err != nil {
		return false, errors.Annotate(err, "parse related traits")
	}
	var dts = make(map[string]struct{})
	for _, obj := range e.exchangers {
		var dpCp = obj.DeepCopyObject()
		err := e.Resources(cli, self, dpCp)
		if err != nil {
			return false, err
		}
		var dl dynamicTraitList
		err = parse(dpCp, &dl)
		if err != nil {
			return false, err
		}
		for _, d := range dl.Items {
			dts[d.Name] = struct{}{}
		}
	}
	for _, t := range d.Spec.Traits {
		_, ok := dts[t.Name]
		if t.Init && !ok {
			return false, nil
		}
	}
	return true, nil
}

// WithExchanger returns a ExchangerOption which can register given obj to Exchanger
func WithExchanger(obj runtime.Object) ExchangerOption {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return func(e *Exchanger) error {
		if _, ok := e.exchangers[gvk]; ok {
			return nil
		}
		e.exchangers[gvk] = obj
		return nil
	}
}

// WithLogger is a option for set Exchanger logger
func WithLogger(l logr.Logger) ExchangerOption {
	return func(e *Exchanger) error {
		e.log = l
		return nil
	}
}

func newExchanger(opts ...ExchangerOption) (*Exchanger, error) {
	var e = &Exchanger{
		exchangers: make(map[schema.GroupVersionKind]runtime.Object),
		log:        ctrl.Log.WithName("exchanger"),
	}
	for _, o := range opts {
		err := o(e)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

// NewTraitExchanger return TraitExchanger for trait controller
func NewTraitExchanger(opts ...ExchangerOption) (TraitExchanger, error) {
	return newExchanger(opts...)
}

// NewWorkloadExchanger return rWorkloadExchanger for trait controller
func NewWorkloadExchanger(opts ...ExchangerOption) (WorkloadExchanger, error) {
	return newExchanger(opts...)
}

func parse(o runtime.Object, i interface{}) error {
	bts, err := json.Marshal(o)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bts, i)
	return err
}
