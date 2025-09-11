package parser

type Exists struct {
	path Expression
}

func (e Exists) StatementType() StatementType {
	return STATEMENT_TYPE_EXISTS
}

func (e Exists) ValueType() ValueType {
	return NewValueType(TypeBool{}, false)
}

func (e Exists) IsConstant() bool {
	return false
}

func (e Exists) Path() Expression {
	return e.path
}
