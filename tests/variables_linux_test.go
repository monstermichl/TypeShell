package tests

import (
	"testing"
)

func TestDefineVariableSuccess(t *testing.T) {
	testDefineVariablesSuccess(t, transpileBash)
}

func TestDefineSliceVariable(t *testing.T) {
	testDefineSliceVariable(t, transpileBash)
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
