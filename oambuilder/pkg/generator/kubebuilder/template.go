package kubebuilder

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types"
	"sigs.k8s.io/kubebuilder/pkg/model"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
	"sigs.k8s.io/kubebuilder/plugins/addon"
)

var (
	ControllerFileName = func(u *model.Universe) string {
		return filepath.Join("controllers", strings.ToLower(u.Resource.Kind)+"_controller.go")
	}
	TypeFileName = func(u *model.Universe) string {
		return filepath.Join("api", u.Resource.Version, strings.ToLower(u.Resource.Kind)+"_types.go")
	}
	MainFileName = func(u *model.Universe) string {
		return "main.go"
	}
	GVInfoFileName = func(u *model.Universe) string {
		return filepath.Join("api", u.Resource.Version, "groupversion_info.go")
	}
	MakefileFileName = func(u *model.Universe) string {
		return "Makefile"
	}
)

func Template(t string, ty types.TemplateType) (*tmpl, error) {
	tt := &tmpl{
		template: t,
		Type:     ty,
		FuncMap:  addon.DefaultTemplateFunctions(),
	}
	switch ty {
	case types.TemplateType_GVInfo:
		tt.fileName = GVInfoFileName
	case types.TemplateType_Controller:
		tt.fileName = ControllerFileName
	case types.TemplateType_Type:
		tt.fileName = TypeFileName
		tt.FuncMap["JSONTag"] = addon.JSONTag
	case types.TemplateType_Main:
		tt.fileName = MainFileName
	case types.TemplateType_Makefile:
		tt.fileName = MakefileFileName
	}

	return tt, nil
}

type SetFileName func(u *model.Universe) string

type tmpl struct {
	Type    types.TemplateType
	FuncMap template.FuncMap

	template string
	fileName SetFileName
}

func (t *tmpl) Pipe(u *model.Universe) error {
	templateBody := t.template

	funcs := t.FuncMap
	contents, err := addon.RunTemplate(t.Type.String(), templateBody, u, funcs)
	if err != nil {
		return err
	}

	m := &model.File{
		Path:           t.fileName(u),
		Contents:       contents,
		IfExistsAction: input.Overwrite,
	}

	ok, err := addon.AddFile(u, m)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	return addon.ReplaceFile(u, m)
}
