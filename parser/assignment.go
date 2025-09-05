package parser

type Assignment interface {
	Statement
	AssignmentType() AssignmentType
}
