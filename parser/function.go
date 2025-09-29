package parser

type FunctionDefinition struct {
	ImportableBase
	returnTypes []ValueType
	params      []Param
	body        []Statement
}

func NewFunctionDefinition(name string, prefix string, returnTypes []ValueType, params []Param, body []Statement) FunctionDefinition {
	return FunctionDefinition{
		ImportableBase: NewImportableBase(name, prefix, true),
		returnTypes:    returnTypes,
		params:         params,
		body:           body,
	}
}

func (e FunctionDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_FUNCTION_DEFINITION
}

func (e FunctionDefinition) ValueType() ValueType {
	return functionValueType(e.returnTypes)
}

func (e FunctionDefinition) IsConstant() bool {
	return false
}

func (e FunctionDefinition) ReturnTypes() []ValueType {
	return e.returnTypes
}

func (e FunctionDefinition) Params() []Param {
	return e.params
}

func (e FunctionDefinition) Body() []Statement {
	return e.body
}

type FunctionCall struct {
	FunctionDefinition
	arguments   []Expression
}

func (e FunctionCall) StatementType() StatementType {
	return STATEMENT_TYPE_FUNCTION_CALL
}

func (e FunctionCall) IsConstant() bool {
	return false
}

func (e FunctionCall) Args() []Expression {
	return e.arguments
}

func functionValueType(returnTypes []ValueType) ValueType {
	var valueType ValueType
	length := len(returnTypes)

	if length > 1 {
		valueType = NewValueType(NewTypeMultiple(), false)
	} else if length > 0 {
		valueType = returnTypes[0]
	}
	return valueType
}
