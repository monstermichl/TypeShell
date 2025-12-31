package parser

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
