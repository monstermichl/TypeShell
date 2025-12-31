package parser

import "github.com/monstermichl/typeshell/lexer"

type UnaryOperation struct {
	X             Expression
	OperatorToken lexer.Token
}

func NewUnaryOperation(x Expression, operatorToken lexer.Token) UnaryOperation {
	return UnaryOperation{x, operatorToken}
}

func (b UnaryOperation) Token() lexer.Token {
	return b.OperatorToken
}

func (b UnaryOperation) ExprFn() { b.X.ExprFn() }
