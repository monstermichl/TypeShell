package parser

import "github.com/monstermichl/typeshell/lexer"

type KeyValue struct {
	token lexer.Token
	Key   Expression
	Value Expression
}

func (k KeyValue) Token() lexer.Token {
	return k.token
}

func (k KeyValue) ExprFn() {}
