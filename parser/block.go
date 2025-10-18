package parser

import "github.com/monstermichl/typeshell/lexer"

type Block struct {
	token          lexer.Token
	Statements     []Statement
	OpeningBracket *lexer.Token
	ClosingBracket *lexer.Token
}

func (b Block) StatementType() StatementType {
	return STATEMENT_TYPE_BLOCK
}

func (b Block) Token() lexer.Token {
	return b.token
}
