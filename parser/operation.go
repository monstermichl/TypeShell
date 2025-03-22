package parser

type Operation interface {
	Left() Expression
	Operator() string
	Right() Expression
}
