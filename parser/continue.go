package parser

type Continue struct {
}

func (c Continue) StatementType() StatementType {
	return STATEMENT_TYPE_CONTINUE
}
