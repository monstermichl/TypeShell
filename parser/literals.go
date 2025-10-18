package parser

import "github.com/monstermichl/typeshell/lexer"

type BooleanLiteral struct {
	Value bool
	token lexer.Token
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

func (l BooleanLiteral) Token() lexer.Token {
	return l.token
}

func (l BooleanLiteral) IsConstant() bool {
	return true
}

type IntegerLiteral struct {
	Value int
	token lexer.Token
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

func (l IntegerLiteral) Token() lexer.Token {
	return l.token
}

func (l IntegerLiteral) IsConstant() bool {
	return true
}

type StringLiteral struct {
	Value string
	token lexer.Token
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

func (l StringLiteral) Token() lexer.Token {
	return l.token
}

func (l StringLiteral) IsConstant() bool {
	return true
}
