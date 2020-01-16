package generator

import (
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/generator/kubebuilder"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types/project"
)

type Generator interface {
	Detect(path string) (bool, error)
	Execute(args []string) error
	AttachTemplate(tmpl string, ty types.TemplateType) error
	PartialProject() (*project.PartialProject, error)
}

type generatorEngine struct {
	handlers []Generator
}

var engine = &generatorEngine{
	handlers: []Generator{
		kubebuilder.Builder(),
	},
}

func Engine() *generatorEngine {
	return engine
}

func (e *generatorEngine) Generator(path string) (g Generator, err error) {
	var ok bool
	for _, h := range e.handlers {
		if ok, err = h.Detect(path); ok && err == nil {
			g = h
			break
		}
	}
	if g == nil {
		err = types.Error_FailedToFindGenerator
	}
	return g, err
}
