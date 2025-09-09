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
	fields []StructField
}

func (d StructDeclaration) StatementType() StatementType {
	return STATEMENT_TYPE_STRUCT_DECLARATION
}

func (d StructDeclaration) Fields() []StructField {
	return d.fields
}

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
