package parser

type Unsafe struct {
	code StringLiteral
	args []Expression
}

func (u Unsafe) Code() StringLiteral {
	return u.code
}

func (u Unsafe) Args() []Expression {
	return u.args
}
