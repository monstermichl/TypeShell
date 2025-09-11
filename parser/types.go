package parser

import (
	"fmt"
)

type StatementType string
type AssignmentType string
type TypeKind string
type CompareOperator = string
type UnaryOperator = string
type BinaryOperator = string
type LogicalOperator = string

const (
	TypeKindUnknown  TypeKind = "unknown"
	TypeKindBool     TypeKind = "bool"
	TypeKindInt      TypeKind = "int"
	TypeKindString   TypeKind = "string"
	TypeKindError    TypeKind = "error"
	TypeKindStruct   TypeKind = "struct"
	TypeKindMultiple TypeKind = "multiple"
)

type Type interface {
	Name() string
	IsAlias() bool
	Kind() TypeKind
	Base() Type
	Equals(compareType Type) bool
	ElementaryType() Type
	AliasedType() Type
}

func equals(t1 Type, t2 Type) bool {
	base1 := t1.AliasedType()
	base2 := t2.AliasedType()

	return base1.Name() == base2.Name() && base1.Kind() == base2.Kind()
}

func elementaryType(t Type) Type {
	base := t.Base()

	if base != nil {
		return base.ElementaryType()
	}
	return t
}

func aliasedType(t Type) Type {
	base := t.Base()

	if t.IsAlias() && base != nil {
		return base.AliasedType()
	}
	return t
}

type TypeUnknown struct{}

func (t TypeUnknown) Name() string                 { return "unknown" }
func (t TypeUnknown) IsAlias() bool                { return false }
func (t TypeUnknown) Kind() TypeKind               { return TypeKindUnknown }
func (t TypeUnknown) Base() Type                   { return nil }
func (t TypeUnknown) Equals(compareType Type) bool { return equals(t, compareType) }
func (t TypeUnknown) ElementaryType() Type         { return elementaryType(t) }
func (t TypeUnknown) AliasedType() Type            { return aliasedType(t) }

type TypeCustom struct {
	name    string
	isAlias bool
	kind    TypeKind
	base    Type
}

func NewTypeCustom(name string, isAlias bool, kind TypeKind, base Type) TypeCustom {
	return TypeCustom{
		name,
		isAlias,
		kind,
		base,
	}
}

func (t TypeCustom) Name() string                 { return t.name }
func (t TypeCustom) IsAlias() bool                { return t.isAlias }
func (t TypeCustom) Kind() TypeKind               { return t.kind }
func (t TypeCustom) Base() Type                   { return t.base }
func (t TypeCustom) Equals(compareType Type) bool { return equals(t, compareType) }
func (t TypeCustom) ElementaryType() Type         { return elementaryType(t) }
func (t TypeCustom) AliasedType() Type            { return aliasedType(t) }

type TypeBool struct{}

func (t TypeBool) Name() string                 { return "bool" }
func (t TypeBool) IsAlias() bool                { return false }
func (t TypeBool) Kind() TypeKind               { return TypeKindBool }
func (t TypeBool) Base() Type                   { return nil }
func (t TypeBool) Equals(compareType Type) bool { return equals(t, compareType) }
func (t TypeBool) ElementaryType() Type         { return elementaryType(t) }
func (t TypeBool) AliasedType() Type            { return aliasedType(t) }

type TypeInt struct{}

func (t TypeInt) Name() string                 { return "int" }
func (t TypeInt) IsAlias() bool                { return false }
func (t TypeInt) Kind() TypeKind               { return TypeKindInt }
func (t TypeInt) Base() Type                   { return nil }
func (t TypeInt) Equals(compareType Type) bool { return equals(t, compareType) }
func (t TypeInt) ElementaryType() Type         { return elementaryType(t) }
func (t TypeInt) AliasedType() Type            { return aliasedType(t) }

type TypeString struct{}

func (t TypeString) Name() string                 { return "string" }
func (t TypeString) IsAlias() bool                { return false }
func (t TypeString) Kind() TypeKind               { return TypeKindString }
func (t TypeString) Base() Type                   { return nil }
func (t TypeString) Equals(compareType Type) bool { return equals(t, compareType) }
func (t TypeString) ElementaryType() Type         { return elementaryType(t) }
func (t TypeString) AliasedType() Type            { return aliasedType(t) }

