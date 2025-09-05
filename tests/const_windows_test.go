package tests

import (
	"testing"
)

func TestDefineConstantsSuccess(t *testing.T) {
	testDefineConstantsSuccess(t, transpileBatch)
}

func TestDefineConstantsInFunctionSuccess(t *testing.T) {
	testDefineConstantsInFunctionSuccess(t, transpileBatch)
}

func TestDefineConstantsGroupedSuccess(t *testing.T) {
	testDefineConstantsGroupedSuccess(t, transpileBatch)
}

func TestDefineConstantsMissingValueFail(t *testing.T) {
	testDefineConstantsMissingValueFail(t, transpileBatch)
}

func TestDefineSameConstantFail(t *testing.T) {
	testDefineSameConstantFail(t, transpileBatch)
}

func TestAssignFail(t *testing.T) {
	testAssignFail(t, transpileBatch)
}

func TestAssignFromFunctionFail(t *testing.T) {
	testAssignFromFunctionFail(t, transpileBatch)
}
