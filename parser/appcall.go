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
	return ValueType{dataType: DATA_TYPE_STRING}
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
