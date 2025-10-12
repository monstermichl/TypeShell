package parser

type Program struct {
	Statements []Statement
}

func (p Program) StatementType() StatementType {
	return STATEMENT_TYPE_PROGRAM
}
