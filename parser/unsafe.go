package parser

type Unsafe struct {
	code StringLiteral
	args []Expression
}

func (u Unsafe) StatementType() StatementType {
	return STATEMENT_TYPE_UNSAFE
}

func (u Unsafe) Code() StringLiteral {
	return u.code
}

func (u Unsafe) Args() []Expression {
	return u.args
}
