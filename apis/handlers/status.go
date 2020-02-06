package handlers

import (
	"fmt"
	"sync"

	"github.com/oam-dev/oam-go-sdk/apis/flags"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type StatusHander func(rsrc metav1.Object) string

var statusHandlers = make(map[string]StatusHander)
var statusHandlerLock sync.Mutex

func FormatGVK(gvk schema.GroupVersionKind) string {
	return fmt.Sprintf("%s/%s.%s", gvk.Group, gvk.Version, gvk.Kind)
}

func RegisterStatusHandler(gvk schema.GroupVersionKind, handler StatusHander) {
	statusHandlerLock.Lock()
	defer statusHandlerLock.Unlock()
	statusHandlers[FormatGVK(gvk)] = handler
}

func TryStatusHandler(r metav1.Object) string {
	if ro, ok := r.(runtime.Object); ok {
		statusHandlerLock.Lock()
		handler, ok := statusHandlers[FormatGVK(ro.GetObjectKind().GroupVersionKind())]
		statusHandlerLock.Unlock()
		if ok {
			return handler(r)
		}
	}
	return flags.StatusUnknown
}
