package parser

type BinaryOperation struct {
	left      Expression
	operator  BinaryOperator
	right     Expression
}

func (b BinaryOperation) StatementType() StatementType {
	return STATEMENT_TYPE_BINARY_OPERATION
}

func (b BinaryOperation) ValueType() ValueType {
	return b.left.ValueType()
}

func (b BinaryOperation) Left() Expression {
	return b.left
}

func (b BinaryOperation) Right() Expression {
	return b.right
}

func (b BinaryOperation) Operator() BinaryOperator {
	return b.operator
}
