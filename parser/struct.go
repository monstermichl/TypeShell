package parser

type StructMember struct {
	name      string
	valueType ValueType
}

func (m StructMember) Name() string {
	return m.name
}

func (m StructMember) ValueType() ValueType {
	return m.valueType
}

type StructDefinition struct {
	name    string
	members []StructMember
}

func (d StructDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_STRUCT_DEFINITION
}

func (d StructDefinition) Name() string {
	return d.name
}

func (d StructDefinition) Members() []StructMember {
	return d.members
}
