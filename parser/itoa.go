package parser

type Itoa struct {
	value Expression
}

func (o Itoa) ExprFn() {}

func (e Itoa) Value() Expression {
	return e.value
}
