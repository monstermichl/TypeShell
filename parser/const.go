package parser

type Const struct {
	name      string
	valueType ValueType
	layer     int
	public    bool
}

func NewConst(name string, valueType ValueType, layer int, public bool) Const {
	return Const{
		name,
		valueType,
		layer,
		public,
	}
}

func (c Const) Name() string {
	return c.name
}

func (c Const) ValueType() ValueType {
	return c.valueType
}

func (c Const) Layer() int {
	return c.layer
}

func (c Const) IsConstant() bool {
	return true
}

func (c *Const) SetValueType(valueType ValueType) {
	c.valueType = valueType
}

func (c Const) Global() bool {
	return c.layer == 0
}

func (c Const) Public() bool {
	return c.public
}

type ConstDefinition struct {
	constants []Const
	values    []Expression
}

func (c ConstDefinition) StatementType() StatementType {
	return STATEMENT_TYPE_CONST_DEFINITION
}

func (c ConstDefinition) AssignmentType() AssignmentType {
	return ASSIGNMENT_TYPE_VALUE
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
