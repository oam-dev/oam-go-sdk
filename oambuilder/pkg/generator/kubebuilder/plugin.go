package kubebuilder

import (
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types"
	"sigs.k8s.io/kubebuilder/pkg/model"
)

type Plugin struct {
	tmpls map[types.TemplateType]*tmpl
}

func (p *Plugin) Pipe(u *model.Universe) error {
	for _, t := range p.tmpls {
		if err := t.Pipe(u); err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) Attach(ty types.TemplateType, t *tmpl) {
	p.tmpls[ty] = t
}
