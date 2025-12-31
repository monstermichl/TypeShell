package parser

type StringSubscript struct {
	value      Expression
	startIndex Expression
	endIndex   Expression
}

func (s StringSubscript) ExprFn() {}

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
