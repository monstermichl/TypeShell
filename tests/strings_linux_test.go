package tests

import (
	"testing"
)

func TestStringConcatSuccess(t *testing.T) {
	testStringConcatSuccess(t, transpileBash)
}

func TestStringLengthSuccess(t *testing.T) {
	testStringLengthSuccess(t, transpileBash)
}

func TestStringSingleSubscriptSuccess(t *testing.T) {
	testStringSingleSubscriptSuccess(t, transpileBash)
}

func TestStringStartSubscriptSuccess(t *testing.T) {
	testStringStartSubscriptSuccess(t, transpileBash)
}

func TestStringEndSubscriptSuccess(t *testing.T) {
	testStringEndSubscriptSuccess(t, transpileBash)
}

func TestStringRangeSubscriptSuccess(t *testing.T) {
	testStringRangeSubscriptSuccess(t, transpileBash)
}

func TestStringRangeNoIndicesSubscriptSuccess(t *testing.T) {
	testStringRangeNoIndicesSubscriptSuccess(t, transpileBash)
}

func TestStringWithNewlineSuccess(t *testing.T) {
	testStringWithNewlineSuccess(t, transpileBash)
}

func TestStringWithoutNewlineSuccess(t *testing.T) {
	testStringWithoutNewlineSuccess(t, transpileBash)
}

func TestMultilineStringSuccess(t *testing.T) {
	testMultilineStringSuccess(t, transpileBash)
}

func TestItoaSuccess(t *testing.T) {
	testItoaSuccess(t, transpileBash)
}
