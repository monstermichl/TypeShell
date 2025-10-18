package parser

import "github.com/monstermichl/typeshell/lexer"

type Selector struct {
	X        Expression
	Selector string
}

func NewSelector(x Expression, name string) Selector {
	return Selector{x, name}
}

func (s Selector) StatementType() StatementType {
	return STATEMENT_TYPE_SELECTOR
}

func (s Selector) ValueType() ValueType {
	return NewValueType(NewTypeUnknown(), false) // TODO: Remove. Type is not relevant in parser.
}

func (s Selector) Token() lexer.Token {
	return s.X.Token()
}

func (s Selector) IsConstant() bool {
	return false
}
