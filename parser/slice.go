package parser

import "github.com/monstermichl/typeshell/lexer"

type SliceType struct {
	Type Type
}

func (s SliceType) Token() lexer.Token {
	return s.Type.Token()
}

func (s SliceType) TypeFn() {}

type SliceInstantiation struct {
	t      Type
	values []Expression
}

func (s SliceInstantiation) ExprFn() {}

func (s SliceInstantiation) Values() []Expression {
	return s.values
}

type SliceEvaluation struct {
	value Expression
	index Expression
}

func (s SliceEvaluation) Value() Expression {
	return s.value
}

func (s SliceEvaluation) Index() Expression {
	return s.index
}

func (s SliceEvaluation) ExprFn() {}

type SliceAssignment struct {
	value      Expression
	index      Expression
	assignment Expression
}

func (s SliceAssignment) Value() Expression {
	return s.value
}

func (s SliceAssignment) Index() Expression {
	return s.index
}

func (s SliceAssignment) Assignment() Expression {
	return s.assignment
}
