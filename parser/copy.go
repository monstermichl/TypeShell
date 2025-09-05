package parser

type Copy struct {
	destination Variable
	source      Expression
}

func (c Copy) StatementType() StatementType {
	return STATEMENT_TYPE_COPY
}

func (c Copy) ValueType() ValueType {
	return NewValueType(DATA_TYPE_INTEGER, false)
}

func (c Copy) IsConstant() bool {
	return false
}

func (c Copy) Source() Expression {
	return c.source
}

func (c Copy) Destination() Variable {
	return c.destination
}
