package tests

import (
	"testing"
)

func TestIntEqualSuccess(t *testing.T) {
	testIntEqualSuccess(t, transpileBash)
}

func TestIntNotEqualSuccess(t *testing.T) {
	testIntNotEqualSuccess(t, transpileBash)
}

func TestIntLessSuccess(t *testing.T) {
	testIntLessSuccess(t, transpileBash)
}

func TestIntLessOrEqualSuccess(t *testing.T) {
	testIntLessOrEqualSuccess(t, transpileBash)
}

func TestIntGreaterSuccess(t *testing.T) {
	testIntGreaterSuccess(t, transpileBash)
}

func TestIntGreaterOrEqualSuccess(t *testing.T) {
	testIntGreaterOrEqualSuccess(t, transpileBash)
}

func TestComplexIntComparisonSuccess(t *testing.T) {
	testComplexIntComparisonSuccess(t, transpileBash)
}

func TestStringEqualSuccess(t *testing.T) {
	testStringEqualSuccess(t, transpileBash)
}

func TestStringNotEqualSuccess(t *testing.T) {
	testStringNotEqualSuccess(t, transpileBash)
}

func TestBooleanEqualSuccess(t *testing.T) {
	testBooleanEqualSuccess(t, transpileBash)
}

func TestBooleanNotEqualSuccess(t *testing.T) {
	testBooleanNotEqualSuccess(t, transpileBash)
}
