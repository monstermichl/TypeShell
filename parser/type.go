package parser

type TypeDefinition struct {
	name string
}

func (t TypeDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_TYPE_DEFINITION
}

func (t TypeDefinition) Name() string {
	return t.name
}

type TypeInstantiation struct {
	value     Expression
	valueType ValueType
}

func (t TypeInstantiation) StatementType() StatementType {
	return STATEMENT_TYPE_TYPE_INSTANTIATION
}

func (t TypeInstantiation) ValueType() ValueType {
	return t.valueType
}

func (t TypeInstantiation) Value() Expression {
	return t.value
}
