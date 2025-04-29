package parser

type StringSubscript struct {
	value Expression
	index Expression
}

func (s StringSubscript) StatementType() StatementType {
	return STATEMENT_TYPE_STRING_SUBSCRIPT
}

func (s StringSubscript) ValueType() ValueType {
	return NewValueType(DATA_TYPE_STRING, false)
}

func (s StringSubscript) Value() Expression {
	return s.value
}

func (s StringSubscript) Index() Expression {
	return s.index
}
