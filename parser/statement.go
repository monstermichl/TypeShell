package parser

type Statement interface {
	StatementType() StatementType
}
