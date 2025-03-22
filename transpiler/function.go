package transpiler

import "github.com/monstermichl/typeshell/parser"

type function struct {
	name      string
	valueType valueType
	params    []parser.Variable
	body      []parser.Statement
}

func (e function) Name() string {
	return e.name
}

func (e function) ValueType() valueType {
	return e.valueType
}

func (e function) Body() []parser.Statement {
	return e.body
}
