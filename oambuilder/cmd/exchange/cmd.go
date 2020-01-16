package exchange

import (
	"os"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"

	ccmd "github.com/oam-dev/oam-go-sdk/oambuilder/cmd"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types"
)

var (
	logger = ctrl.Log.WithName("oambuilder.exchange")
)

var Exchange = &cobra.Command{
	Use:                "exchange",
	Short:              "generate echange crd for oam trait & workload.",
	Long:               "generate echange crd for oam trait & workload.",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, err := ccmd.Getwd()
		if err != nil {
			logger.Error(err, "cmd.Getwd")
		}
		g, err := ctx.GetGenerator()
		if err != nil {
			logger.Error(err, "ctx.GetGenerator")
			os.Exit(1)
		}

		g.AttachTemplate(typeTemplate, types.TemplateType_Type)
		// exchange crd without controller
		args = append(args, "--controller=false")

		if err := g.Execute(args); err != nil {
			logger.Error(err, "g.Execute")
			os.Exit(1)
		}

		if err := ctx.UpdateProject(g, types.ResourceType_Exchange); err != nil {
			logger.Error(err, "ctx.UpdateProject")
			os.Exit(1)
		}
	},
}

var typeTemplate = `{{ .Boilerplate }}

package {{ .Resource.Version }}

import (
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// {{.Resource.Kind}}Spec defines the desired state which set by trait
type {{.Resource.Kind}}Spec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state shown for workload
	// Important: Run "make" to regenerate code after modifying this file
}

// {{.Resource.Kind}}Status defines the observed state which set by workload
type {{.Resource.Kind}}Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state shown for trait
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
{{ if not .Resource.Namespaced }} // +kubebuilder:resource:scope=Cluster {{ end }}
// {{.Resource.Kind}} is the Schema for the {{ .Resource.Resource }} API
type {{.Resource.Kind}} struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
	metav1.ObjectMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
	Spec   {{.Resource.Kind}}Spec   ` + "`" + `json:"spec,omitempty"` + "`" + `
	Status {{.Resource.Kind}}Status ` + "`" + `json:"status,omitempty"` + "`" + `
}

// +kubebuilder:object:root=true
// {{.Resource.Kind}}List contains a list of {{.Resource.Kind}}
type {{.Resource.Kind}}List struct {
	metav1.TypeMeta ` + "`" + `json:",inline"` + "`" + `
	metav1.ListMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
	Items           []{{ .Resource.Kind }} ` + "`" + `json:"items"` + "`" + `
}

func init() {
	SchemeBuilder.Register(&{{.Resource.Kind}}{}, &{{.Resource.Kind}}List{})
	oambuilder := &scheme.Builder{GroupVersion: GroupVersion}
	oambuilder.Register(&{{.Resource.Kind}}{}, &{{.Resource.Kind}}List{})
	runtime.Register(oambuilder)
}
`
