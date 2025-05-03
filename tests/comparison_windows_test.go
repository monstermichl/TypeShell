package tests

import (
	"testing"
)

func TestIntEqualSuccess(t *testing.T) {
	testIntEqualSuccess(t, transpileBatch)
}

func TestIntNotEqualSuccess(t *testing.T) {
	testIntNotEqualSuccess(t, transpileBatch)
}

func TestIntLessSuccess(t *testing.T) {
	testIntLessSuccess(t, transpileBatch)
}

func TestIntLessOrEqualSuccess(t *testing.T) {
	testIntLessOrEqualSuccess(t, transpileBatch)
}

func TestIntGreaterSuccess(t *testing.T) {
	testIntGreaterSuccess(t, transpileBatch)
}

func TestIntGreaterOrEqualSuccess(t *testing.T) {
	testIntGreaterOrEqualSuccess(t, transpileBatch)
}

func TestComplexIntComparisonSuccess(t *testing.T) {
	testComplexIntComparisonSuccess(t, transpileBatch)
}

func TestStringEqualSuccess(t *testing.T) {
	testStringEqualSuccess(t, transpileBatch)
}

func TestStringNotEqualSuccess(t *testing.T) {
	testStringNotEqualSuccess(t, transpileBatch)
}

func TestBooleanEqualSuccess(t *testing.T) {
	testBooleanEqualSuccess(t, transpileBatch)
}

func TestBooleanNotEqualSuccess(t *testing.T) {
	testBooleanNotEqualSuccess(t, transpileBatch)
}
