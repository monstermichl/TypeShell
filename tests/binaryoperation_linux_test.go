package tests

import (
	"testing"
)

func TestAdditionSuccess(t *testing.T) {
	testAdditionSuccess(t, transpileBash)
}

func TestSubtractionSuccess(t *testing.T) {
	testSubtractionSuccess(t, transpileBash)
}

func TestMultiplicationSuccess(t *testing.T) {
	testMultiplicationSuccess(t, transpileBash)
}

func TestDivisionSuccess(t *testing.T) {
	testDivisionSuccess(t, transpileBash)
}

func TestModuloSuccess(t *testing.T) {
	testModuloSuccess(t, transpileBash)
}

func TestMoreComplexCalculationSuccess(t *testing.T) {
	testMoreComplexCalculationSuccess(t, transpileBash)
}

func TestMoreComplexCalculationWithBracketsSuccess(t *testing.T) {
	testMoreComplexCalculationWithBracketsSuccess(t, transpileBash)
}

func TestComplexCalculationSuccess(t *testing.T) {
	testComplexCalculationSuccess(t, transpileBash)
}

func TestCompoundAssignmentAdditionSuccess(t *testing.T) {
	testCompoundAssignmentAdditionSuccess(t, transpileBash)
}

func TestCompoundAssignmentSubtractionSuccess(t *testing.T) {
	testCompoundAssignmentSubtractionSuccess(t, transpileBash)
}

func TestCompoundAssignmentMultiplicationSuccess(t *testing.T) {
	testCompoundAssignmentMultiplicationSuccess(t, transpileBash)
}

func TestCompoundAssignmentDivisionSuccess(t *testing.T) {
	testCompoundAssignmentDivisionSuccess(t, transpileBash)
}

func TestCompoundAssignmentModuloSuccess(t *testing.T) {
	testCompoundAssignmentModuloSuccess(t, transpileBash)
}
