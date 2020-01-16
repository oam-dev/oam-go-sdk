package project

import (
	"io/ioutil"
	"os"

	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/yaml"
)

var (
	logger = ctrl.Log.WithName("oambuilder.project")
)

const (
	PROJECT        = "OAM"
	Version1Alpha1 = "v1alpha1"
)

type DomainRepo struct {
	Domain string `yaml"domain" json:"domain"`
	Repo   string `yaml:"repo" json:"repo"`
}

func (dr DomainRepo) IsSame(ddr DomainRepo) bool {
	return dr.Domain == ddr.Domain && dr.Repo == ddr.Repo
}

type OAMProject struct {
	Version string `yaml:"version" json:"version"`

	DomainRepo `yaml:",inline" json:",inline"`
	Workloads  []types.GroupVersionKind `yaml:"workloads,omitempty" json:"workloads,omitempty"`
	Traits     []types.GroupVersionKind `yaml:"traits,omitempty" json:"traits,omitempty"`
	Exchanges  []types.GroupVersionKind `yaml:"exchange,omitempty" json:"exchange,omitempty"`
}

type PartialProject struct {
	DomainRepo             `yaml:",inline" json:",inline"`
	types.GroupVersionKind `yaml:",inline" json:",inline"`
}

func UpdateProject(p string, pp *PartialProject, ty types.ResourceType) error {
	project, err := LoadProject(p)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		project = &OAMProject{
			DomainRepo: pp.DomainRepo,
			Version:    Version1Alpha1,
		}
		a := []types.GroupVersionKind{pp.GroupVersionKind}
		switch ty {
		case types.ResourceType_Exchange:
			project.Exchanges = a
		case types.ResourceType_Workload:
			project.Workloads = a
		case types.ResourceType_Trait:
			project.Traits = a
		}
		return CreateProject(project, p)
	}

	if isSame := project.DomainRepo.IsSame(pp.DomainRepo); !isSame {
		return types.Error_VDRMismatch
	}
	var a *[]types.GroupVersionKind
	switch ty {
	case types.ResourceType_Exchange:
		if project.Exchanges == nil {
			project.Exchanges = []types.GroupVersionKind{}
		}
		a = &project.Exchanges
	case types.ResourceType_Workload:
		if project.Workloads == nil {
			project.Workloads = []types.GroupVersionKind{}
		}
		a = &project.Workloads
	case types.ResourceType_Trait:
		if project.Traits == nil {
			project.Traits = []types.GroupVersionKind{}
		}
		a = &project.Traits
	}

	var dup bool
	for _, gvk := range *a {
		if same := gvk.IsSame(pp.GroupVersionKind); same {
			dup = true
			break
		}
	}
	if !dup {
		*a = append(*a, pp.GroupVersionKind)
	}

	return FlushProject(project, p)
}

func LoadProject(p string) (*OAMProject, error) {
	bts, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	project := &OAMProject{}
	if err := yaml.Unmarshal(bts, project); err != nil {
		return nil, err
	}
	return project, nil
}

func CreateProject(project *OAMProject, p string) error {
	return writeProject(project, p, os.O_CREATE|os.O_WRONLY)
}

func FlushProject(project *OAMProject, p string) error {
	return writeProject(project, p, os.O_RDWR|os.O_TRUNC)
}

func writeProject(project *OAMProject, p string, flag int) error {
	bts, err := yaml.Marshal(project)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p, flag, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	writed, err := f.Write(bts)
	if err != nil {
		return err
	} else if writed < len(bts) {
		return types.Error_PartialWrite
	}
	return nil
}
