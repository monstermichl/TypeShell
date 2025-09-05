package parser

type AppCall struct {
	name string
	args []Expression
	next *AppCall
}

func (a AppCall) StatementType() StatementType {
	return STATEMENT_TYPE_APP_CALL
}

func (a AppCall) ValueType() ValueType {
	return NewValueType(DATA_TYPE_MULTIPLE, false)
}

func (a AppCall) IsConstant() bool {
	return false
}

func (a AppCall) Name() string {
	return a.name
}

func (a AppCall) Args() []Expression {
	return a.args
}

func (a AppCall) Next() *AppCall {
	return a.next
}

func (a AppCall) ReturnTypes() []ValueType {
	return []ValueType{
		NewValueType(DATA_TYPE_STRING, false),  // stdout
		NewValueType(DATA_TYPE_STRING, false),  // stderr
		NewValueType(DATA_TYPE_INTEGER, false), // error code
	}
}
