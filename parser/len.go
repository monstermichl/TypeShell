package parser

type Len struct {
	expression Expression
}

func (l Len) ExprFn() {}

func (l Len) Expression() Expression {
	return l.expression
}
