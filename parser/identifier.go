package parser

import "github.com/monstermichl/typeshell/lexer"

type Identifier struct {
	Name  string
	Token lexer.Token
}

func NewIdentifier(name string, token lexer.Token) Identifier {
	return Identifier{name, token}
}

func (i Identifier) StatementType() StatementType {
	return STATEMENT_TYPE_IDENTIFIER
}

func (i Identifier) ValueType() ValueType {
	return NewValueType(NewTypeUnknown(), false) // TODO: Remove. Type is not relevant in parser.
}

func (i Identifier) IsConstant() bool {
	return false
}
