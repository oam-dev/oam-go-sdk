package common

import (
	"encoding/json"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"

	"github.com/stretchr/testify/assert"
)

func Test_matchPattern(t *testing.T) {
	tests := []struct {
		value    string
		expMatch bool
		expKey   string
	}{
		{
			value:    "${value}",
			expMatch: true,
			expKey:   "value",
		},
		{
			value:    "${}",
			expMatch: true,
			expKey:   "",
		},
		{
			value:    "${",
			expMatch: false,
			expKey:   "",
		},
		{
			value:    "{vv}",
			expMatch: false,
			expKey:   "",
		},
		{
			value:    "${${xxx}}",
			expMatch: true,
			expKey:   "${xxx}",
		},
	}
	for _, ti := range tests {
		gotMatch, gotKey := matchPattern(ti.value)
		assert.Equal(t, ti.expMatch, gotMatch)
		assert.Equal(t, ti.expKey, gotKey)
	}
}

func Test_getParamValue(t *testing.T) {
	params := []v1alpha1.ParameterValue{
		{Name: "k1", Value: "v1"},
		{Name: "k2", Value: ""},
		{Name: "k3", Value: "xxx"},
	}
	tests := []struct {
		params []v1alpha1.ParameterValue
		key    string
		value  string
	}{
		{
			params: params,
			key:    "k1",
			value:  "v1",
		},
		{
			params: params,
			key:    "k2",
			value:  "",
		},
		{
			params: params,
			key:    "kk",
			value:  "",
		},
	}
	for _, ti := range tests {
		gotValue := getParamValue(ti.params, ti.key)
		assert.Equal(t, ti.value, gotValue)
	}
}

func Test_ExtractFromMap(t *testing.T) {
	params := []v1alpha1.ParameterValue{
		{Name: "k1", Value: "v1"},
		{Name: "k2", Value: ""},
		{Name: "k3", Value: "xxx"},
	}
	tests := []struct {
		params    []v1alpha1.ParameterValue
		values    map[string]interface{}
		expValues map[string]interface{}
	}{
		{
			params: params,
			values: map[string]interface{}{
				"k1": "${k1}",
				"k2": "${",
			},
			expValues: map[string]interface{}{
				"k1": "v1",
				"k2": "${",
			},
		},
		{
			params: params,
			values: map[string]interface{}{
				"k1": "${k1}",
				"k2": map[string]interface{}{
					"k3": "${k3}",
					"k4": 12,
				},
			},
			expValues: map[string]interface{}{
				"k1": "v1",
				"k2": map[string]interface{}{
					"k3": "xxx",
					"k4": 12,
				},
			},
		},
	}
	for _, ti := range tests {
		gotValue := ExtractFromMap(ti.params, ti.values)
		assert.Equal(t, ti.expValues, gotValue)
	}
}

func TestExtractParams(t *testing.T) {
	j1, _ := json.Marshal(map[string]interface{}{
		"k1": "${k1}",
		"k2": "${",
	})
	j2, _ := json.Marshal(map[string]interface{}{
		"k1": "${k1}",
		"k2": map[string]interface{}{
			"k3": "${k3}",
			"k4": "12",
		},
	})
	params := []v1alpha1.ParameterValue{
		{Name: "k1", Value: "v1"},
		{Name: "k2", Value: ""},
		{Name: "k3", Value: "xxx"},
	}
	tests := []struct {
		params    []v1alpha1.ParameterValue
		values    runtime.RawExtension
		expValues map[string]interface{}
	}{
		{
			params: params,
			values: runtime.RawExtension{Raw: j1},
			expValues: map[string]interface{}{
				"k1": "v1",
				"k2": "${",
			},
		},
		{
			params: params,
			values: runtime.RawExtension{Raw: j2},
			expValues: map[string]interface{}{
				"k1": "v1",
				"k2": map[string]interface{}{
					"k3": "xxx",
					"k4": "12",
				},
			},
		},
	}
	for _, ti := range tests {
		gotValue, err := ExtractParams(ti.params, ti.values)
		assert.NoError(t, err)
		assert.Equal(t, ti.expValues, gotValue)
	}
}
