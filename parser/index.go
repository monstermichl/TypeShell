package parser

import "github.com/monstermichl/typeshell/lexer"

type Index struct {
	X            Expression
	Index        Expression
	LeftBracket  *lexer.Token
	RightBracket *lexer.Token
}

func NewIndex(x Expression, index Expression, leftBracket *lexer.Token, rightBracket *lexer.Token) Index {
	return Index{x, index, leftBracket, rightBracket}
}

func (i Index) Token() lexer.Token {
	return i.X.Token()
}

func (i Index) ExprFn() {}
