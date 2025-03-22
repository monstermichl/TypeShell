package parser

type Return struct {
	value Expression
}

func (r Return) StatementType() StatementType {
	return STATEMENT_TYPE_RETURN
}

func (r Return) Value() Expression {
	return r.value
}
