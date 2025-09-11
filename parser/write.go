package parser

type Write struct {
	path   Expression
	data   Expression
	append Expression
}

func (w Write) StatementType() StatementType {
	return STATEMENT_TYPE_WRITE
}

func (w Write) ValueType() ValueType {
	return NewValueType(TypeError{}, false)
}

func (w Write) IsConstant() bool {
	return false
}

func (w Write) Path() Expression {
	return w.path
}

func (w Write) Data() Expression {
	return w.data
}

func (w Write) Append() Expression {
	return w.append
}
