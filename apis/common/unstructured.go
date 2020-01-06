package common

import (
	encodingjson "encoding/json"
	"fmt"
)

// Minimal unstructured type to add in Kubernetes types.
type Unstructured map[string]interface{}

func (in Unstructured) DeepCopy() Unstructured {
	if in == nil {
		return nil
	}
	out := DeepCopyJSON(in)
	return out
}

// DeepCopyJSON deep copies the passed value, assuming it is a valid JSON representation i.e. only contains
// types produced by json.Unmarshal() and also int64.
// bool, int64, float64, string, []interface{}, map[string]interface{}, json.Number and nil
func DeepCopyJSON(x map[string]interface{}) map[string]interface{} {
	return DeepCopyJSONValue(x).(map[string]interface{})
}

// DeepCopyJSONValue deep copies the passed value, assuming it is a valid JSON representation i.e. only contains
// types produced by json.Unmarshal() and also int64.
// bool, int64, float64, string, []interface{}, map[string]interface{}, json.Number and nil
func DeepCopyJSONValue(x interface{}) interface{} {
	switch x := x.(type) {
	case map[string]interface{}:
		if x == nil {
			// Typed nil - an interface{} that contains a type map[string]interface{} with a value of nil
			return x
		}
		clone := make(map[string]interface{}, len(x))
		for k, v := range x {
			clone[k] = DeepCopyJSONValue(v)
		}
		return clone
	case []interface{}:
		if x == nil {
			// Typed nil - an interface{} that contains a type []interface{} with a value of nil
			return x
		}
		clone := make([]interface{}, len(x))
		for i, v := range x {
			clone[i] = DeepCopyJSONValue(v)
		}
		return clone
	case string, int64, bool, float64, nil, encodingjson.Number:
		return x
	default:
		panic(fmt.Errorf("cannot deep copy %T", x))
	}
}
