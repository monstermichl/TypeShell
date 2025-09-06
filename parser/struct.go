package parser

type StructField struct {
	name      string
	valueType ValueType
}

func (f StructField) Name() string {
	return f.name
}

func (f StructField) ValueType() ValueType {
	return f.valueType
}

type StructDeclaration struct {
	fields []StructField
}

func (d StructDeclaration) StatementType() StatementType {
	return STATEMENT_TYPE_STRUCT_DECLARATION
}

func (d StructDeclaration) Fields() []StructField {
	return d.fields
}
