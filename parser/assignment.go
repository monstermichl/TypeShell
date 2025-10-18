package parser

import "github.com/monstermichl/typeshell/lexer"

type Assignment struct {
	Left          []Expression
	OperatorToken lexer.Token
	Right         []Expression
}

func (a Assignment) StatementType() StatementType {
	return STATEMENT_TYPE_ASSIGNMENT
}

func (a Assignment) Token() lexer.Token {
	return a.OperatorToken
}
