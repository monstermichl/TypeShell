package parser

import (
	"fmt"
)

// TODO: Just keep this for now to be able to compile but remove later.
type Variable struct {
	ImportableBase
	valueType ValueType
	layer     int
}

func NewVariable(name string, prefix string, valueType ValueType, layer int) Variable {
	return Variable{
		ImportableBase: NewImportableBase(name, prefix, layer == 0),
		valueType:      valueType,
		layer:          layer,
	}
}

func (v Variable) ValueType() ValueType {
	return v.valueType
}

func (v Variable) Layer() int {
	return v.layer
}

func (v Variable) LayerName() string {
	return fmt.Sprintf("%s_%d", v.name, v.layer)
}

func (v Variable) IsConstant() bool {
	return false
}

func (v *Variable) SetValueType(valueType ValueType) {
	v.valueType = valueType
}
