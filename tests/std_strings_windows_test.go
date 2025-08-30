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

func TestStdStringsCutPrefixSuccess(t *testing.T) {
	testStdStringsCutPrefixSuccess(t, transpileBatchFunc)
}

func TestStdStringsCutSuffixSuccess(t *testing.T) {
	testStdStringsCutSuffixSuccess(t, transpileBatchFunc)
}

func TestStdStringsCutSuccess(t *testing.T) {
	testStdStringsCutSuccess(t, transpileBatchFunc)
}

func TestStdStringsTrimPrefixSuccess(t *testing.T) {
	testStdStringsTrimPrefixSuccess(t, transpileBatchFunc)
}

func TestStdStringsTrimSuffixSuccess(t *testing.T) {
	testStdStringsTrimSuffixSuccess(t, transpileBatchFunc)
}

func TestStdStringsTrimLeftSuccess(t *testing.T) {
	testStdStringsTrimLeftSuccess(t, transpileBatchFunc)
}

func TestStdStringsTrimRightSuccess(t *testing.T) {
	testStdStringsTrimRightSuccess(t, transpileBatchFunc)
}

func TestStdStringsTrimSuccess(t *testing.T) {
	testStdStringsTrimSuccess(t, transpileBatchFunc)
}

func TestStdStringsTrimSpaceSuccess(t *testing.T) {
	testStdStringsTrimSpaceSuccess(t, transpileBatchFunc)
}
