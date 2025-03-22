package parser

type FunctionDefinition struct {
	name      string
	valueType ValueType
	params    []Variable
	body      []Statement
}

func (e FunctionDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_FUNCTION_DEFINITION
}

func (e FunctionDefinition) Name() string {
	return e.name
}

func (e FunctionDefinition) ValueType() ValueType {
	return e.valueType
}

func (e FunctionDefinition) Params() []Variable {
	return e.params
}

func (e FunctionDefinition) Body() []Statement {
	return e.body
}

type FunctionCall struct {
	name      string
	valueType ValueType
	arguments []Expression
}

func (e FunctionCall) StatementType() StatementType {
	return STATEMENT_TYPE_FUNCTION_CALL
}

func (e FunctionCall) Name() string {
	return e.name
}

func (e FunctionCall) ValueType() ValueType {
	return e.valueType
}

func (e FunctionCall) Args() []Expression {
	return e.arguments
}
