package parser

import "fmt"

type Const struct {
	ImportableBase
	valueType ValueType
	layer     int
}

func NewConst(name string, prefix string, valueType ValueType, layer int) Const {
	return Const{
		ImportableBase: NewImportableBase(name, prefix, layer == 0),
		valueType:      valueType,
		layer:          layer,
	}
}

func (c Const) ValueType() ValueType {
	return c.valueType
}

func (c Const) Layer() int {
	return c.layer
}

func (c Const) LayerName() string {
	return fmt.Sprintf("%s_%d", c.name, c.layer)
}

func (c Const) IsConstant() bool {
	return true
}

func (c *Const) SetValueType(valueType ValueType) {
	c.valueType = valueType
}

type ConstDefinition struct {
	constants []Const
	values    []Expression
}

func (c ConstDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_CONST_DEFINITION
}

func (c ConstDefinition) AssignmentType() AssignmentType {
	return ASSIGNMENT_TYPE_VALUE
}

func (c ConstDefinition) Constants() []Const {
	return c.constants
}

func (c ConstDefinition) Values() []Expression {
	return c.values
}

type ConstEvaluation struct {
	Const
}

func (c ConstEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_CONST_EVALUATION
}
