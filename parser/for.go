package parser

import "github.com/monstermichl/typeshell/lexer"

type For struct {
	token     lexer.Token
	Init      Statement
	Condition Expression
	Post      Statement
	Body      Block
}

func (f For) StatementType() StatementType {
	return STATEMENT_TYPE_FOR
}

func (f For) Token() lexer.Token {
	return f.token
}

type ForRange struct {
	token      lexer.Token
	Key        Expression
	Value      Expression
	RangeToken lexer.Token
	X          Expression
	Body       Block
}

func (f ForRange) StatementType() StatementType {
	return STATEMENT_TYPE_FOR_RANGE
}

func (f ForRange) Token() lexer.Token {
	return f.token
}
