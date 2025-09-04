package parser

type NamedValue interface {
	Name() string
	ValueType() ValueType
	Global() bool
	Public() bool
	IsConstant() bool
}
