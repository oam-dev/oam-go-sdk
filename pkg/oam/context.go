package oam

type ActionContext struct {
	PreActions  []Action
	Actions     []Action
	PostActions []Action
	Values      map[string]interface{}
}

// add actions executed before actions added through Add method
func (o *ActionContext) AddPre(a ...Action) {
	o.PreActions = append(o.PreActions, a...)
}

// add actions executed before actions added through AddPost method
func (o *ActionContext) Add(a ...Action) {
	o.Actions = append(o.Actions, a...)
}

func (o *ActionContext) AddPost(a ...Action) {
	o.PostActions = append(o.PostActions, a...)
}

func (o *ActionContext) AddValue(k string, v interface{}) {
	if o.Values == nil {
		o.Values = map[string]interface{}{}
	}
	o.Values[k] = v
}

func (o *ActionContext) GetValue(k string) interface{} {
	if o.Values == nil {
		return nil
	}
	return o.Values[k]
}

// clear and gather all actions according to action order.
func (o *ActionContext) Gather() []Action {
	var actions []Action
	actions = append(actions, o.PreActions...)
	o.PreActions = nil
	actions = append(actions, o.Actions...)
	o.Actions = nil
	actions = append(actions, o.PostActions...)
	o.PostActions = nil
	return actions
}
