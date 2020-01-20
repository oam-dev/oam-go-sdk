package common

import (
	"encoding/json"
	"strings"

	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

func matchPattern(value string) (match bool, key string) {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		key = strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
		return true, key
	}
	return false, ""
}

func getParamValue(params []v1alpha1.ParameterValue, key string) string {
	for _, v := range params {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}

func ExtractFromMap(params []v1alpha1.ParameterValue, values map[string]interface{}) map[string]interface{} {
	for k, v := range values {
		switch subVal := v.(type) {
		case string:
			match, key := matchPattern(subVal)
			if match {
				values[k] = getParamValue(params, key)
			}
		case map[string]interface{}:
			values[k] = ExtractFromMap(params, subVal)
		}
	}
	return values
}

// ExtractParams will extract param from Pattern "${parameter_key}"
func ExtractParams(params []v1alpha1.ParameterValue, raw runtime.RawExtension) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	if err := json.Unmarshal(raw.Raw, &values); err != nil {
		return nil, err
	}
	return ExtractFromMap(params, values), nil
}
