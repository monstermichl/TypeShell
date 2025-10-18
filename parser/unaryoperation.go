package parser

import "github.com/monstermichl/typeshell/lexer"

type UnaryOperation struct {
	X             Expression
	OperatorToken lexer.Token
}

func NewUnaryOperation(x Expression, operatorToken lexer.Token) UnaryOperation {
	return UnaryOperation{x, operatorToken}
}

func (b UnaryOperation) StatementType() StatementType {
	return STATEMENT_TYPE_UNARY_OPERATION
}

func (b UnaryOperation) ValueType() ValueType {
	return NewValueType(NewTypeUnknown(), false) // TODO: Remove. Type is not relevant in parser.
}

func (b UnaryOperation) Token() lexer.Token {
	return b.OperatorToken
}

func (b UnaryOperation) IsConstant() bool {
	return b.X.IsConstant()
}
