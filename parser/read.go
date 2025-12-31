package parser

type Read struct {
	path Expression
}

func (r Read) ExprFn() {}

func (r Read) Path() Expression {
	return r.path
}
