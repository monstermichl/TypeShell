package parser

type Group struct {
	child Expression
}

func (e Group) StatementType() StatementType {
	return STATEMENT_TYPE_GROUP
}

func (e Group) ValueType() ValueType {
	return e.Child().ValueType()
}

func (e Group) IsConstant() bool {
	return e.Child().IsConstant()
}

func (e Group) Child() Expression {
	return e.child
}
