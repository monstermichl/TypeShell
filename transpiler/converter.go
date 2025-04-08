package transpiler

import "github.com/monstermichl/typeshell/parser"

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

type Converter interface {
	// Common methods
	BoolToString(value bool) string
	IntToString(value int) string
	StringToString(value string) string
	Dump() (string, error)

	// Statement methods
	ProgramStart() error
	ProgramEnd() error
	VarDefinition(name string, value string, global bool) error
	VarAssignment(name string, value string, global bool) error
	SliceAssignment(name string, index string, value string, global bool) error
	FuncStart(name string, params []string, returnType parser.ValueType) error
	FuncEnd() error
	Return(value string, valueType parser.ValueType) error
	IfStart(condition string) error
	IfEnd() error
	ElseIfStart(condition string) error
	ElseIfEnd() error
	ElseStart() error
	ElseEnd() error
	ForStart() error
	ForCondition(condition string) error
	ForEnd() error
	Break() error
	Continue() error
	Print(value []string) error
	Nop() error

	// Expression methods
	UnaryOperation(expr string, operator parser.UnaryOperator, valueType parser.ValueType, valueUsed bool) (string, error)
	BinaryOperation(left string, operator parser.BinaryOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error)
	Comparison(left string, operator parser.CompareOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error)
	LogicalOperation(left string, operator parser.LogicalOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error)
	VarEvaluation(name string, valueUsed bool, global bool) (string, error)
	SliceInstantiation(values []string, valueUsed bool) (string, error)
	SliceEvaluation(name string, index string, valueUsed bool, global bool) (string, error)
	SliceLen(name string, valueUsed bool, global bool) (string, error)
	Group(value string, valueUsed bool) (string, error)
	FuncCall(name string, args []string, valueType parser.ValueType, valueUsed bool) (string, error)
	AppCall(calls []AppCall, valueUsed bool) (string, error)
	Input(prompt string, valueUsed bool) (string, error)
}
