package parser

import "github.com/monstermichl/typeshell/lexer"

type Selector struct {
	X        Expression
	Selector string
}

func NewSelector(x Expression, name string) Selector {
	return Selector{x, name}
}

func (s Selector) Token() lexer.Token {
	return s.X.Token()
}

func (s Selector) ExprFn() {}
