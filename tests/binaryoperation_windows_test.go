package tests

import (
	"testing"
)

func TestAdditionSuccess(t *testing.T) {
	testAdditionSuccess(t, transpileBatch)
}

func TestSubtractionSuccess(t *testing.T) {
	testSubtractionSuccess(t, transpileBatch)
}

func TestMultiplicationSuccess(t *testing.T) {
	testMultiplicationSuccess(t, transpileBatch)
}

func TestDivisionSuccess(t *testing.T) {
	testDivisionSuccess(t, transpileBatch)
}

func TestModuloSuccess(t *testing.T) {
	testModuloSuccess(t, transpileBatch)
}

func TestMoreComplexCalculationSuccess(t *testing.T) {
	testMoreComplexCalculationSuccess(t, transpileBatch)
}

func TestMoreComplexCalculationWithBracketsSuccess(t *testing.T) {
	testMoreComplexCalculationWithBracketsSuccess(t, transpileBatch)
}

func TestComplexCalculationSuccess(t *testing.T) {
	testComplexCalculationSuccess(t, transpileBatch)
}

func TestCompoundAssignmentAdditionSuccess(t *testing.T) {
	testCompoundAssignmentAdditionSuccess(t, transpileBatch)
}

func TestCompoundAssignmentSubtractionSuccess(t *testing.T) {
	testCompoundAssignmentSubtractionSuccess(t, transpileBatch)
}

func TestCompoundAssignmentMultiplicationSuccess(t *testing.T) {
	testCompoundAssignmentMultiplicationSuccess(t, transpileBatch)
}

func TestCompoundAssignmentDivisionSuccess(t *testing.T) {
	testCompoundAssignmentDivisionSuccess(t, transpileBatch)
}

func TestCompoundAssignmentModuloSuccess(t *testing.T) {
	testCompoundAssignmentModuloSuccess(t, transpileBatch)
}
