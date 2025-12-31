package parser

type Return struct {
	values []Expression
}

func (r Return) Values() []Expression {
	return r.values
}
