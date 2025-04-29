package tests

import (
	"testing"
)

func TestDefineVariableSuccess(t *testing.T) {
	testDefineVariablesSuccess(t, transpileBatch)
}

func TestDefineSameVariableFail(t *testing.T) {
	testDefineSameVariableFail(t, transpileBatch)
}

func TestNoNewVariableFail(t *testing.T) {
	testNoNewVariableFail(t, transpileBatch)
}

func TestAssignSuccessful(t *testing.T) {
	testAssignSuccessful(t, transpileBatch)
}

func TestAssignToUndefinedFail(t *testing.T) {
	testAssignToUndefinedFail(t, transpileBatch)
}
