package transpiler

type valueType string

const (
	VALUE_TYPE_UNKNOWN valueType = "unknown"
	VALUE_TYPE_VOID    valueType = "void"
	VALUE_TYPE_BOOLEAN valueType = "bool"
	VALUE_TYPE_INTEGER valueType = "int"
	VALUE_TYPE_STRING  valueType = "string"
)
