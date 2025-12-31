package parser

type Print struct {
	expressions []Expression
}

func (p Print) Expressions() []Expression {
	return p.expressions
}
