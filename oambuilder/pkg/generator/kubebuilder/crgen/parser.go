package crgen

import (
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-tools/pkg/crd"
	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type Parser struct {
	*crd.Parser

	Roots []*loader.Package
	Path  string
}

func NewParser(p string) (*Parser, error) {
	roots, err := loader.LoadRoots(p)
	if err != nil {
		return nil, err
	}
	reg := &markers.Registry{}
	if err := crdmarkers.Register(reg); err != nil {
		return nil, err
	}
	return &Parser{
		Parser: &crd.Parser{
			Collector: &markers.Collector{
				Registry: reg,
			},
			Checker: &loader.TypeChecker{},
		},
		Roots: roots,
		Path:  p,
	}, nil
}

func (p *Parser) init() {
	crd.AddKnownTypes(p.Parser)
	for _, root := range p.Roots {
		p.NeedPackage(root)
	}
}

func (p *Parser) Load() (map[string]*apiext.JSONSchemaProps, error) {
	p.init()

	metav1Pkg := crd.FindMetav1(p.Roots)
	if metav1Pkg == nil {
		return nil, nil
	}
	kubeKinds := crd.FindKubeKinds(p.Parser, metav1Pkg)
	if len(kubeKinds) == 0 {
		return nil, nil
	}

	m := make(map[string]*apiext.JSONSchemaProps)
	for _, gk := range kubeKinds {
		mm, err := p.loadCRDSettings(gk)
		if err != nil {
			return nil, err
		}
		for k, v := range mm {
			m[k] = v
		}
	}
	return m, nil
}

func (p *Parser) loadCRDSettings(gk schema.GroupKind) (map[string]*apiext.JSONSchemaProps, error) {
	m := make(map[string]*apiext.JSONSchemaProps)
	for pkg, gv := range p.GroupVersions {
		if gv.Group != gk.Group {
			continue
		}

		ident := crd.TypeIdent{Package: pkg, Name: gk.Kind}
		info := p.Types[ident]
		if info == nil {
			continue
		}

		p.NeedFlattenedSchemaFor(ident)
		fs := p.FlattenedSchemata[ident]
		settings := fs.DeepCopy().Properties["spec"].Properties["settings"]
		m[tmplGVKString(gk.Group, p.GroupVersions[pkg].Version, gk.Kind)] = &settings
	}
	return m, nil
}

func tmplGVKString(g, v, k string) string {
	return g + ";" + v + ";" + k
}
