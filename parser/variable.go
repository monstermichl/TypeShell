package parser

type Variable struct {
	name      string
	valueType ValueType
	global    bool
}

func NewVariable(name string, valueType ValueType, global bool) Variable {
	return Variable{
		name,
		valueType,
		global,
	}
}

func (v Variable) Name() string {
	return v.name
}

func (v Variable) ValueType() ValueType {
	return v.valueType
}

func (v Variable) Global() bool {
	return v.global
}

type VariableAssignment struct {
	Variable
	value Expression
}

func (v VariableAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_ASSIGNMENT
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
	Variable
}

func (e VariableEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_EVALUATION
}
