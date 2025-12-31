package parser

import "github.com/monstermichl/typeshell/lexer"

type FunctionDefinition struct {
	token       lexer.Token
	Receiver    *Field
	Name        Expression
	Params      FieldList
	ReturnTypes []Expression
}

func (e FunctionDefinition) Token() lexer.Token {
	return e.token
}
