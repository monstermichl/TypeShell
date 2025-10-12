package parser

import "github.com/monstermichl/typeshell/lexer"

type BinaryOperation struct {
	Left          Expression
	OperatorToken lexer.Token
	Right         Expression
}

func NewBinaryOperation(left Expression, operatorToken lexer.Token, right Expression) BinaryOperation {
	return BinaryOperation{left, operatorToken, right}
}

func (b BinaryOperation) StatementType() StatementType {
	return STATEMENT_TYPE_BINARY_OPERATION
}

func (b BinaryOperation) ValueType() ValueType {
	return NewValueType(NewTypeUnknown(), false) // TODO: Remove. Type is not relevant in parser.
}

func (b BinaryOperation) IsConstant() bool {
	return b.Left.IsConstant() && b.Right.IsConstant()
}
