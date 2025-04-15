package parser

type Return struct {
	values []Expression
}

func (r Return) StatementType() StatementType {
	return STATEMENT_TYPE_RETURN
}

func (r Return) Values() []Expression {
	return r.values
}
