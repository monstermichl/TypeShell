package parser

type BooleanLiteral struct {
	value bool
}

func (l BooleanLiteral) StatementType() StatementType {
	return STATEMENT_TYPE_BOOL_LITERAL
}

func (l BooleanLiteral) ValueType() ValueType {
	return NewValueType(TypeBool{}, false)
}

func (l BooleanLiteral) IsConstant() bool {
	return true
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
	return NewValueType(TypeInt{}, false)
}

func (l IntegerLiteral) IsConstant() bool {
	return true
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
	return NewValueType(TypeString{}, false)
}

func (l StringLiteral) IsConstant() bool {
	return true
}

func (l StringLiteral) Value() string {
	return l.value
}
