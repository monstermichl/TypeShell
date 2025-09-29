package parser

type TypeDeclaration struct {
	ImportableBase
	valueType ValueType
}

func NewTypeDeclaration(name string, prefix string, valueType ValueType, global bool) TypeDeclaration {
	return TypeDeclaration{
		ImportableBase: NewImportableBase(name, prefix, global),
		valueType:      valueType,
	}
}

func (t TypeDeclaration) StatementType() StatementType {
	return STATEMENT_TYPE_TYPE_DECLARATION
}

func (t TypeDeclaration) ValueType() ValueType {
	return t.valueType
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
