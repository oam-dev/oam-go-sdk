package v1alpha1

import (
	"encoding/json"
	"log"
	"testing"
)

var a = `
{
	"name": "xxx",
	"name1": ["1", "2"]
}

`

func TestInter(t *testing.T) {
	var aa = map[string]InterfaceStruct{}
	json.Unmarshal([]byte(a), &aa)
	bts, _ := json.Marshal(&aa)
	log.Print(string(bts))

}
