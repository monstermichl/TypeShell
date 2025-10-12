package parser

import "github.com/monstermichl/typeshell/lexer"

type BooleanLiteral struct {
	Value bool
	Token lexer.Token
}

func NewBooleanLiteral(value bool, token lexer.Token) BooleanLiteral {
	return BooleanLiteral{value, token}
}

func (l BooleanLiteral) StatementType() StatementType {
	return STATEMENT_TYPE_BOOL_LITERAL
}

func (l BooleanLiteral) ValueType() ValueType {
	return NewValueType(NewTypeBool(), false)
}

func (l BooleanLiteral) IsConstant() bool {
	return true
}

type IntegerLiteral struct {
	Value int
	Token lexer.Token
}

func NewIntegerLiteral(value int, token lexer.Token) IntegerLiteral {
	return IntegerLiteral{value, token}
}

func (l IntegerLiteral) StatementType() StatementType {
	return STATEMENT_TYPE_INT_LITERAL
}

func (l IntegerLiteral) ValueType() ValueType {
	return NewValueType(NewTypeInt(), false)
}

func (l IntegerLiteral) IsConstant() bool {
	return true
}

type StringLiteral struct {
	Value string
	Token lexer.Token
}

func NewStringLiteral(value string, token lexer.Token) StringLiteral {
	return StringLiteral{value, token}
}

func (l StringLiteral) StatementType() StatementType {
	return STATEMENT_TYPE_STRING_LITERAL
}

func (l StringLiteral) ValueType() ValueType {
	return NewValueType(NewTypeString(), false)
}

func (l StringLiteral) IsConstant() bool {
	return true
}
