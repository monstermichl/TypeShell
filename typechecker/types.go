package typechecker

type Type interface {
	IsSlice() bool
	String() string
}
