package parser

type Comparison struct {
	left     Expression
	operator CompareOperator
	right    Expression
}

func NewComparison(left Expression, operator CompareOperator, right Expression) Comparison {
	return Comparison{
		left,
		operator,
		right,
	}
}

func (c Comparison) Left() Expression {
	return c.left
}

func (c Comparison) Right() Expression {
	return c.right
}

func (c Comparison) Operator() CompareOperator {
	return c.operator
}
