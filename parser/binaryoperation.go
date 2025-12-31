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

func (b BinaryOperation) Token() lexer.Token {
	return b.OperatorToken
}

func (b BinaryOperation) ExprFn() {}
