package parser

import "fmt"

type Const struct {
	name      string
	prefix    string
	valueType ValueType
	layer     int
}

func NewConst(name string, prefix string, valueType ValueType, layer int) Const {
	return Const{
		name,
		prefix,
		valueType,
		layer,
	}
}

func (c Const) Name() string {
	return c.name
}

func (c Const) Prefix() string {
	return c.prefix
}

func (c Const) ValueType() ValueType {
	return c.valueType
}

func (c Const) ImportableType() ImportableType {
	return ImportableTypeConstant
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

func (c Const) Global() bool {
	return c.layer == 0
}

func (c Const) Public() bool {
	return isPublic(c.Name())
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