type TypeError struct{}

func (t TypeError) Name() string                 { return "error" }
func (t TypeError) IsAlias() bool                { return true }
func (t TypeError) Kind() TypeKind               { return TypeKindError }
func (t TypeError) Base() Type                   { return TypeString{} }
func (t TypeError) Equals(compareType Type) bool { return equals(t, compareType) }
func (t TypeError) ElementaryType() Type         { return elementaryType(t) }
func (t TypeError) AliasedType() Type            { return aliasedType(t) }

type TypeMultiple struct {
	types []Type
}

func (t TypeMultiple) Name() string   { return "multiple" }
func (t TypeMultiple) IsAlias() bool  { return false }
func (t TypeMultiple) Kind() TypeKind { return TypeKindMultiple }
func (t TypeMultiple) Base() Type     { return nil }
func (t TypeMultiple) Equals(c Type) bool {
	ct, isType := c.(TypeMultiple)

	if isType {
		typesT1 := t.Types()
		typesT2 := ct.Types()

		if len(typesT1) != len(typesT2) {
			return false
		}

		for i, t1 := range typesT1 {
			if !t1.Equals(typesT2[i]) {
				return false
			}
		}
	}
	return false
}
func (t TypeMultiple) ElementaryType() Type { return elementaryType(t) }
func (t TypeMultiple) AliasedType() Type    { return aliasedType(t) }

func (t TypeMultiple) Types() []Type {
	return t.types
}

type ValueType struct {
	t       Type
	isSlice bool
}

func NewValueType(t Type, isSlice bool) ValueType {
	return ValueType{
		t,
		isSlice,
	}
}

func (vt ValueType) Type() Type {
	return vt.t
}

func (vt ValueType) IsSlice() bool {
	return vt.isSlice
}

func (vt ValueType) String() string {
	s := string(vt.Type().Name())

	if vt.isSlice {
		s = fmt.Sprintf("[]%s", s)
	}
	return s
}

func (vt ValueType) Equals(valueType ValueType) bool {
	return vt.Type().Equals(valueType.Type()) && vt.IsSlice() == valueType.IsSlice()
}

func (vt ValueType) IsBool() bool {
	return vt.isNonSliceType(TypeBool{})
}

func (vt ValueType) IsInt() bool {
	return vt.isNonSliceType(TypeInt{})
}

func (vt ValueType) IsString() bool {
	return vt.isNonSliceType(TypeString{})
}

func (vt ValueType) isNonSliceType(t Type) bool {
	return vt.Type().Equals(t) && !vt.IsSlice()
}

