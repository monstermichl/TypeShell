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

type StructValue struct {
	name  string
	value Expression
}

func (v StructValue) Name() string {
	return v.name
}

func (v StructValue) ValueType() ValueType {
	return v.Value().ValueType()
}

func (v StructValue) Value() Expression {
	return v.value
}

type StructDefinition struct {
	valueType ValueType
	values    []StructValue
}

func (d StructDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_STRUCT_DEFINITION
}

func (d StructDefinition) ValueType() ValueType {
	return d.valueType
}

func (d StructDefinition) IsConstant() bool {
	return false
}

func (d StructDefinition) Values() []StructValue {
	return d.values
}
