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

func TestStdStringsHasPrefixSuccess(t *testing.T) {
	testStdStringsHasPrefixSuccess(t, transpileBashFunc)
}

func TestStdStringsHasSuffixSuccess(t *testing.T) {
	testStdStringsHasSuffixSuccess(t, transpileBashFunc)
}

func TestStdStringsCountSuccess(t *testing.T) {
	testStdStringsCountSuccess(t, transpileBashFunc)
}

func TestStdStringsRepeatSuccess(t *testing.T) {
	testStdStringsRepeatSuccess(t, transpileBashFunc)
}

func TestStdStringsReplaceSuccess(t *testing.T) {
	testStdStringsReplaceSuccess(t, transpileBashFunc)
}

func TestStdStringsReplaceAllSuccess(t *testing.T) {
	testStdStringsReplaceAllSuccess(t, transpileBashFunc)
}
