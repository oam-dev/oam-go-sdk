package util

import (
	"encoding/json"
	"reflect"

	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	OAMObjectCurrentSpecAnnotationKey = v1alpha1.Group + v1alpha1.Separator + "oam-object-current-spec"
)

// SpecEqual compare object (such as *deployment, need to be pointer) and spec (such as DeploymentSpec)
// If flag update is set, the input object will change annotation if current spec is not equal to target spec.
// This function assumes object has field `ObjectMeta` which must be satisfied to use it.
// This function is designed to reduce traffic in EDAS reconcile. Other controller (such as rollout-controller)
// should not make assumption that object with same spec would not be updated repetitively although this function
// should filter these situations.
func SpecEqual(object runtime.Object, curSpec interface{}, update bool) bool {
	meta := reflect.ValueOf(object).Elem().FieldByName("ObjectMeta").Interface().(metav1.ObjectMeta)
	lastSpecLiteral := ""
	if meta.Annotations != nil {
		lastSpecLiteral = meta.Annotations[OAMObjectCurrentSpecAnnotationKey]
	}
	lastSpec := reflect.New(reflect.TypeOf(curSpec)).Interface()
	currentSpec := reflect.New(reflect.TypeOf(curSpec)).Interface()
	if bytes, err := json.Marshal(curSpec); err != nil {
		return false
	} else if err := json.Unmarshal(bytes, currentSpec); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(lastSpecLiteral), lastSpec); err != nil || !reflect.DeepEqual(lastSpec, currentSpec) {
		if update {
			bytes, err := json.Marshal(curSpec)
			if err != nil {
				bytes = []byte{}
			}
			meta := reflect.ValueOf(object).Elem().FieldByName("ObjectMeta").Addr().Interface().(*metav1.ObjectMeta)
			if meta.Annotations == nil {
				meta.Annotations = make(map[string]string)
			}
			meta.Annotations[OAMObjectCurrentSpecAnnotationKey] = string(bytes)
		}
		return false
	}
	return true
}
