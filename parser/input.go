package parser

type Input struct {
	prompt Expression
}

func (i Input) StatementType() StatementType {
	return STATEMENT_TYPE_INPUT
}

func (i Input) ValueType() ValueType {
	return NewValueType(TypeString{}, false)
}

func (i Input) IsConstant() bool {
	return false
}

func (i Input) Prompt() Expression {
	return i.prompt
}
