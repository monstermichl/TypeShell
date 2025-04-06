package parser

type UnaryOperation struct {
	expr      Expression
	operator  UnaryOperator
	valueType ValueType
}

func (b UnaryOperation) StatementType() StatementType {
	return STATEMENT_TYPE_UNARY_OPERATION
}

func (b UnaryOperation) ValueType() ValueType {
	return b.valueType
}

func (b UnaryOperation) Expression() Expression {
	return b.expr
}

func (b UnaryOperation) Operator() UnaryOperator {
	return b.operator
}
