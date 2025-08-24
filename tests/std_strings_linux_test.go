package tests

import "testing"

func TestStdStringsIndexSuccess(t *testing.T) {
	testStdStringsIndexSuccess(t, transpileBashFunc)
}

func TestStdStringsContainsSuccess(t *testing.T) {
	testStdStringsContainsSuccess(t, transpileBashFunc)
}

func TestStdStringsJoinSuccess(t *testing.T) {
	testStdStringsJoinSuccess(t, transpileBashFunc)
}
