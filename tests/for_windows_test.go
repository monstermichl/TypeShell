package tests

import (
	"testing"
)

func TestForComparisonSuccess(t *testing.T) {
	testForComparisonSuccess(t, transpileBatch)
}

func TestNonBoolForConditionFail(t *testing.T) {
	testNonBoolForConditionFail(t, transpileBatch)
}

func TestForWithAndComparisonSuccess(t *testing.T) {
	testForWithAndComparisonSuccess(t, transpileBatch)
}

func TestForWithOrComparisonSuccess(t *testing.T) {
	testForWithOrComparisonSuccess(t, transpileBatch)
}

func TestForWithCountingVariableSuccess(t *testing.T) {
	testForWithCountingVariableSuccess(t, transpileBatch)
}

func TestForWithSeparateCountingVariableSuccess(t *testing.T) {
	testForWithSeparateCountingVariableSuccess(t, transpileBatch)
}

func TestForWithSeparateCountingVariableAndSeparateIncrementSuccess(t *testing.T) {
	testForWithSeparateCountingVariableAndSeparateIncrementSuccess(t, transpileBatch)
}

func TestForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementSuccess(t *testing.T) {
	testForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementSuccess(t, transpileBatch)
}

func TestForWithNoConditionSuccess(t *testing.T) {
	testForWithNoConditionSuccess(t, transpileBatch)
}

func TestForRangeSliceSuccess(t *testing.T) {
	testForRangeSliceSuccess(t, transpileBatch)
}

func TestForRangeStringSuccess(t *testing.T) {
	testForRangeStringSuccess(t, transpileBatch)
}

func TestForRangeNonIterableFail(t *testing.T) {
	testForRangeNonIterableFail(t, transpileBatch)
}

func TestForComparisonInFunctionSuccess(t *testing.T) {
	testForComparisonInFunctionSuccess(t, transpileBatch)
}

func TestNonBoolForConditionInFunctionFail(t *testing.T) {
	testNonBoolForConditionInFunctionFail(t, transpileBatch)
}

func TestForWithAndComparisonInFunctionSuccess(t *testing.T) {
	testForWithAndComparisonInFunctionSuccess(t, transpileBatch)
}

func TestForWithOrComparisonInFunctionSuccess(t *testing.T) {
	testForWithOrComparisonInFunctionSuccess(t, transpileBatch)
}

func TestForWithCountingVariableInFunctionSuccess(t *testing.T) {
	testForWithCountingVariableInFunctionSuccess(t, transpileBatch)
}

func TestForWithSeparateCountingVariableInFunctionSuccess(t *testing.T) {
	testForWithSeparateCountingVariableInFunctionSuccess(t, transpileBatch)
}

func TestForWithSeparateCountingVariableAndSeparateIncrementInFunctionSuccess(t *testing.T) {
	testForWithSeparateCountingVariableAndSeparateIncrementInFunctionSuccess(t, transpileBatch)
}

func TestForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementInFunctionSuccess(t *testing.T) {
	testForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementInFunctionSuccess(t, transpileBatch)
}

func TestForWithNoConditionInFunctionSuccess(t *testing.T) {
	testForWithNoConditionInFunctionSuccess(t, transpileBatch)
}

func TestForRangeSliceInFunctionSuccess(t *testing.T) {
	testForRangeSliceInFunctionSuccess(t, transpileBatch)
}

func TestForRangeStringInFunctionSuccess(t *testing.T) {
	testForRangeStringInFunctionSuccess(t, transpileBatch)
}
