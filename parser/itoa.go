package parser

type Itoa struct {
	value Expression
}

func (e Itoa) StatementType() StatementType {
	return STATEMENT_TYPE_ITOA
}

func (e Itoa) ValueType() ValueType {
	return NewValueType(TypeString{}, false)
}

func (o Itoa) IsConstant() bool {
	return false
}

func (e Itoa) Value() Expression {
	return e.value
}
