package parser

type Len struct {
	expression Expression
}

func (l Len) StatementType() StatementType {
	return STATEMENT_TYPE_LEN
}

func (l Len) ValueType() ValueType {
	return NewValueType(DATA_TYPE_INTEGER, false)
}

func (l Len) IsConstant() bool {
	return false
}

func (l Len) Expression() Expression {
	return l.expression
}
