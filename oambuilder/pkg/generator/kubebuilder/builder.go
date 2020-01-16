package kubebuilder

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	flag "github.com/spf13/pflag"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types"
	bp "github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types/project"
	oamutil "github.com/oam-dev/oam-go-sdk/oambuilder/pkg/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/kubebuilder/cmd/util"
	"sigs.k8s.io/kubebuilder/pkg/scaffold"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/project"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/resource"
)

var (
	logger = ctrl.Log.WithName("oambuilder.builder")
)

type builder struct {
	scaffold                     *scaffold.API
	plugin                       *Plugin
	flags                        *flag.FlagSet
	resourceFlag, controllerFlag *flag.Flag
	project                      *bp.PartialProject

	runMake bool

	executed bool
}

func Builder() *builder {
	p := &Plugin{
		tmpls: make(map[types.TemplateType]*tmpl),
	}
	b := &builder{
		scaffold: &scaffold.API{
			Plugins: []scaffold.Plugin{p},
		},
		plugin:  p,
		flags:   flag.NewFlagSet("kubebuilder", flag.ContinueOnError),
		project: &bp.PartialProject{},
	}
	b.setupFlags()
	return b
}

// 这里有大段的复制自kubebuilder的代码，这里这么做的原因是kubebuilder目前尚未将plugin功能实现，仅在scaffold层加了一个interface以供拓展，
// 这里要实现kububuilder style的生成逻辑，暂时无法通过外挂plugin command实现，只能使用scaffold package。
// 直到kubebuilder支持plugin。
func (b *builder) setupFlags() {
	b.flags.BoolVar(&b.runMake, "make", true,
		"if true, run make after generating files")
	b.flags.BoolVar(&b.scaffold.DoResource, "resource", true,
		"if set, generate the resource without prompting the user")
	b.resourceFlag = b.flags.Lookup("resource")
	b.flags.BoolVar(&b.scaffold.DoController, "controller", true,
		"if set, generate the controller without prompting the user")
	b.controllerFlag = b.flags.Lookup("controller")
	b.flags.BoolVar(&b.scaffold.Force, "force", false,
		"attempt to create resource even if it already exists")
	r := &resource.Resource{}
	b.flags.StringVar(&r.Kind, "kind", "", "resource Kind")
	b.flags.StringVar(&r.Group, "group", "", "resource Group")
	b.flags.StringVar(&r.Version, "version", "", "resource Version")
	b.flags.BoolVar(&r.Namespaced, "namespaced", true, "resource is namespaced")
	b.flags.BoolVar(&r.CreateExampleReconcileBody, "example", true,
		"if true an example reconcile body should be written while scaffolding a resource.")

	b.scaffold.Resource = r
}

func (b *builder) Detect(p string) (bool, error) {
	pp, err := scaffold.LoadProjectFile(path.Join(p, "PROJECT"))
	if err != nil {
		return false, nil
	}
	if pp.Version != project.Version2 {
		return false, types.Error_OnlySupportScaffoldV2
	}
	logger.Info("kubebuilder project detected.")
	b.project.DomainRepo = bp.DomainRepo{
		Domain: pp.Domain,
		Repo:   pp.Repo,
	}
	return true, nil
}

func (b *builder) Execute(args []string) error {
	if err := b.flags.Parse(args); err != nil {
		return err
	}

	if err := b.scaffold.Validate(); err != nil {
		log.Fatalln(err)
	}

	reader := bufio.NewReader(os.Stdin)
	if !b.resourceFlag.Changed {
		fmt.Println("Create Resource [y/n]")
		b.scaffold.DoResource = util.Yesno(reader)
	}

	if !b.controllerFlag.Changed {
		fmt.Println("Create Controller [y/n]")
		b.scaffold.DoController = util.Yesno(reader)
	}

	logger.Info("Writing scaffold for you to edit...")

	if err := b.scaffold.Scaffold(); err != nil {
		return err
	}

	if err := b.postScaffold(); err != nil {
		return err
	}

	var adds = map[string][]string{
		oamutil.ApiSchemeScaffoldMarker: []string{
			"oamruntime.AddToScheme(scheme)",
		},
		oamutil.ApiPkgImportScaffoldMarker: []string{
			`oamruntime "github.com/oam-dev/oam-go-sdk/oambuilder/pkg/runtime"`,
		},
	}

	if err := oamutil.InsertStringsInFile(path.Join("main.go"), adds); err != nil {
		return err
	}
	b.project.GroupVersionKind = types.GroupVersionKind{
		Group:   b.scaffold.Resource.Group,
		Version: b.scaffold.Resource.Version,
		Kind:    b.scaffold.Resource.Kind,
	}
	b.executed = true
	return nil
}

func (b *builder) postScaffold() error {
	if b.runMake {
		logger.Info("Running make...")
		cm := exec.Command("make") // #nosec
		cm.Stderr = os.Stderr
		cm.Stdout = os.Stdout
		if err := cm.Run(); err != nil {
			return fmt.Errorf("error running make: %v", err)
		}
	}
	return nil
}

func (b *builder) AttachTemplate(t string, ty types.TemplateType) error {
	tt, err := Template(t, ty)
	if err != nil {
		return err
	}
	b.plugin.Attach(ty, tt)
	return nil
}

func (b *builder) PartialProject() (*bp.PartialProject, error) {
	if !b.executed {
		return nil, types.Error_NeedExecuteFirst
	}
	return b.project, nil
}
