package parser

type For struct {
	condition Expression
	body      []Statement
}

func (f For) StatementType() StatementType {
	return STATEMENT_TYPE_FOR
}

func (f For) Condition() Expression {
	return f.condition
}

func (f For) Body() []Statement {
	return f.body
}
