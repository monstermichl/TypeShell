package parser

import "github.com/monstermichl/typeshell/lexer"

type Selector struct {
	X        Expression
	Selector Expression
}

func NewSelector(x Expression, selector Expression) Selector {
	return Selector{x, selector}
}

func (s Selector) Token() lexer.Token {
	return s.X.Token()
}

func (s Selector) ExprFn() {}
func (s Selector) TypeFn() {}
