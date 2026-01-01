package parser

import "github.com/monstermichl/typeshell/lexer"

type Continue struct {
	token lexer.Token
}

func (c Continue) Token() lexer.Token {
	return c.token
}
