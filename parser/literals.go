package parser

import "github.com/monstermichl/typeshell/lexer"

type BooleanLiteral struct {
	Value bool
	token lexer.Token
}

func NewBooleanLiteral(value bool, token lexer.Token) BooleanLiteral {
	return BooleanLiteral{value, token}
}

func (l BooleanLiteral) Token() lexer.Token {
	return l.token
}

func (l BooleanLiteral) ExprFn() {}

type IntegerLiteral struct {
	Value int
	token lexer.Token
}

func NewIntegerLiteral(value int, token lexer.Token) IntegerLiteral {
	return IntegerLiteral{value, token}
}

func (l IntegerLiteral) Token() lexer.Token {
	return l.token
}

func (l IntegerLiteral) ExprFn() {}

type StringLiteral struct {
	Value string
	token lexer.Token
}

func NewStringLiteral(value string, token lexer.Token) StringLiteral {
	return StringLiteral{value, token}
}

func (l StringLiteral) Token() lexer.Token {
	return l.token
}

func (l StringLiteral) ExprFn() {}
