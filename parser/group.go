package parser

import "github.com/monstermichl/typeshell/lexer"

type Group struct {
	X              Expression
	OpeningBracket lexer.Token
	ClosingBracket *lexer.Token
}

func NewGroup(x Expression, openingBracket lexer.Token, closingBracket *lexer.Token) Group {
	return Group{x, openingBracket, closingBracket}
}

func (e Group) Token() lexer.Token {
	return e.OpeningBracket
}

func (e Group) ExprFn() { e.X.ExprFn() }
