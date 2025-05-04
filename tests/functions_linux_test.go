package tests

import (
	"testing"
)

func TestVoidFunctionSuccess(t *testing.T) {
	testVoidFunctionSuccess(t, transpileBash)
}

func TestSingleReturnValueFunctionSuccess(t *testing.T) {
	testSingleReturnValueFunctionSuccess(t, transpileBash)
}

func TestMultiReturnValueFunctionSuccess(t *testing.T) {
	testMultiReturnValueFunctionSuccess(t, transpileBash)
}

func TestSingleParamFunctionSuccess(t *testing.T) {
	testSingleParamFunctionSuccess(t, transpileBash)
}

func TestMultiParamFunctionSuccess(t *testing.T) {
	testMultiParamFunctionSuccess(t, transpileBash)
}

func TestSliceParamFunctionSuccess(t *testing.T) {
	testSliceParamFunctionSuccess(t, transpileBash)
}

func TestCallFunctionFromFunctionSuccess(t *testing.T) {
	testCallFunctionFromFunctionSuccess(t, transpileBash)
}
