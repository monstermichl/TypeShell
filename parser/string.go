package parser

type StringSubscript struct {
	Variable
	index Expression
}

func (s StringSubscript) StatementType() StatementType {
	return STATEMENT_TYPE_STRING_SUBSCRIPT
}

func (s StringSubscript) ValueType() ValueType {
	return NewValueType(DATA_TYPE_STRING, false)
}

func (s StringSubscript) Index() Expression {
	return s.index
}

func (s StringSubscript) Global() bool {
	return s.global
}
