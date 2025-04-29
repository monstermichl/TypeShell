package parser

type SliceInstantiation struct {
	dataType DataType
	values   []Expression
}

func (s SliceInstantiation) StatementType() StatementType {
	return STATEMENT_TYPE_SLICE_INSTANTIATION
}

func (s SliceInstantiation) ValueType() ValueType {
	return ValueType{dataType: s.dataType, isSlice: true}
}

func (s SliceInstantiation) Values() []Expression {
	return s.values
}

type SliceEvaluation struct {
	value    Expression
	index    Expression
	dataType DataType
}

func (s SliceEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_SLICE_EVALUATION
}

func (s SliceEvaluation) Value() Expression {
	return s.value
}

func (s SliceEvaluation) Index() Expression {
	return s.index
}

func (s SliceEvaluation) ValueType() ValueType {
	return ValueType{dataType: s.dataType}
}

type SliceAssignment struct {
	Variable
	index Expression
	value Expression
}

func (s SliceAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_SLICE_ASSIGNMENT
}

func (s SliceAssignment) Name() string {
	return s.name
}

func (s SliceAssignment) Index() Expression {
	return s.index
}

func (s SliceAssignment) Value() Expression {
	return s.value
}
