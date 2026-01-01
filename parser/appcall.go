package parser

import "github.com/monstermichl/typeshell/lexer"

type AppCall struct {
	token     lexer.Token
	Name      Expression
	Arguments []Expression
	Next      *AppCall
}

func (a AppCall) Token() lexer.Token {
	return a.token
}

func (a AppCall) Call() Expression {
	return a.Name
}

func (a AppCall) Args() []Expression {
	return a.Arguments
}

func (a *AppCall) SetArgs(args []Expression) {
	a.Arguments = args
}

func (a AppCall) ExprFn() {}
