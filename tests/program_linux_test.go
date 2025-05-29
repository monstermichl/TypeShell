package tests

import (
	"testing"
)

func TestComplexProgram1Success(t *testing.T) {
	testComplexProgram1Success(t, transpileBash)
}

func TestComplexProgram2Success(t *testing.T) {
	testComplexProgram2Success(t, transpileBash)
}

func TestComplexProgram3Success(t *testing.T) {
	testComplexProgram3Success(t, transpileBash)
}
