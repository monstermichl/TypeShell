package parser

import "github.com/monstermichl/typeshell/lexer"

type CaseClause struct {
	token   lexer.Token
	List    []Expression
	Body    Block
	Default bool
}

func (c CaseClause) Token() lexer.Token {
	return c.token
}

type Switch struct {
	token lexer.Token
	Tag   Expression
	Body  Block
}

func (s Switch) Token() lexer.Token {
	return s.token
}
