package parser

type For struct {
	init      Statement
	condition Expression
	increment Statement
	body      []Statement
}

func (f For) StatementType() StatementType {
	return STATEMENT_TYPE_FOR
}

func (f For) Init() Statement {
	return f.init
}

func (f For) Condition() Expression {
	return f.condition
}

func (f For) Increment() Statement {
	return f.increment
}

func (f For) Body() []Statement {
	return f.body
}
