package parser

import "fmt"

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
	name   string
	fields []StructField
}

func NewStructDeclaration(name string, fields []StructField) StructDeclaration {
	return StructDeclaration{
		name,
		fields,
	}
}

func (d StructDeclaration) Fields() []StructField {
	return d.fields
}

func (d StructDeclaration) Name() string {
	return d.name
}

func (d StructDeclaration) IsAlias() bool {
	return false
}

func (d StructDeclaration) Kind() TypeKind {
	return TypeKindStruct
}

func (d StructDeclaration) Base() Type {
	return nil
}

func (d StructDeclaration) Equals(c Type) bool {
	compareType, isDeclaration := c.(StructDeclaration)

	if !isDeclaration {
		return false
	}
	fieldsD1 := d.Fields()
	fieldsD2 := compareType.Fields()

	if len(fieldsD1) != len(fieldsD2) {
		return false
	}

	for i, fieldD1 := range fieldsD1 {
		fieldD2 := fieldsD2[i]

		if fieldD1.Name() != fieldD2.Name() || !fieldD1.ValueType().Equals(fieldD2.ValueType()) {
			return false
		}
	}
	return true
}

func (t StructDeclaration) ElementaryType() Type { return elementaryType(t) }
func (t StructDeclaration) AliasedType() Type    { return aliasedType(t) }

func (d StructDeclaration) FindField(name string) (StructField, error) {
	for _, field := range d.Fields() {
		if field.Name() == name {
			return field, nil
		}
	}
	return StructField{}, fmt.Errorf(`struct field %s doesn't exist`, name)
}

type StructValue struct {
	StructField
	value Expression
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

type StructAssignment struct {
	Variable
	value StructValue
}

func (a StructAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_STRUCT_ASSIGNMENT
}

func (a StructAssignment) Value() StructValue {
	return a.value
}

type StructEvaluation struct {
	Variable
	field StructField
}

func (e StructEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_STRUCT_EVALUATION
}

func (e StructEvaluation) Field() StructField {
	return e.field
}
