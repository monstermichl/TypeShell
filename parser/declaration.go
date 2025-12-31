package parser

import "github.com/monstermichl/typeshell/lexer"

type ValueSpec struct {
	Names  []Expression
	Type   Type
	Values []Expression
}

type Declaration struct {
	Keyword        lexer.Token
	Specs          []ValueSpec
	OpeningBracket *lexer.Token
	ClosingBracket *lexer.Token
}

func (d Declaration) Token() lexer.Token {
	return d.Keyword
}
