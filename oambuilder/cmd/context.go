package cmd

import (
	"os"
	"path/filepath"

	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/generator"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types/project"
)

type Context struct {
	Wd string
}

func Getwd() (*Context, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Context{
		Wd: wd,
	}, nil
}

func (ctx *Context) GetGenerator() (generator.Generator, error) {
	g, err := generator.Engine().Generator(ctx.Wd)
	return g, err
}

func (ctx *Context) UpdateProject(g generator.Generator, ty types.ResourceType) error {
	pp, err := g.PartialProject()
	if err != nil {
		return err
	}
	return project.UpdateProject(filepath.Join(ctx.Wd, project.PROJECT), pp, ty)
}
