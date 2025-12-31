package parser

type Input struct {
	prompt Expression
}

func (i Input) ExprFn() {}

func (i Input) Prompt() Expression {
	return i.prompt
}
