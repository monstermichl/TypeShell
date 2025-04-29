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

func TestAssignFromFunctionSuccessful(t *testing.T) {
	testAssignFromFunctionSuccessful(t, transpileBatch)
}

func TestDefineVariableInFunctionSuccess(t *testing.T) {
	testDefineVariablesInFunctionSuccess(t, transpileBatch)
}

func TestDefineSameVariableInFunctionFail(t *testing.T) {
	testDefineSameVariableInFunctionFail(t, transpileBatch)
}

func TestNoNewVariableInFunctionFail(t *testing.T) {
	testNoNewVariableInFunctionFail(t, transpileBatch)
}

func TestAssignInFunctionSuccessful(t *testing.T) {
	testAssignInFunctionSuccessful(t, transpileBatch)
}

func TestAssignToUndefinedInFunctionFail(t *testing.T) {
	testAssignToUndefinedInFunctionFail(t, transpileBatch)
}

func TestAssignFromFunctionInFunctionSuccessful(t *testing.T) {
	testAssignFromFunctionInFunctionSuccessful(t, transpileBatch)
}
