package parser

type Read struct {
	path Expression
}

func (r Read) StatementType() StatementType {
	return STATEMENT_TYPE_READ
}

func (r Read) ValueType() ValueType {
	return NewValueType(DATA_TYPE_STRING, false)
}

func (r Read) IsConstant() bool {
	return false
}

func (r Read) Path() Expression {
	return r.path
}
