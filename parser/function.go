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

type FunctionCall struct {
	token               lexer.Token
	openingBracketToken lexer.Token
	closingBracketToken lexer.Token
	Func                Expression
	Arguments           []Expression
}

func (f FunctionCall) Token() lexer.Token {
	return f.token
}

func (f FunctionCall) Call() Expression {
	return f.Func
}

func (f FunctionCall) Args() []Expression {
	return f.Arguments
}

func (f *FunctionCall) SetArgs(args []Expression) {
	f.Arguments = args
}

func (f FunctionCall) ExprFn() {}
