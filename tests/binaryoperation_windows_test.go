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
