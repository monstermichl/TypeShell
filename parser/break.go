package parser

import "github.com/monstermichl/typeshell/lexer"

type Break struct {
	token lexer.Token
}

func (b Break) Token() lexer.Token {
	return b.token
}
