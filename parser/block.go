package parser

type Block interface {
	Body() []Statement
}
