package parser

type Input struct {
	prompt Expression
}

func (i Input) StatementType() StatementType {
	return STATEMENT_TYPE_INPUT
}

func (i Input) ValueType() ValueType {
	return ValueType{dataType: DATA_TYPE_STRING}
}

func (i Input) Prompt() Expression {
	return i.prompt
}
