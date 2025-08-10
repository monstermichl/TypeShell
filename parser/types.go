package parser

import "fmt"

type StatementType string
type DataType string
type CompareOperator = string
type UnaryOperator = string
type BinaryOperator = string
type LogicalOperator = string

type ValueType struct {
	dataType DataType
	isSlice  bool
}

func NewValueType(dataType DataType, isSlice bool) ValueType {
	return ValueType{
		dataType,
		isSlice,
	}
}

func (vt ValueType) DataType() DataType {
	return vt.dataType
}

func (vt ValueType) IsSlice() bool {
	return vt.isSlice
}

func (vt ValueType) String() string {
	s := string(vt.dataType)

	if vt.isSlice {
		s = fmt.Sprintf("[]%s", s)
	}
	return s
}

func (vt ValueType) Equals(valueType ValueType) bool {
	return vt.DataType() == valueType.DataType() && vt.IsSlice() == valueType.IsSlice()
}

func (vt ValueType) IsBool() bool {
	return vt.isNonSliceType(DATA_TYPE_BOOLEAN)
}

func (vt ValueType) IsInt() bool {
	return vt.isNonSliceType(DATA_TYPE_INTEGER)
}

func (vt ValueType) IsString() bool {
	return vt.isNonSliceType(DATA_TYPE_STRING)
}

func (vt ValueType) isNonSliceType(dataType DataType) bool {
	return vt.DataType() == dataType && !vt.IsSlice()
}

const (
	STATEMENT_TYPE_PROGRAM                        StatementType = "program"
	STATEMENT_TYPE_BOOL_LITERAL                   StatementType = "boolean"
	STATEMENT_TYPE_INT_LITERAL                    StatementType = "integer"
	STATEMENT_TYPE_STRING_LITERAL                 StatementType = "string"
	STATEMENT_TYPE_STRING_SUBSCRIPT               StatementType = "string subscript"
	STATEMENT_TYPE_NIL_LITERAL                    StatementType = "nil"
	STATEMENT_TYPE_UNARY_OPERATION                StatementType = "unary operation"
	STATEMENT_TYPE_BINARY_OPERATION               StatementType = "binary operation"
	STATEMENT_TYPE_LOGICAL_OPERATION              StatementType = "logical operation"
	STATEMENT_TYPE_COMPARISON                     StatementType = "comparison"
	STATEMENT_TYPE_VAR_DEFINITION                 StatementType = "variable definition"
	STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT StatementType = "variable definition func assignment"
	STATEMENT_TYPE_VAR_ASSIGNMENT                 StatementType = "variable assignment"
	STATEMENT_TYPE_VAR_ASSIGNMENT_CALL_ASSIGNMENT StatementType = "variable assignment func assignment"
	STATEMENT_TYPE_VAR_EVALUATION                 StatementType = "variable evaluation"
	STATEMENT_TYPE_GROUP                          StatementType = "group"
	STATEMENT_TYPE_FUNCTION_DEFINITION            StatementType = "function definition"
	STATEMENT_TYPE_FUNCTION_CALL                  StatementType = "function call"
	STATEMENT_TYPE_APP_CALL                       StatementType = "app call"
	STATEMENT_TYPE_RETURN                         StatementType = "return"
	STATEMENT_TYPE_IF                             StatementType = "if"
	STATEMENT_TYPE_FOR                            StatementType = "for"
	STATEMENT_TYPE_FOR_RANGE                      StatementType = "for range"
	STATEMENT_TYPE_BREAK                          StatementType = "break"
	STATEMENT_TYPE_CONTINUE                       StatementType = "continue"
	STATEMENT_TYPE_STRUCT_DEFINITION              StatementType = "struct definition"
	STATEMENT_TYPE_INSTANTIATION                  StatementType = "instantiation"
	STATEMENT_TYPE_PRINT                          StatementType = "print"
	STATEMENT_TYPE_ITOA                           StatementType = "itoa"
	STATEMENT_TYPE_EXISTS                         StatementType = "exists"
	STATEMENT_TYPE_PANIC                          StatementType = "panic"
	STATEMENT_TYPE_LEN                            StatementType = "len"
	STATEMENT_TYPE_INPUT                          StatementType = "input"
	STATEMENT_TYPE_COPY                           StatementType = "copy"
	STATEMENT_TYPE_READ                           StatementType = "read"
	STATEMENT_TYPE_WRITE                          StatementType = "write"
	STATEMENT_TYPE_SLICE_INSTANTIATION            StatementType = "slice instantiation"
	STATEMENT_TYPE_SLICE_ASSIGNMENT               StatementType = "slice assignment"
	STATEMENT_TYPE_SLICE_EVALUATION               StatementType = "slice evaluation"
)

const (
	DATA_TYPE_UNKNOWN  DataType = "unknown"
	DATA_TYPE_MULTIPLE DataType = "multiple"
	DATA_TYPE_BOOLEAN  DataType = "bool"
	DATA_TYPE_INTEGER  DataType = "int"
	DATA_TYPE_STRING   DataType = "string"
	DATA_TYPE_ERROR    DataType = DATA_TYPE_STRING
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
