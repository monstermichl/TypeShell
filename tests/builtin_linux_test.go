package tests

import (
	"testing"
)

func TestLenSliceSuccess(t *testing.T) {
	testLenSliceSuccess(t, transpileBash)
}

func TestLenStringSuccess(t *testing.T) {
	testLenStringSuccess(t, transpileBash)
}

func TestCopySuccess(t *testing.T) {
	testCopySuccess(t, transpileBash)
}

func TestReadSuccess(t *testing.T) {
	testReadSuccess(t, transpileBash)
}

func TestWriteSuccess(t *testing.T) {
	testWriteSuccess(t, transpileBash)
}

func TestPanicSuccess(t *testing.T) {
	testPanicSuccess(t, transpileBash)
}

func TestLenSliceInFunctionSuccess(t *testing.T) {
	testLenSliceInFunctionSuccess(t, transpileBash)
}

func TestLenStringInFunctionSuccess(t *testing.T) {
	testLenStringInFunctionSuccess(t, transpileBash)
}

func TestCopyInFunctionSuccess(t *testing.T) {
	testCopyInFunctionSuccess(t, transpileBash)
}

func TestReadInFunctionSuccess(t *testing.T) {
	testReadInFunctionSuccess(t, transpileBash)
}

func TestWriteInFunctionSuccess(t *testing.T) {
	testWriteInFunctionSuccess(t, transpileBash)
}

func TestPanicInFunctionSuccess(t *testing.T) {
	testPanicInFunctionSuccess(t, transpileBash)
}
