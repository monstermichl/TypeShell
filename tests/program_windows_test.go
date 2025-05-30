package tests

import (
	"testing"
)

func TestComplexProgram1Success(t *testing.T) {
	testComplexProgram1Success(t, transpileBatch)
}

func TestComplexProgram2Success(t *testing.T) {
	testComplexProgram2Success(t, transpileBatch)
}

func TestComplexProgram3Success(t *testing.T) {
	testComplexProgram3Success(t, transpileBatch)
}

func TestComplexProgram4Success(t *testing.T) {
	testComplexProgram4Success(t, transpileBatch)
}
