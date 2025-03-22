package parser

type Program struct {
	body []Statement
}

func (p Program) StatementType() StatementType {
	return STATEMENT_TYPE_PROGRAM
}

func (p Program) Body() []Statement {
	return p.body
}
