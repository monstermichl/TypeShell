package parser

import "github.com/monstermichl/typeshell/lexer"

type For struct {
	Init      Statement
	Condition Expression
	Post      Statement
	Body      Block
}

func (f For) StatementType() StatementType {
	return STATEMENT_TYPE_FOR
}

type ForRange struct {
	Key        Expression
	Value      Expression
	RangeToken lexer.Token
	X          Expression
	Body       Block
}

func (f ForRange) StatementType() StatementType {
	return STATEMENT_TYPE_FOR_RANGE
}
