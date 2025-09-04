package parser

type Const struct {
	name      string
	valueType ValueType
	global    bool
	public    bool
}

func (c Const) Name() string {
	return c.name
}

func (c Const) ValueType() ValueType {
	return c.valueType
}

func (c Const) Global() bool {
	return c.global
}

func (c Const) Public() bool {
	return c.public
}

func (c Const) IsConstant() bool {
	return true
}

type ConstDefinition struct {
	constants []Const
	values    []Expression
}

func (c ConstDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_CONST_DEFINITION
}

func (c ConstDefinition) Constants() []Const {
	return c.constants
}

func (c ConstDefinition) Values() []Expression {
	return c.values
}

type ConstEvaluation struct {
	Const
}

func (c ConstEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_CONST_EVALUATION
}
