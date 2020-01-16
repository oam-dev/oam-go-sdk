package v1alpha1

type InterfaceParam struct {
	Name string           `json:"name"`
	Val  *InterfaceStruct `json:"value"`
}

type InterfaceStruct struct {
	val []byte `json:"-"`
}

func (inter *InterfaceStruct) Real() []byte {
	return inter.val
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (inter *InterfaceStruct) UnmarshalJSON(value []byte) error {
	inter.val = value
	return nil
}

// MarshalJSON implements the json.Marshaller interface.
func (inter *InterfaceStruct) MarshalJSON() ([]byte, error) {
	return inter.val, nil
}

func (InterfaceStruct) OpenAPISchemaType() []string {
	return []string{"string", "array", "object", "int"}
}
func (InterfaceStruct) OpenAPISchemaFormat() []string {
	return []string{"interface-struct"}
}
