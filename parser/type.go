package parser

import "github.com/monstermichl/typeshell/lexer"

type Type interface {
	Statement
	TypeFn()
}

type TypeDeclaration struct {
	token         lexer.Token
	Name          Expression
	AssignOpToken *lexer.Token
	Type          Type
}

func (t TypeDeclaration) Token() lexer.Token {
	return t.token
}

type TypeDefinition struct {
	value Expression
}

func (t TypeDefinition) Value() Expression {
	return t.value
}
