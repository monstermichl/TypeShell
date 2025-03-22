package parser

type BooleanLiteral struct {
	value bool
}

func (l BooleanLiteral) StatementType() StatementType {
	return STATEMENT_TYPE_BOOL_LITERAL
}

func (l BooleanLiteral) ValueType() ValueType {
	return VALUE_TYPE_BOOLEAN
}

func (l BooleanLiteral) Value() bool {
	return l.value
}

type IntegerLiteral struct {
	value int
}

func (l IntegerLiteral) StatementType() StatementType {
	return STATEMENT_TYPE_INT_LITERAL
}

func (l IntegerLiteral) ValueType() ValueType {
	return VALUE_TYPE_INTEGER
}

func (l IntegerLiteral) Value() int {
	return l.value
}

type StringLiteral struct {
	value string
}

func (l StringLiteral) StatementType() StatementType {
	return STATEMENT_TYPE_STRING_LITERAL
}

func (l StringLiteral) ValueType() ValueType {
	return VALUE_TYPE_STRING
}

func (l StringLiteral) Value() string {
	return l.value
}
