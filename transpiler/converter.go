package transpiler

import (
	"github.com/monstermichl/typeshell/parser"
	"github.com/monstermichl/typeshell/typechecker"
)

type AppCall struct {
	name string
	args []string
}

func (c AppCall) Name() string {
	return c.name
}

func (c AppCall) Args() []string {
	return c.args
}

type Condition struct {
	condition string
	operator  parser.LogicalOperator
	next      *Condition
}

func (c Condition) Condition() string {
	return c.condition
}

func (c Condition) Operator() parser.LogicalOperator {
	return c.operator
}

func (c Condition) Next() *Condition {
	return c.next
}

type ReturnValue struct {
	value     string
	valueType typechecker.Type
}

func (rv ReturnValue) Value() string {
	return rv.value
}

func (rv ReturnValue) ValueType() typechecker.Type {
	return rv.valueType
}

type StructValue struct {
	name  string
	value string
}

func (sv StructValue) Name() string {
	return sv.name
}

func (sv StructValue) Value() string {
	return sv.value
}

type Converter interface {
	// Common methods
	StringToString(value string) string
	VarName(name string, global bool) string
	Dump() (string, error)
	Extension() string

	// Statement methods
	ProgramStart() error
	ProgramEnd() error
	VarDefinition(name string, value string, global bool) error
	VarAssignment(name string, value string, global bool) error
	SliceAssignment(name string, index string, value string, defaultValue string, global bool) error
	StructAssignment(name string, field string, value string, global bool) error
	FuncStart(name string, params []string, returnTypes []typechecker.Type) error
	FuncEnd() error
	Return(values []ReturnValue) error
	IfStart(condition string) error
	IfEnd() error
	ElseIfStart(condition string) error
	ElseIfEnd() error
	ElseStart() error
	ElseEnd() error
	ForStart() error
	ForIncrementStart() error
	ForIncrementEnd() error
	ForCondition(condition string) error
	ForEnd() error
	Break() error
	Continue() error
	Print(value []string) error
	Panic(value string) error
	Unsafe(code string) error
	WriteFile(path string, content string, append string) error
	Nop() error

	// Expression methods
	UnaryOperation(expr string, operator parser.UnaryOperator, valueType typechecker.Type, valueUsed bool) (string, error)
	BinaryOperation(left string, operator parser.BinaryOperator, right string, valueType typechecker.Type, valueUsed bool) (string, error)
	Comparison(left string, operator parser.CompareOperator, right string, valueType typechecker.Type, valueUsed bool) (string, error)
	LogicalOperation(left string, operator parser.LogicalOperator, right string, valueType typechecker.Type, valueUsed bool) (string, error)
	VarEvaluation(name string, valueUsed bool, global bool) (string, error)
	SliceInstantiation(values []string, valueUsed bool) (string, error)
	SliceEvaluation(name string, index string, valueUsed bool) (string, error)
	SliceLen(name string, valueUsed bool) (string, error)
	StructInitialization(values []StructValue, valueUsed bool) (string, error)
	StructEvaluation(name string, field string, valueUsed bool) (string, error)
	StringSubscript(value string, startIndex string, endIndex string, valueUsed bool) (string, error)
	StringLen(value string, valueUsed bool) (string, error)
	Group(value string, valueUsed bool) (string, error)
	FuncCall(name string, args []string, returnTypes []typechecker.Type, valueUsed bool) ([]string, error)
	AppCall(calls []AppCall, valueUsed bool) ([]string, error)
	Input(prompt string, valueUsed bool) (string, error)
	Copy(destination string, source string, valueUsed bool, global bool) (string, error)
	Exists(path string, valueUsed bool) (string, error)
	ReadFile(path string, valueUsed bool) (string, error)
}
