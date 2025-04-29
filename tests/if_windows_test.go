package tests

import (
	"testing"
)

func TestIfComparisonSuccess(t *testing.T) {
	testIfComparisonSuccess(t, transpileBatch)
}

func TestNonBoolIfConditionFail(t *testing.T) {
	testNonBoolIfConditionFail(t, transpileBatch)
}

func TestIfWithAndComparisonSuccess(t *testing.T) {
	testIfWithAndComparisonSuccess(t, transpileBatch)
}

func TestIfWithOrComparisonSuccess(t *testing.T) {
	testIfWithOrComparisonSuccess(t, transpileBatch)
}

func TestElseIfSuccess(t *testing.T) {
	testElseIfSuccess(t, transpileBatch)
}

func TestElseSuccess(t *testing.T) {
	testElseSuccess(t, transpileBatch)
}

func TestIfComparisonInFunctionSuccess(t *testing.T) {
	testIfComparisonInFunctionSuccess(t, transpileBatch)
}

func TestIfWithAndComparisonInFunctionSuccess(t *testing.T) {
	testIfWithAndComparisonInFunctionSuccess(t, transpileBatch)
}

func TestIfWithOrComparisonInFunctionSuccess(t *testing.T) {
	testIfWithOrComparisonInFunctionSuccess(t, transpileBatch)
}

func TestElseIfInFunctionSuccess(t *testing.T) {
	testElseIfInFunctionSuccess(t, transpileBatch)
}

func TestElseInFunctionSuccess(t *testing.T) {
	testElseInFunctionSuccess(t, transpileBatch)
}
