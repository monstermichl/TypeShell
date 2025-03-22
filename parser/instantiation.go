package parser

type Instantiation struct {
	valueType ValueType
	args      []Expression
}

func (i Instantiation) StatementType() StatementType {
	return STATEMENT_TYPE_INSTANTIATION
}

func (i Instantiation) ValueType() ValueType {
	return i.valueType
}
