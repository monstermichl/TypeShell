package parser

import "github.com/monstermichl/typeshell/lexer"

type StructField struct {
	Names []Expression
	Type  Type
}

type StructDeclaration struct {
	keywordToken lexer.Token
	Fields       []StructField
}

func (d StructDeclaration) Token() lexer.Token {
	return d.keywordToken
}

func (d StructDeclaration) TypeFn() {}

type StructValue struct {
	StructField
	value Expression
}

type StructInitialization struct {
	t      Type
	values []StructValue
}

type StructAssignment struct {
	value      Expression
	assignment StructValue
}

type StructEvaluation struct {
	value Expression
	field StructField
}
