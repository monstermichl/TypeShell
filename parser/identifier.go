package parser

import "github.com/monstermichl/typeshell/lexer"

type Identifier struct {
	Name  string
	token lexer.Token
}

func NewIdentifier(name string, token lexer.Token) Identifier {
	return Identifier{name, token}
}

func (i Identifier) Token() lexer.Token {
	return i.token
}

func (i Identifier) ExprFn() {}
func (i Identifier) TypeFn()  {}
