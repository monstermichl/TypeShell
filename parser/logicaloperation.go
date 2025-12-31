package parser

type LogicalOperation struct {
	left     Expression
	operator LogicalOperator
	right    Expression
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
