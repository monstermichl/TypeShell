package tests

import (
	"testing"
)

func TestDefineConstantsSuccess(t *testing.T) {
	testDefineConstantsSuccess(t, transpileBash)
}

func TestDefineConstantsInFunctionSuccess(t *testing.T) {
	testDefineConstantsInFunctionSuccess(t, transpileBash)
}

func TestDefineConstantsGroupedSuccess(t *testing.T) {
	testDefineConstantsGroupedSuccess(t, transpileBash)
}

func TestDefineConstantsMissingValueFail(t *testing.T) {
	testDefineConstantsMissingValueFail(t, transpileBash)
}

func TestDefineSameConstantFail(t *testing.T) {
	testDefineSameConstantFail(t, transpileBash)
}

func TestAssignFail(t *testing.T) {
	testAssignFail(t, transpileBash)
}

func TestAssignFromFunctionFail(t *testing.T) {
	testAssignFromFunctionFail(t, transpileBash)
}
