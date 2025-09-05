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

func (v Variable) IsConstant() bool {
	return false
}

func (v *Variable) SetValueType(valueType ValueType) {
	v.valueType = valueType
}

func (v Variable) Global() bool {
	return v.global
}

func (v Variable) Public() bool {
	return v.public
}

type VariableDefinitionValueAssignment struct {
	variables []Variable
	values    []Expression
}

func NewVariableDefinition(variables []Variable, values []Expression) VariableDefinitionValueAssignment {
	return VariableDefinitionValueAssignment{
		variables,
		values,
	}
}

func (v VariableDefinitionValueAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_DEFINITION_VALUE_ASSIGNMENT
}

func (v VariableDefinitionValueAssignment) AssignmentType() AssignmentType {
	return ASSIGNMENT_TYPE_VALUE
}

func (v VariableDefinitionValueAssignment) Variables() []Variable {
	return v.variables
}

func (v VariableDefinitionValueAssignment) Values() []Expression {
	return v.values
}

type VariableDefinitionCallAssignment struct {
	variables []Variable
	call      Call
}

func (v VariableDefinitionCallAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT
}

func (v VariableDefinitionCallAssignment) AssignmentType() AssignmentType {
	return ASSIGNMENT_TYPE_CALL
}

func (v VariableDefinitionCallAssignment) Variables() []Variable {
	return v.variables
}

func (v VariableDefinitionCallAssignment) Call() Call {
	return v.call
}

type VariableAssignment struct {
	assignments []Assignment
}

func (v VariableAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_ASSIGNMENT
}

func (v *VariableAssignment) AddAssignment(assignment Assignment) {
	v.assignments = append(v.assignments, assignment)
}

func (v VariableAssignment) Assignments() []Assignment {
	return v.assignments
}

type VariableAssignmentValueAssignment struct {
	variables []Variable
	values    []Expression
}

func (v VariableAssignmentValueAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_ASSIGNMENT_VALUE_ASSIGNMENT
}

func (v VariableAssignmentValueAssignment) AssignmentType() AssignmentType {
	return ASSIGNMENT_TYPE_VALUE
}

func (v VariableAssignmentValueAssignment) Variables() []Variable {
	return v.variables
}

func (v VariableAssignmentValueAssignment) Values() []Expression {
	return v.values
}

type VariableAssignmentCallAssignment struct {
	variables []Variable
	call      Call
}

func (v VariableAssignmentCallAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_ASSIGNMENT_CALL_ASSIGNMENT
}

func (v VariableAssignmentCallAssignment) AssignmentType() AssignmentType {
	return ASSIGNMENT_TYPE_CALL
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

func NewVariableEvaluation(name string, valueType ValueType, global bool, public bool) VariableEvaluation {
	return VariableEvaluation{
		Variable{
			name,
			valueType,
			global,
			public,
		},
	}
}

func (e VariableEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_VAR_EVALUATION
}
