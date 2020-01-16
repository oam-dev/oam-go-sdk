package workload

type Trait struct {
	Name string `json:"name"`
	Init bool   `json:"init,omitempty"`
}

func (in *Trait) DeepCopyInto(out *Trait) {
	*out = *in
	return
}

func (in *Trait) DeepCopy() *Trait {
	if in == nil {
		return nil
	}
	out := new(Trait)
	in.DeepCopyInto(out)
	return out
}
