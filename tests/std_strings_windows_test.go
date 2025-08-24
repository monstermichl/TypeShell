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
