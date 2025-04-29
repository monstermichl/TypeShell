package tests

import (
	"testing"
)

func TestIfComparisonSuccess(t *testing.T) {
	testIfComparisonSuccess(t, transpileBash)
}

func TestNonBoolIfConditionFail(t *testing.T) {
	testNonBoolIfConditionFail(t, transpileBash)
}

func TestIfWithAndComparisonSuccess(t *testing.T) {
	testIfWithAndComparisonSuccess(t, transpileBash)
}

func TestIfWithOrComparisonSuccess(t *testing.T) {
	testIfWithOrComparisonSuccess(t, transpileBash)
}

func TestElseIfSuccess(t *testing.T) {
	testElseIfSuccess(t, transpileBash)
}

func TestElseSuccess(t *testing.T) {
	testElseSuccess(t, transpileBash)
}

func TestIfComparisonInFunctionSuccess(t *testing.T) {
	testIfComparisonInFunctionSuccess(t, transpileBash)
}

func TestIfWithAndComparisonInFunctionSuccess(t *testing.T) {
	testIfWithAndComparisonInFunctionSuccess(t, transpileBash)
}

func TestIfWithOrComparisonInFunctionSuccess(t *testing.T) {
	testIfWithOrComparisonInFunctionSuccess(t, transpileBash)
}

func TestElseIfInFunctionSuccess(t *testing.T) {
	testElseIfInFunctionSuccess(t, transpileBash)
}

func TestElseInFunctionSuccess(t *testing.T) {
	testElseInFunctionSuccess(t, transpileBash)
}
