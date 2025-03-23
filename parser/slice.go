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
	dataType DataType
	index    Expression
}

func (s SliceEvaluation) StatementType() StatementType {
	return STATEMENT_TYPE_SLICE_EVALUATION
}

func (s SliceEvaluation) ValueType() ValueType {
	return ValueType{dataType: s.dataType, isSlice: true}
}

func (s SliceEvaluation) Index() Expression {
	return s.index
}

type SliceAssignment struct {
	index Expression
	value Expression
}

func (s SliceAssignment) StatementType() StatementType {
	return STATEMENT_TYPE_SLICE_ASSIGNMENT
}

func (s SliceAssignment) Index() Expression {
	return s.index
}

func (s SliceAssignment) Value() Expression {
	return s.value
}