const (
	STATEMENT_TYPE_NOP                             StatementType = "nop"
	STATEMENT_TYPE_PROGRAM                         StatementType = "program"
	STATEMENT_TYPE_TYPE_DECLARATION                StatementType = "type declaration"
	STATEMENT_TYPE_TYPE_DEFINITION                 StatementType = "type definition"
	STATEMENT_TYPE_STRUCT_DECLARATION              StatementType = "struct declaration"
	STATEMENT_TYPE_STRUCT_DEFINITION               StatementType = "struct definition"
	STATEMENT_TYPE_STRUCT_ASSIGNMENT               StatementType = "struct assignment"
	STATEMENT_TYPE_STRUCT_EVALUATION               StatementType = "struct evaluation"
	STATEMENT_TYPE_BOOL_LITERAL                    StatementType = "boolean"
	STATEMENT_TYPE_INT_LITERAL                     StatementType = "integer"
	STATEMENT_TYPE_STRING_LITERAL                  StatementType = "string"
	STATEMENT_TYPE_STRING_SUBSCRIPT                StatementType = "string subscript"
	STATEMENT_TYPE_NIL_LITERAL                     StatementType = "nil"
	STATEMENT_TYPE_UNARY_OPERATION                 StatementType = "unary operation"
	STATEMENT_TYPE_BINARY_OPERATION                StatementType = "binary operation"
	STATEMENT_TYPE_LOGICAL_OPERATION               StatementType = "logical operation"
	STATEMENT_TYPE_COMPARISON                      StatementType = "comparison"
	STATEMENT_TYPE_CONST_DEFINITION                StatementType = "constant definition"
	STATEMENT_TYPE_CONST_EVALUATION                StatementType = "constant evaluation"
	STATEMENT_TYPE_NAMED_VALUES_DEFINITION         StatementType = "named values definition"
	STATEMENT_TYPE_VAR_DEFINITION_VALUE_ASSIGNMENT StatementType = "variable definition value assignment"
	STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT  StatementType = "variable definition call assignment"
	STATEMENT_TYPE_VAR_ASSIGNMENT                  StatementType = "variable assignment"
	STATEMENT_TYPE_VAR_ASSIGNMENT_VALUE_ASSIGNMENT StatementType = "variable assignment value assignment"
	STATEMENT_TYPE_VAR_ASSIGNMENT_CALL_ASSIGNMENT  StatementType = "variable assignment call assignment"
	STATEMENT_TYPE_VAR_EVALUATION                  StatementType = "variable evaluation"
	STATEMENT_TYPE_GROUP                           StatementType = "group"
	STATEMENT_TYPE_FUNCTION_DEFINITION             StatementType = "function definition"
	STATEMENT_TYPE_FUNCTION_CALL                   StatementType = "function call"
	STATEMENT_TYPE_APP_CALL                        StatementType = "app call"
	STATEMENT_TYPE_RETURN                          StatementType = "return"
	STATEMENT_TYPE_IF                              StatementType = "if"
	STATEMENT_TYPE_FOR                             StatementType = "for"
	STATEMENT_TYPE_FOR_RANGE                       StatementType = "for range"
	STATEMENT_TYPE_BREAK                           StatementType = "break"
	STATEMENT_TYPE_CONTINUE                        StatementType = "continue"
	STATEMENT_TYPE_IOTA                            StatementType = "iota"
	STATEMENT_TYPE_INSTANTIATION                   StatementType = "instantiation"
	STATEMENT_TYPE_PRINT                           StatementType = "print"
	STATEMENT_TYPE_ITOA                            StatementType = "itoa"
	STATEMENT_TYPE_EXISTS                          StatementType = "exists"
	STATEMENT_TYPE_PANIC                           StatementType = "panic"
	STATEMENT_TYPE_LEN                             StatementType = "len"
	STATEMENT_TYPE_INPUT                           StatementType = "input"
	STATEMENT_TYPE_COPY                            StatementType = "copy"
	STATEMENT_TYPE_READ                            StatementType = "read"
	STATEMENT_TYPE_WRITE                           StatementType = "write"
	STATEMENT_TYPE_SLICE_INSTANTIATION             StatementType = "slice instantiation"
	STATEMENT_TYPE_SLICE_ASSIGNMENT                StatementType = "slice assignment"
	STATEMENT_TYPE_SLICE_EVALUATION                StatementType = "slice evaluation"
)

const (
	ASSIGNMENT_TYPE_VALUE AssignmentType = "value"
	ASSIGNMENT_TYPE_CALL  AssignmentType = "call"
)

const (
	COMPARE_OPERATOR_EQUAL            CompareOperator = "=="
	COMPARE_OPERATOR_NOT_EQUAL        CompareOperator = "!="
	COMPARE_OPERATOR_LESS             CompareOperator = "<"
	COMPARE_OPERATOR_LESS_OR_EQUAL    CompareOperator = "<="
	COMPARE_OPERATOR_GREATER          CompareOperator = ">"
	COMPARE_OPERATOR_GREATER_OR_EQUAL CompareOperator = ">="
)

const (
	UNARY_OPERATOR_NEGATE UnaryOperator = "!"
)

const (
	BINARY_OPERATOR_MULTIPLICATION BinaryOperator = "*"
	BINARY_OPERATOR_DIVISION       BinaryOperator = "/"
	BINARY_OPERATOR_MODULO         BinaryOperator = "%"
	BINARY_OPERATOR_ADDITION       BinaryOperator = "+"
	BINARY_OPERATOR_SUBTRACTION    BinaryOperator = "-"
)

const (
	LOGICAL_OPERATOR_AND LogicalOperator = "&&"
	LOGICAL_OPERATOR_OR  LogicalOperator = "||"
)
