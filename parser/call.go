package parser

type Call interface {
	Expression
	Call() Expression
	Args() []Expression
	SetArgs([]Expression)
}
