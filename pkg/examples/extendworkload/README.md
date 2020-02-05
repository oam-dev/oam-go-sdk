# Component Example

## Install

1. install all crd files by  `make install`
2. `go run pkg/examples/extendworkload/main.go`

## Apply


When we install `component.yaml`, we will get such output:

```shell script
2019-12-23T16:12:31.720+0800	INFO	setup	hello oam from pre hook: comp
settings: {"Description":"xxx","Protocol":"Example","Type":"Performance"}
Description: xxx
Protocol: Example
Type: Performance
```

Raw json data is stored in `component.Spec.WorkloadSettings.Raw`. In fact you could get the real type.

You could see our workloadType in `pkg/examples/extendworkload/workloadtype.yaml`, our schema is:

```
    {
       "$schema":"http://json-schema.org/draft-07/schema#",
       "type":"object",
       "description":"",
       "required":[
          "Protocol"
       ],
       "properties":{
          "Protocol":{
             "type":"string",
             "description":""
          },
          "Type":{
             "type":"string",
             "description":""
          },
          "Description":{
             "type":"string",
             "description":""
          }
       }
    }
```

So we could use `map[string]interface{}` to parse our output, so we could get more concrete data struct.


## New CRD 

### install new CRD

```shell script
kubectl apply -f pkg/examples/extendworkload/new-crd.yaml
```

### Run demo with new CRD

```shell script
go run main.go --new-crd=true
```

## Apply app in new crd

```shell script
kubectl apply -f pkg/examples/extendworkload/app-new-crd.yaml
```