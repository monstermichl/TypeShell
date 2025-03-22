package parser

type Call interface {
	Expression
	Args() []Expression
}
