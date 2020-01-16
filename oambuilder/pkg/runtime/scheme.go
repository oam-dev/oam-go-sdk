package runtime

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	builders []*scheme.Builder
)

func AddToScheme(scheme *runtime.Scheme) *runtime.Scheme {
	for _, b := range builders {
		b.AddToScheme(scheme)
	}
	return scheme
}

func Register(builder *scheme.Builder) {
	builders = append(builders, builder)
}
