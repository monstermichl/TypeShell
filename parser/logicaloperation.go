package parser

type LogicalOperation struct {
	left     Expression
	operator LogicalOperator
	right    Expression
}

func (l LogicalOperation) StatementType() StatementType {
	return STATEMENT_TYPE_LOGICAL_OPERATION
}

func (l LogicalOperation) ValueType() ValueType {
	return VALUE_TYPE_BOOLEAN
}

func (l LogicalOperation) Left() Expression {
	return l.left
}

func (l LogicalOperation) Operator() LogicalOperator {
	return l.operator
}

func (l LogicalOperation) Right() Expression {
	return l.right
}
