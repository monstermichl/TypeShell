package tests

import (
	"testing"
)

func TestDefineSliceSuccess(t *testing.T) {
	testDefineSliceSuccess(t, transpileBatch)
}

func TestSliceAssignValuesSuccess(t *testing.T) {
	testSliceAssignValuesSuccess(t, transpileBatch)
}

func TestSliceAssignUndefinedSubscriptSuccess(t *testing.T) {
	testSliceAssignUndefinedSubscriptSuccess(t, transpileBatch)
}

func TestSliceLengthSuccess(t *testing.T) {
	testSliceLengthSuccess(t, transpileBatch)
}

func TestIterateSliceSuccess(t *testing.T) {
	testIterateSliceSuccess(t, transpileBatch)
}

func TestReassignSliceSuccess(t *testing.T) {
	testReassignSliceSuccess(t, transpileBatch)
}

func TestCopySliceSuccess(t *testing.T) {
	testCopySliceSuccess(t, transpileBatch)
}

func TestDefineSliceInFunctionSuccess(t *testing.T) {
	testDefineSliceInFunctionSuccess(t, transpileBatch)
}

func TestSliceAssignValuesInFunctionSuccess(t *testing.T) {
	testSliceAssignValuesInFunctionSuccess(t, transpileBatch)
}

func TestSliceAssignUndefinedSubscriptInFunctionSuccess(t *testing.T) {
	testSliceAssignUndefinedSubscriptInFunctionSuccess(t, transpileBatch)
}

func TestSliceLengthInFunctionSuccess(t *testing.T) {
	testSliceLengthInFunctionSuccess(t, transpileBatch)
}

func TestIterateSliceInFunctionSuccess(t *testing.T) {
	testIterateSliceInFunctionSuccess(t, transpileBatch)
}

func TestReassignSliceInFunctionSuccess(t *testing.T) {
	testReassignSliceInFunctionSuccess(t, transpileBatch)
}

func TestCopySliceInFunctionSuccess(t *testing.T) {
	testCopySliceInFunctionSuccess(t, transpileBatch)
}

func TestSliceReturnedFromFunctionSuccess(t *testing.T) {
	testSliceReturnedFromFunctionSuccess(t, transpileBatch)
}

func TestComplexSliceOperationsSuccess(t *testing.T) {
	testComplexSliceOperationsSuccess(t, transpileBatch)
}
