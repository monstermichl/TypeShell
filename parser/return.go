package parser

import "github.com/monstermichl/typeshell/lexer"

type Return struct {
	token  lexer.Token
	Values []Expression
}

func (r Return) Token() lexer.Token {
	return r.token
}
