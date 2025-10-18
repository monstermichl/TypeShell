package parser

import "github.com/monstermichl/typeshell/lexer"

type Expression interface {
	// An expression is a super-type of statement which results in a value.
	Statement
	ValueType() ValueType
	Token() lexer.Token
	IsConstant() bool
}
