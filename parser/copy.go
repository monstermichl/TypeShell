package parser

type Copy struct {
	destination Expression
	source      Expression
}

func (c Copy) ExprFn() {}

func (c Copy) Source() Expression {
	return c.source
}

func (c Copy) Destination() Expression {
	return c.destination
}
