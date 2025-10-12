package parser

import "github.com/monstermichl/typeshell/lexer"

type Group struct {
	X              Expression
	OpeningBracket *lexer.Token
	ClosingBracket *lexer.Token
}

func NewGroup(x Expression, openingBracket *lexer.Token, closingBracket *lexer.Token) Group {
	return Group{x, openingBracket, closingBracket}
}

func (e Group) StatementType() StatementType {
	return STATEMENT_TYPE_GROUP
}

func (e Group) ValueType() ValueType {
	return e.X.ValueType()
}

func (e Group) IsConstant() bool {
	return e.X.IsConstant()
}
