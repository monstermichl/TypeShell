package parser

import "github.com/monstermichl/typeshell/lexer"

type Import struct {
	token lexer.Token
	Name  Expression
	Path  Expression
}

func (i Import) Token() lexer.Token {
	return i.token
}

type Imports struct {
	token   lexer.Token
	Imports []Import
}

func (i Imports) Token() lexer.Token {
	return i.token
}
