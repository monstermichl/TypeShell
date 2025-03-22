package parser

type Print struct {
	expressions []Expression
}

func (p Print) StatementType() StatementType {
	return STATEMENT_TYPE_PRINT
}

func (p Print) Expressions() []Expression {
	return p.expressions
}
