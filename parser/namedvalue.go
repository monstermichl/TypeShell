package parser

type NamedValue interface {
	Name() string
	ValueType() ValueType
	Global() bool
	Public() bool
	IsConstant() bool
}

type NamedValuesDefinition struct {
	assignments []Assignment
}

func (v NamedValuesDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_NAMED_VALUES_DEFINITION
}

func (v *NamedValuesDefinition) AddAssignment(assignment Assignment) {
	v.assignments = append(v.assignments, assignment)
}

func (v NamedValuesDefinition) Assignments() []Assignment {
	return v.assignments
}
