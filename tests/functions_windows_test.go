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

func TestMultiParamWithSameTypeFunctionSuccess(t *testing.T) {
	testMultiParamWithSameTypeFunctionSuccess(t, transpileBatch)
}

func TestSliceParamFunctionSuccess(t *testing.T) {
	testSliceParamFunctionSuccess(t, transpileBatch)
}

func TestStructParamFunctionSuccess(t *testing.T) {
	testStructParamFunctionSuccess(t, transpileBatch)
}

func TestStructPointerParamFunctionSuccess(t *testing.T) {
	testStructPointerParamFunctionSuccess(t, transpileBatch)
}

func TestStructReceiverFunctionSuccess(t *testing.T) {
	testStructReceiverFunctionSuccess(t, transpileBatch)
}

func TestStructPointerReceiverFunctionSuccess(t *testing.T) {
	testStructPointerReceiverFunctionSuccess(t, transpileBatch)
}

func TestOptionalParamsFunctionWithParamsSuccess(t *testing.T) {
	testOptionalParamsFunctionWithParamsSuccess(t, transpileBatch)
}

func TestOptionalParamsFunctionWithoutParamsSuccess(t *testing.T) {
	testOptionalParamsFunctionWithoutParamsSuccess(t, transpileBatch)
}

func TestOptionalParamsFunctionWithWrongTypeFail(t *testing.T) {
	testOptionalParamsFunctionWithWrongTypeFail(t, transpileBatch)
}

func TestCallFunctionFromFunctionSuccess(t *testing.T) {
	testCallFunctionFromFunctionSuccess(t, transpileBatch)
}
