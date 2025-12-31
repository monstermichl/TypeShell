package parser

import "github.com/monstermichl/typeshell/lexer"

type Statement interface {
	Token() lexer.Token
}
