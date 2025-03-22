package transpiler

type evaluationValue struct {
	valueType valueType
	value     any
}

func (v evaluationValue) ValueType() valueType {
	return v.valueType
}

func (v evaluationValue) Value() any {
	return v.value
}
