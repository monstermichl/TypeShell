package parser

type FunctionDefinition struct {
	name        string
	returnTypes []ValueType
	params      []Variable
	body        []Statement
}

func (e FunctionDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_FUNCTION_DEFINITION
}

func (e FunctionDefinition) Name() string {
	return e.name
}

func (e FunctionDefinition) ValueType() ValueType {
	return functionValueType(e.returnTypes)
}

func (e FunctionDefinition) ReturnTypes() []ValueType {
	return e.returnTypes
}

func (e FunctionDefinition) Params() []Variable {
	return e.params
}

func (e FunctionDefinition) Body() []Statement {
	return e.body
}

type FunctionCall struct {
	name        string
	returnTypes []ValueType
	arguments   []Expression
}

func (e FunctionCall) StatementType() StatementType {
	return STATEMENT_TYPE_FUNCTION_CALL
}

func (e FunctionCall) Name() string {
	return e.name
}

func (e FunctionCall) ValueType() ValueType {
	return functionValueType(e.returnTypes)
}

func (e FunctionCall) ReturnTypes() []ValueType {
	return e.returnTypes
}

func (e FunctionCall) Args() []Expression {
	return e.arguments
}

func functionValueType(returnTypes []ValueType) ValueType {
	var valueType ValueType

	if len(returnTypes) > 1 {
		valueType = NewValueType(DATA_TYPE_MULTIPLE, returnTypes[0].IsSlice())
	} else {
		valueType = returnTypes[0]
	}
	return valueType
}
