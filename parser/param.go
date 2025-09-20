package parser

type Param struct {
	Variable
	pointer bool
}

func NewParam(name string, valueType ValueType, layer int, public bool, pointer bool) Param {
	return Param{
		Variable: Variable{
			name,
			valueType,
			layer,
			public,
		},
		pointer: pointer,
	}
}

func (p Param) Pointer() bool {
	return p.pointer
}
