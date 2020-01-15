# Concepts

## Handler
Handler triggered by components, traits, scopes modify event.

User can implement Handler and make business logic done in handler and add action to action context, add values to action context.

## Action

Action is "resource" create action.

```
type Action struct {
	// plugin do action, e.g: k8s, helm, ...
	Provider PType
	// action command, e.g: Create, Update
	Command CmdType
	// action content, for k8s plugin, this is k8s object, for helm plugin, this is helm chart address.
	Plan interface{}
}
```

For action Action{Provider: "k8s", Command: "create", Plan: &Deployment{...}}, oam-runtime will create Deployment for you to k8s platform.


## ActionContext

ActionContext used for store actions and context values.

For actions need to be processed early, use ctx.AddPre; 

For actions need to be processed late, use ctx.AddPost;

For normal actions, just use ctx.Add.

OAM framework will do preActions -> actions -> postActions for you.
