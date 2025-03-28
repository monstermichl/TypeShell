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

type ForRange struct {
	indexVar Variable
	valueVar Variable
	slice    Expression
	body     []Statement
}

func (f ForRange) StatementType() StatementType {
	return STATEMENT_TYPE_FOR_RANGE
}

func (f ForRange) IndexVar() Variable {
	return f.indexVar
}

func (f ForRange) ValueVar() Variable {
	return f.valueVar
}

func (f ForRange) Slice() Expression {
	return f.slice
}

func (f ForRange) Body() []Statement {
	return f.body
}
