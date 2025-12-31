package parser

type Panic struct {
	expression Expression
}

func (p Panic) Expression() Expression {
	return p.expression
}
