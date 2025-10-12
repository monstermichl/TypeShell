package parser

import "github.com/monstermichl/typeshell/lexer"

type Index struct {
	X            Expression
	Index        Expression
	LeftBracket  *lexer.Token
	RightBracket *lexer.Token
}

func NewIndex(x Expression, index Expression, leftBracket *lexer.Token, rightBracket *lexer.Token) Index {
	return Index{x, index, leftBracket, rightBracket}
}

func (i Index) StatementType() StatementType {
	return STATEMENT_TYPE_INDEX
}

func (i Index) ValueType() ValueType {
	return NewValueType(NewTypeUnknown(), false) // TODO: Remove. Type is not relevant in parser.
}

func (i Index) IsConstant() bool {
	return false
}
