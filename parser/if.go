package parser

import "github.com/monstermichl/typeshell/lexer"

type If struct {
	token     lexer.Token
	Condition Expression
	Body      Block
	Else      Statement
}

func (i If) Token() lexer.Token {
	return i.token
}
