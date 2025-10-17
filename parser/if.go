package parser

type If struct {
	Condition Expression
	Body      Block
	Else      Statement
}

func (i If) StatementType() StatementType {
	return STATEMENT_TYPE_IF
}
