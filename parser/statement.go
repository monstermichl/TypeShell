package parser

import "github.com/monstermichl/typeshell/lexer"

type Statement interface {
	StatementType() StatementType
	Token() lexer.Token
}
