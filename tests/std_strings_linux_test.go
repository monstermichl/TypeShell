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

func TestStdStringsCutPrefixSuccess(t *testing.T) {
	testStdStringsCutPrefixSuccess(t, transpileBashFunc)
}

func TestStdStringsCutSuffixSuccess(t *testing.T) {
	testStdStringsCutSuffixSuccess(t, transpileBashFunc)
}

func TestStdStringsCutSuccess(t *testing.T) {
	testStdStringsCutSuccess(t, transpileBashFunc)
}

func TestStdStringsTrimPrefixSuccess(t *testing.T) {
	testStdStringsTrimPrefixSuccess(t, transpileBashFunc)
}

func TestStdStringsTrimSuffixSuccess(t *testing.T) {
	testStdStringsTrimSuffixSuccess(t, transpileBashFunc)
}

func TestStdStringsTrimLeftSuccess(t *testing.T) {
	testStdStringsTrimLeftSuccess(t, transpileBashFunc)
}

func TestStdStringsTrimRightSuccess(t *testing.T) {
	testStdStringsTrimRightSuccess(t, transpileBashFunc)
}

func TestStdStringsTrimSuccess(t *testing.T) {
	testStdStringsTrimSuccess(t, transpileBashFunc)
}

func TestStdStringsTrimSpaceSuccess(t *testing.T) {
	testStdStringsTrimSpaceSuccess(t, transpileBashFunc)
}
