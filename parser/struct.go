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

type StructDefinition struct {
	name   string
	fields []StructField
}

func NewStructDefinition(name string, fields []StructField) StructDefinition {
	return StructDefinition{
		name,
		fields,
	}
}

func (d StructDefinition) Fields() []StructField {
	return d.fields
}

func (d StructDefinition) Name() string {
	return d.name
}

func (d StructDefinition) IsAlias() bool {
	return false
}

func (d StructDefinition) Kind() TypeKind {
	return TypeKindStruct
}

func (d StructDefinition) Base() Type {
	return nil
}

func (d StructDefinition) Equals(c Type) bool {
	compareType, isDeclaration := c.(StructDefinition)

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

func (t StructDefinition) ElementaryType() Type { return elementaryType(t) }
func (t StructDefinition) AliasedType() Type    { return aliasedType(t) }

func (d StructDefinition) FindField(name string) (StructField, error) {
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

func NewStructValue(name string, valueType ValueType, value Expression) StructValue {
	return StructValue{
		StructField: StructField{
			name,
			valueType,
		},
		value: value,
	}
}

func (v StructValue) Value() Expression {
	return v.value
}

type StructInitialization struct {
	t      Type
	values []StructValue
}

func NewStructInitialization(t Type, values ...StructValue) StructInitialization {
	return StructInitialization{
		t,
		values,
	}
}

func (d StructInitialization) StatementType() StatementType {
	return STATEMENT_TYPE_STRUCT_DEFINITION
}

func (d StructInitialization) ValueType() ValueType {
	return NewValueType(d.t, false)
}

func (d StructInitialization) IsConstant() bool {
	return false
}

func (d StructInitialization) Values() []StructValue {
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
	value Expression
	field StructField
}

func (e StructEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_STRUCT_EVALUATION
}

func (e StructEvaluation) Value() Expression {
	return e.value
}

func (e StructEvaluation) ValueType() ValueType {
	return e.field.ValueType()
}

func (e StructEvaluation) IsConstant() bool {
	return false
}

func (e StructEvaluation) Field() StructField {
	return e.field
}
