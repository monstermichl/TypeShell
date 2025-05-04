package parser

type StringSubscript struct {
	value      Expression
	startIndex Expression
	endIndex   Expression
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

func (s StringSubscript) StartIndex() Expression {
	return s.startIndex
}

func (s StringSubscript) EndIndex() Expression {
	endIndex := s.endIndex

	// Automatically define end-index.
	if endIndex == nil {
		endIndex = s.startIndex
	}
	return endIndex
}
