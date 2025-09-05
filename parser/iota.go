package parser

type Iota struct {
}

func (i Iota) StatementType() StatementType {
	return STATEMENT_TYPE_IOTA
}

func (i Iota) ValueType() ValueType {
	return NewValueType(DATA_TYPE_INTEGER, false)
}

func (i Iota) IsConstant() bool {
	return true
}
