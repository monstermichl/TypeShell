package tests

import (
	"testing"
)

func TestDefineSliceSuccess(t *testing.T) {
	testDefineSliceSuccess(t, transpileBash)
}

func TestSliceAssignValuesSuccess(t *testing.T) {
	testSliceAssignValuesSuccess(t, transpileBash)
}

func TestSliceAssignUndefinedSubscriptSuccess(t *testing.T) {
	testSliceAssignUndefinedSubscriptSuccess(t, transpileBash)
}

func TestSliceLengthSuccess(t *testing.T) {
	testSliceLengthSuccess(t, transpileBash)
}

func TestIterateSliceSuccess(t *testing.T) {
	testIterateSliceSuccess(t, transpileBash)
}

func TestCopySliceSuccess(t *testing.T) {
	testCopySliceSuccess(t, transpileBash)
}

func TestDefineSliceInFunctionSuccess(t *testing.T) {
	testDefineSliceInFunctionSuccess(t, transpileBash)
}

func TestSliceAssignValuesInFunctionSuccess(t *testing.T) {
	testSliceAssignValuesInFunctionSuccess(t, transpileBash)
}

func TestSliceAssignUndefinedSubscriptInFunctionSuccess(t *testing.T) {
	testSliceAssignUndefinedSubscriptInFunctionSuccess(t, transpileBash)
}

func TestSliceLengthInFunctionSuccess(t *testing.T) {
	testSliceLengthInFunctionSuccess(t, transpileBash)
}

func TestIterateSliceInFunctionSuccess(t *testing.T) {
	testIterateSliceInFunctionSuccess(t, transpileBash)
}

func TestCopySliceInFunctionSuccess(t *testing.T) {
	testCopySliceInFunctionSuccess(t, transpileBash)
}
