package parser

type Exists struct {
	path Expression
}

func (e Exists) ExprFn() {}

func (e Exists) Path() Expression {
	return e.path
}
