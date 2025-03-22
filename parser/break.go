package parser

type Break struct {
}

func (b Break) StatementType() StatementType {
	return STATEMENT_TYPE_BREAK
}
