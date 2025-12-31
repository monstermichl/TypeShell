package parser

type AppCall struct {
	name string
	args []Expression
	next *AppCall
}

func (a AppCall) ExprFn() {}

func (a AppCall) Name() string {
	return a.name
}

func (a AppCall) Args() []Expression {
	return a.args
}

func (a AppCall) Next() *AppCall {
	return a.next
}
