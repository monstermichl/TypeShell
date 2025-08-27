package tests

import "testing"

func TestStdStringsIndexSuccess(t *testing.T) {
	testStdStringsIndexSuccess(t, transpileBatchFunc)
}

func TestStdStringsContainsSuccess(t *testing.T) {
	testStdStringsContainsSuccess(t, transpileBatchFunc)
}

func TestStdStringsJoinSuccess(t *testing.T) {
	testStdStringsJoinSuccess(t, transpileBatchFunc)
}

func TestStdStringsHasPrefixSuccess(t *testing.T) {
	testStdStringsHasPrefixSuccess(t, transpileBatchFunc)
}

func TestStdStringsHasSuffixSuccess(t *testing.T) {
	testStdStringsHasSuffixSuccess(t, transpileBatchFunc)
}

func TestStdStringsCountSuccess(t *testing.T) {
	testStdStringsCountSuccess(t, transpileBatchFunc)
}

func TestStdStringsRepeatSuccess(t *testing.T) {
	testStdStringsRepeatSuccess(t, transpileBatchFunc)
}

func TestStdStringsReplaceSuccess(t *testing.T) {
	testStdStringsReplaceSuccess(t, transpileBatchFunc)
}

func TestStdStringsReplaceAllSuccess(t *testing.T) {
	testStdStringsReplaceAllSuccess(t, transpileBatchFunc)
}
