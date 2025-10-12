package parser

import "github.com/monstermichl/typeshell/lexer"

type Block struct {
	Statements     []Statement
	OpeningBracket *lexer.Token
	ClosingBracket *lexer.Token
}
