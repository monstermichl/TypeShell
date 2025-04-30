package parser

type Panic struct {
	expression Expression
}

func (p Panic) StatementType() StatementType {
	return STATEMENT_TYPE_PANIC
}

func (p Panic) Expression() Expression {
	return p.expression
}
