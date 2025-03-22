package parser

type Variable struct {
	name      string
	valueType ValueType
}

func (v Variable) Name() string {
	return v.name
}

func (v Variable) ValueType() ValueType {
	return v.valueType
}

type VariableAssignment struct {
	variable Variable
	value    Expression
}

func (v VariableAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_ASSIGNMENT
}

func (v VariableAssignment) ValueType() ValueType {
	return v.variable.ValueType()
}

func (v VariableAssignment) Variable() Variable {
	return v.variable
}

func (v VariableAssignment) Value() Expression {
	return v.value
}

type VariableDefinition struct {
	variable Variable
	value    Expression
}

func (v VariableDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_DEFINITION
}

func (v VariableDefinition) Variable() Variable {
	return v.variable
}

func (v VariableDefinition) Value() Expression {
	return v.value
}

type VariableEvaluation struct {
	name      string
	valueType ValueType
}

func (e VariableEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_EVALUATION
}

func (e VariableEvaluation) ValueType() ValueType {
	return e.valueType
}

func (e VariableEvaluation) Name() string {
	return e.name
}
