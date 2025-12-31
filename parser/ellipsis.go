package parser

import "github.com/monstermichl/typeshell/lexer"

type Ellipsis struct {
	token lexer.Token
	Name  Expression
	Type  Expression
}

func (e Ellipsis) Token() lexer.Token {
	return e.token
}
