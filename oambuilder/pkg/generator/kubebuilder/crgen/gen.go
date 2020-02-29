package crgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/flect"
	oamv1alpha1 "github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types/project"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	ParameterTypeMap = map[string]oamv1alpha1.ParameterType{
		"boolean": oamv1alpha1.Boolean,
		"string":  oamv1alpha1.String,
		"integer": oamv1alpha1.Number,
	}
)

type Walker func(p *project.OAMProject) ([]*File, error)

type Generator struct {
	Path   string
	OAM    string
	Output string

	schemaMap map[string]*apiext.JSONSchemaProps
}

type File struct {
	Name    string
	Content []byte
}

func (g *Generator) Run() error {
	if err := g.Validate(); err != nil {
		return err
	}
	p, err := project.LoadProject(g.OAM)
	if err != nil {
		return err
	}

	parser, err := NewParser(g.Path)
	if err != nil {
		return err
	}
	m, err := parser.Load()
	if err != nil {
		return err
	}
	g.schemaMap = m

	files, err := g.Walk(p)
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := ioutil.WriteFile(filepath.Join(g.Output, f.Name), f.Content, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) Walk(p *project.OAMProject) ([]*File, error) {
	fs := []*File{}

	for _, h := range []Walker{
		g.walkWorkloads,
		g.walkTraits,
	} {
		if wfs, err := h(p); err != nil {
			return nil, err
		} else {
			fs = append(fs, wfs...)
		}
	}

	return fs, nil
}

func (g *Generator) walkWorkloads(p *project.OAMProject) ([]*File, error) {
	fs := []*File{}
	for _, t := range p.Workloads {
		group := fmt.Sprintf("%s.%s", t.Group, p.Domain)
		lowerKind := strings.ToLower(t.Kind)
		tt := &oamv1alpha1.WorkloadType{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Workloadtype",
				APIVersion: fmt.Sprintf("%s/%s", oamv1alpha1.Group, oamv1alpha1.Version),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:              lowerKind,
				CreationTimestamp: metav1.Time{time.Now()},
			},
			Spec: oamv1alpha1.WorkloadTypeSpec{
				Names: oamv1alpha1.Names{
					Kind:     t.Kind,
					Singular: lowerKind,
					Plural:   flect.Pluralize(lowerKind),
				},
				Group:   group,
				Version: t.Version,
			},
		}
		if s, ok := g.schemaMap[tmplGVKString(group, t.Version, t.Kind)]; ok {
			settings := []oamv1alpha1.Parameter{}
			rm := map[string]interface{}{}
			for _, k := range s.Required {
				rm[k] = nil
			}
			for k, d := range s.Properties {
				_, required := rm[k]
				p := oamv1alpha1.Parameter{
					Name:        k,
					Description: d.Description,
					Required:    required,
				}
				t, ok := ParameterTypeMap[d.Type]
				if !ok {
					return nil, types.Error_InvalidParameterType
				}
				p.ParameterType = t
				settings = append(settings, p)
			}
			var settingsString string
			settingsBytes, err := json.Marshal(settings)
			if err != nil {
				return nil, err
			}
			settingsString = string(settingsBytes)
			tt.Spec.Settings = settingsString
		}
		c, err := yaml.Marshal(tt)
		if err != nil {
			return nil, err
		}
		fs = append(fs, &File{
			Name:    fmt.Sprintf("workloadtype_%s.%s.yaml", lowerKind, group),
			Content: c,
		})
	}
	return fs, nil
}

func (g *Generator) walkTraits(p *project.OAMProject) ([]*File, error) {
	fs := []*File{}
	for _, t := range p.Traits {
		group := fmt.Sprintf("%s.%s", t.Group, p.Domain)
		lowerKind := strings.ToLower(t.Kind)
		tt := &oamv1alpha1.Trait{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Trait",
				APIVersion: fmt.Sprintf("%s/%s", oamv1alpha1.Group, oamv1alpha1.Version),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:              lowerKind,
				CreationTimestamp: metav1.Time{time.Now()},
			},
			Spec: oamv1alpha1.TraitSpec{
				Group:   group,
				Version: t.Version,
				Names: oamv1alpha1.Names{
					Kind:     t.Kind,
					Singular: lowerKind,
					Plural:   lowerKind + "s",
				},
				// todo appliesTo
				AppliesTo: []string{},
			},
		}
		if s, ok := g.schemaMap[tmplGVKString(group, t.Version, t.Kind)]; ok {
			b, err := json.Marshal(s)
			if err != nil {
				return nil, err
			}
			tt.Spec.Properties = string(b)
		}
		c, err := yaml.Marshal(tt)
		if err != nil {
			return nil, err
		}
		fs = append(fs, &File{
			Name:    fmt.Sprintf("trait_%s.%s.yaml", lowerKind, group),
			Content: c,
		})
	}
	return fs, nil
}

func (g *Generator) loadPackage() {
}

func (g *Generator) Validate() error {
	if g.OAM == "" {
		return types.Error_NeedOAM
	}
	if g.Path == "" {
		return types.Error_NeedPath
	}
	if g.Output == "" {
		return types.Error_NeedOutput
	}
	return nil
}
