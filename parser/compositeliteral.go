package parser

import "github.com/monstermichl/typeshell/lexer"

type CompositeLiteral struct {
	Elements []Expression
	Type     Type
	token    lexer.Token
}

func (c CompositeLiteral) Token() lexer.Token {
	return c.token
}

func (c CompositeLiteral) ExprFn() {}
