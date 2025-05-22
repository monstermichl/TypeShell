package parser

type Call interface {
	Expression
	Name() string
	Args() []Expression
	ReturnTypes() []ValueType
}
