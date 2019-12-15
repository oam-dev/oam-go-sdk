package v1alpha1

import (
	"encoding/json"
	"fmt"
)

/// A value that is substituted into a parameter.
type ParameterValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

/**
Translate []ParameterValue to struct, the struct always should be empty struct pointer.
e.g: Translate(&RollOutParameter{}, p)
*/
func Translate(v interface{}, p []ParameterValue) error {
	props := make(map[string]string)
	for _, v := range p {
		props[v.Name] = v.Value
	}
	propsJsonBytes, _ := json.Marshal(props)
	// init from json bytes.
	return json.Unmarshal(propsJsonBytes, v)
}

/**
Translate struct to []ParameterValue.
*/
func TranslateReverse(v interface{}) []ParameterValue {
	propsJsonBytes, _ := json.Marshal(v)
	kvs := make(map[string]string)
	if err := json.Unmarshal(propsJsonBytes, &kvs); err != nil {
		fmt.Println("reverse translate error", err)
	}

	var pvs []ParameterValue

	for k, v := range kvs {
		pvs = append(pvs, ParameterValue{Name: k, Value: v})
	}
	return pvs
}
