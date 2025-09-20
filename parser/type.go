package parser

type TypeDeclaration struct {
	name      string
	valueType ValueType
	public    bool
}

func (t TypeDeclaration) StatementType() StatementType {
	return STATEMENT_TYPE_TYPE_DECLARATION
}

func (t TypeDeclaration) Name() string {
	return t.name
}

func (t TypeDeclaration) ValueType() ValueType {
	return t.valueType
}

func (t TypeDeclaration) Public() bool {
	return t.public
}

type TypeDefinition struct {
	value     Expression
	valueType ValueType
}

func (t TypeDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_TYPE_DEFINITION
}

func (t TypeDefinition) ValueType() ValueType {
	return t.valueType
}

func (t TypeDefinition) IsConstant() bool {
	return t.Value().IsConstant()
}

func (t TypeDefinition) Value() Expression {
	return t.value
}
