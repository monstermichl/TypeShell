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
	variables []Variable
	values    []Expression
}

func (v VariableDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_DEFINITION
}

func (v VariableDefinition) Variables() []Variable {
	return v.variables
}

func (v VariableDefinition) Values() []Expression {
	return v.values
}

type VariableDefinitionCallAssignment struct {
	variables []Variable
	call      FunctionCall
}

func (v VariableDefinitionCallAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT
}

func (v VariableDefinitionCallAssignment) Variables() []Variable {
	return v.variables
}

func (v VariableDefinitionCallAssignment) Call() FunctionCall {
	return v.call
}

type VariableEvaluation struct {
	Variable
}

func (e VariableEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_EVALUATION
}
