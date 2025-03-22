package parser

type Expression interface {
	// An expression is a super-type of statement which results in a value.
	Statement
	ValueType() ValueType
}
