package tests

import (
	"testing"
)

func TestDefineVariableSuccess(t *testing.T) {
	testDefineVariablesSuccess(t, transpileBash)
}

func TestDefineSameVariableFail(t *testing.T) {
	testDefineSameVariableFail(t, transpileBash)
}

func TestNoNewVariableFail(t *testing.T) {
	testNoNewVariableFail(t, transpileBash)
}

func TestAssignSuccessful(t *testing.T) {
	testAssignSuccessful(t, transpileBash)
}

func TestAssignToUndefinedFail(t *testing.T) {
	testAssignToUndefinedFail(t, transpileBash)
}

func TestAssignFromFunctionSuccessful(t *testing.T) {
	testAssignFromFunctionSuccessful(t, transpileBash)
}

func TestDefineVariableInFunctionSuccess(t *testing.T) {
	testDefineVariablesInFunctionSuccess(t, transpileBash)
}

func TestDefineSameVariableInFunctionFail(t *testing.T) {
	testDefineSameVariableInFunctionFail(t, transpileBash)
}

func TestNoNewVariableInFunctionFail(t *testing.T) {
	testNoNewVariableInFunctionFail(t, transpileBash)
}

func TestAssignInFunctionSuccessful(t *testing.T) {
	testAssignInFunctionSuccessful(t, transpileBash)
}

func TestAssignToUndefinedInFunctionFail(t *testing.T) {
	testAssignToUndefinedInFunctionFail(t, transpileBash)
}

func TestAssignFromFunctionInFunctionSuccessful(t *testing.T) {
	testAssignFromFunctionInFunctionSuccessful(t, transpileBash)
}
