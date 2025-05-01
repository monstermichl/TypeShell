package tests

import (
	"testing"
)

func TestLenSliceSuccess(t *testing.T) {
	testLenSliceSuccess(t, transpileBatch)
}

func TestLenStringSuccess(t *testing.T) {
	testLenStringSuccess(t, transpileBatch)
}

func TestCopySuccess(t *testing.T) {
	testCopySuccess(t, transpileBatch)
}

func TestReadSuccess(t *testing.T) {
	testReadSuccess(t, transpileBatch)
}

func TestWriteSuccess(t *testing.T) {
	testWriteSuccess(t, transpileBatch)
}

func TestPanicSuccess(t *testing.T) {
	testPanicSuccess(t, transpileBatch)
}

func TestLenSliceInFunctionSuccess(t *testing.T) {
	testLenSliceInFunctionSuccess(t, transpileBatch)
}

func TestLenStringInFunctionSuccess(t *testing.T) {
	testLenStringInFunctionSuccess(t, transpileBatch)
}

func TestCopyInFunctionSuccess(t *testing.T) {
	testCopyInFunctionSuccess(t, transpileBatch)
}

func TestReadInFunctionSuccess(t *testing.T) {
	testReadInFunctionSuccess(t, transpileBatch)
}

func TestWriteInFunctionSuccess(t *testing.T) {
	testWriteInFunctionSuccess(t, transpileBatch)
}

func TestPanicInFunctionSuccess(t *testing.T) {
	testPanicInFunctionSuccess(t, transpileBatch)
}
