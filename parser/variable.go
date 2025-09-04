package parser

type Variable struct {
	name      string
	valueType ValueType
	global    bool
	public    bool
}

func NewVariable(name string, valueType ValueType, global bool, public bool) Variable {
	return Variable{
		name,
		valueType,
		global,
		public,
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

func (v Variable) Public() bool {
	return v.public
}

func (v Variable) IsConstant() bool {
	return false
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
	call      Call
}

func (v VariableDefinitionCallAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT
}

func (v VariableDefinitionCallAssignment) Variables() []Variable {
	return v.variables
}

func (v VariableDefinitionCallAssignment) Call() Call {
	return v.call
}

type VariableAssignment struct {
	variables []Variable
	values    []Expression
}

func (v VariableAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_ASSIGNMENT
}

func (v VariableAssignment) Variables() []Variable {
	return v.variables
}

func (v VariableAssignment) Values() []Expression {
	return v.values
}

type VariableAssignmentCallAssignment struct {
	variables []Variable
	call      Call
}

func (v VariableAssignmentCallAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_ASSIGNMENT_CALL_ASSIGNMENT
}

func (v VariableAssignmentCallAssignment) Variables() []Variable {
	return v.variables
}

func (v VariableAssignmentCallAssignment) Call() Call {
	return v.call
}

type VariableEvaluation struct {
	Variable
}

func (e VariableEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_EVALUATION
}
