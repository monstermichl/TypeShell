package parser

import "github.com/monstermichl/typeshell/lexer"

type Field struct {
	token lexer.Token
	Name  Expression
	Type  Statement
}

func (p Field) Token() lexer.Token {
	return p.token
}

type FieldList struct {
	token          lexer.Token
	Fields         []Field
	OpeningBracket *lexer.Token
	ClosingBracket *lexer.Token
}

func (p FieldList) Token() lexer.Token {
	return p.token
}
