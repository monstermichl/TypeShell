package tests

import (
	"testing"
)

func TestVoidFunctionSuccess(t *testing.T) {
	testVoidFunctionSuccess(t, transpileBatch)
}

func TestSingleReturnValueFunctionSuccess(t *testing.T) {
	testSingleReturnValueFunctionSuccess(t, transpileBatch)
}

func TestMultiReturnValueFunctionSuccess(t *testing.T) {
	testMultiReturnValueFunctionSuccess(t, transpileBatch)
}

func TestSingleParamFunctionSuccess(t *testing.T) {
	testSingleParamFunctionSuccess(t, transpileBatch)
}

func TestMultiParamFunctionSuccess(t *testing.T) {
	testMultiParamFunctionSuccess(t, transpileBatch)
}

func TestSliceParamFunctionSuccess(t *testing.T) {
	testSliceParamFunctionSuccess(t, transpileBatch)
}

func TestCallFunctionFromFunctionSuccess(t *testing.T) {
	testCallFunctionFromFunctionSuccess(t, transpileBatch)
}
