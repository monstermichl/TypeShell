package tests

import (
	"testing"
)

func TestDeclareAndDefineStructSuccess(t *testing.T) {
	testDeclareAndDefineStructSuccess(t, transpileBatch)
}

func TestDeclareAndDefineStructWithValuesSuccess(t *testing.T) {
	testDeclareAndDefineStructWithValuesSuccess(t, transpileBatch)
}

func TestDeclareAndDefineStructWithValuesOneLineSuccess(t *testing.T) {
	testDeclareAndDefineStructWithValuesOneLineSuccess(t, transpileBatch)
}

func TestDeclareAndDefineStructSliceSuccess(t *testing.T) {
	testDeclareAndDefineStructSliceSuccess(t, transpileBatch)
}

func TestPassStructToFunctionSuccess(t *testing.T) {
	testPassStructToFunctionSuccess(t, transpileBatch)
}

func TestReturnDifferentStructsSuccess(t *testing.T) {
	testReturnDifferentStructsSuccess(t, transpileBatch)
}

func TestNestedStructSuccess(t *testing.T) {
	testNestedStructSuccess(t, transpileBatch)
}

func TestStructEvaluationChainingSuccess(t *testing.T) {
	testStructEvaluationChainingSuccess(t, transpileBatch)
}

func TestStructAssignmentChainingSuccess(t *testing.T) {
	testStructAssignmentChainingSuccess(t, transpileBatch)
}

func TestStructFieldAssignedTwiceInInitializationFail(t *testing.T) {
	testStructFieldAssignedTwiceInInitializationFail(t, transpileBatch)
}

func TestStructUnknownFieldInInitializationFail(t *testing.T) {
	testStructUnknownFieldInInitializationFail(t, transpileBatch)
}

func TestStructFieldWrongTypeAssignmentInInitializationFail(t *testing.T) {
	testStructFieldWrongTypeAssignmentInInitializationFail(t, transpileBatch)
}

func TestStructUnknownFieldInAssignmentFail(t *testing.T) {
	testStructUnknownFieldInAssignmentFail(t, transpileBatch)
}

func TestStructFieldWrongTypeAssignmentInAssignmentFail(t *testing.T) {
	testStructFieldWrongTypeAssignmentInAssignmentFail(t, transpileBatch)
}
