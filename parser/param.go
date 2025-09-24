package parser

type Param struct {
	Variable
	pointer  bool
	optional bool
}

func NewParam(name string, valueType ValueType, layer int, public bool, pointer bool, optional bool) Param {
	return Param{
		Variable: Variable{
			name,
			valueType,
			layer,
			public,
		},
		pointer:  pointer,
		optional: optional,
	}
}

func (p Param) Pointer() bool {
	return p.pointer
}

func (p Param) Optional() bool {
	return p.optional
}
