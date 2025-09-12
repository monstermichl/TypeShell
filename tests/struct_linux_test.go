package tests

import (
	"testing"
)

func TestDeclareAndDefineStructSuccess(t *testing.T) {
	testDeclareAndDefineStructSuccess(t, transpileBash)
}

func TestDeclareAndDefineStructWithValuesSuccess(t *testing.T) {
	testDeclareAndDefineStructWithValuesSuccess(t, transpileBash)
}

func TestDeclareAndDefineStructWithValuesOneLineSuccess(t *testing.T) {
	testDeclareAndDefineStructWithValuesOneLineSuccess(t, transpileBash)
}

func TestDeclareAndDefineStructSliceSuccess(t *testing.T) {
	testDeclareAndDefineStructSliceSuccess(t, transpileBash)
}

func TestStructFieldAssignedTwiceInInitializationFail(t *testing.T) {
	testStructFieldAssignedTwiceInInitializationFail(t, transpileBash)
}

func TestStructUnknownFieldInInitializationFail(t *testing.T) {
	testStructUnknownFieldInInitializationFail(t, transpileBash)
}

func TestStructFieldWrongTypeAssignmentInInitializationFail(t *testing.T) {
	testStructFieldWrongTypeAssignmentInInitializationFail(t, transpileBash)
}

func TestStructUnknownFieldInAssignmentFail(t *testing.T) {
	testStructUnknownFieldInAssignmentFail(t, transpileBash)
}

func TestStructFieldWrongTypeAssignmentInAssignmentFail(t *testing.T) {
	testStructFieldWrongTypeAssignmentInAssignmentFail(t, transpileBash)
}
