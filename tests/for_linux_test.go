package tests

import (
	"testing"
)

func TestForComparisonSuccess(t *testing.T) {
	testForComparisonSuccess(t, transpileBash)
}

func TestNonBoolForConditionFail(t *testing.T) {
	testNonBoolForConditionFail(t, transpileBash)
}

func TestForWithAndComparisonSuccess(t *testing.T) {
	testForWithAndComparisonSuccess(t, transpileBash)
}

func TestForWithOrComparisonSuccess(t *testing.T) {
	testForWithOrComparisonSuccess(t, transpileBash)
}

func TestForWithCountingVariableSuccess(t *testing.T) {
	testForWithCountingVariableSuccess(t, transpileBash)
}

func TestForWithSeparateCountingVariableSuccess(t *testing.T) {
	testForWithSeparateCountingVariableSuccess(t, transpileBash)
}

func TestForWithSeparateCountingVariableAndSeparateIncrementSuccess(t *testing.T) {
	testForWithSeparateCountingVariableAndSeparateIncrementSuccess(t, transpileBash)
}

func TestForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementSuccess(t *testing.T) {
	testForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementSuccess(t, transpileBash)
}

func TestForWithNoConditionSuccess(t *testing.T) {
	testForWithNoConditionSuccess(t, transpileBash)
}

func TestForRangeSliceSuccess(t *testing.T) {
	testForRangeSliceSuccess(t, transpileBash)
}

func TestForRangeStringSuccess(t *testing.T) {
	testForRangeStringSuccess(t, transpileBash)
}

func TestForRangeNonIterableFail(t *testing.T) {
	testForRangeNonIterableFail(t, transpileBash)
}

func TestForComparisonInFunctionSuccess(t *testing.T) {
	testForComparisonInFunctionSuccess(t, transpileBash)
}

func TestNonBoolForConditionInFunctionFail(t *testing.T) {
	testNonBoolForConditionInFunctionFail(t, transpileBash)
}

func TestForWithAndComparisonInFunctionSuccess(t *testing.T) {
	testForWithAndComparisonInFunctionSuccess(t, transpileBash)
}

func TestForWithOrComparisonInFunctionSuccess(t *testing.T) {
	testForWithOrComparisonInFunctionSuccess(t, transpileBash)
}

func TestForWithCountingVariableInFunctionSuccess(t *testing.T) {
	testForWithCountingVariableInFunctionSuccess(t, transpileBash)
}

func TestForWithSeparateCountingVariableInFunctionSuccess(t *testing.T) {
	testForWithSeparateCountingVariableInFunctionSuccess(t, transpileBash)
}

func TestForWithSeparateCountingVariableAndSeparateIncrementInFunctionSuccess(t *testing.T) {
	testForWithSeparateCountingVariableAndSeparateIncrementInFunctionSuccess(t, transpileBash)
}

func TestForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementInFunctionSuccess(t *testing.T) {
	testForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementInFunctionSuccess(t, transpileBash)
}

func TestForWithNoConditionInFunctionSuccess(t *testing.T) {
	testForWithNoConditionInFunctionSuccess(t, transpileBash)
}

func TestForRangeSliceInFunctionSuccess(t *testing.T) {
	testForRangeSliceInFunctionSuccess(t, transpileBash)
}

func TestForRangeStringInFunctionSuccess(t *testing.T) {
	testForRangeStringInFunctionSuccess(t, transpileBash)
}
