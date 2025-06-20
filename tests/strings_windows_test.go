package tests

import (
	"testing"
)

func TestStringConcatSuccess(t *testing.T) {
	testStringConcatSuccess(t, transpileBatch)
}

func TestStringLengthSuccess(t *testing.T) {
	testStringLengthSuccess(t, transpileBatch)
}

func TestStringSingleSubscriptSuccess(t *testing.T) {
	testStringSingleSubscriptSuccess(t, transpileBatch)
}

func TestStringStartSubscriptSuccess(t *testing.T) {
	testStringStartSubscriptSuccess(t, transpileBatch)
}

func TestStringEndSubscriptSuccess(t *testing.T) {
	testStringEndSubscriptSuccess(t, transpileBatch)
}

func TestStringRangeSubscriptSuccess(t *testing.T) {
	testStringRangeSubscriptSuccess(t, transpileBatch)
}

func TestStringRangeNoIndicesSubscriptSuccess(t *testing.T) {
	testStringRangeNoIndicesSubscriptSuccess(t, transpileBatch)
}

func TestStringWithNewlineSuccess(t *testing.T) {
	testStringWithNewlineSuccess(t, transpileBatch)
}

func TestStringWithoutNewlineSuccess(t *testing.T) {
	testStringWithoutNewlineSuccess(t, transpileBatch)
}

func TestMultilineStringSuccess(t *testing.T) {
	testMultilineStringSuccess(t, transpileBatch)
}

func TestItoaSuccess(t *testing.T) {
	testItoaSuccess(t, transpileBatch)
}
