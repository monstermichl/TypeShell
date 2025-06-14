package tests

import (
	"testing"
)

func TestSwitchWithBoolSuccess(t *testing.T) {
	testSwitchWithBoolSuccess(t, transpileBatch)
}

func TestSwitchWithBoolDefaultSuccess(t *testing.T) {
	testSwitchWithBoolDefaultSuccess(t, transpileBatch)
}

func TestSwitchWithImplicitBoolSuccess(t *testing.T) {
	testSwitchWithImplicitBoolSuccess(t, transpileBatch)
}

func TestSwitchWithComparisonsSuccess(t *testing.T) {
	testSwitchWithComparisonsSuccess(t, transpileBatch)
}

func TestSwitchOnlyDefaultSuccess(t *testing.T) {
	testSwitchOnlyDefaultSuccess(t, transpileBatch)
}

func TestSwitchStringsSuccess(t *testing.T) {
	testSwitchStringsSuccess(t, transpileBatch)
}

func TestSwitchWithBoolInFunctionSuccess(t *testing.T) {
	testSwitchWithBoolInFunctionSuccess(t, transpileBatch)
}

func TestSwitchWithBoolDefaultInFunctionSuccess(t *testing.T) {
	testSwitchWithBoolDefaultInFunctionSuccess(t, transpileBatch)
}

func TestSwitchWithImplicitBoolInFunctionSuccess(t *testing.T) {
	testSwitchWithImplicitBoolInFunctionSuccess(t, transpileBatch)
}

func TestSwitchWithComparisonsInFunctionSuccess(t *testing.T) {
	testSwitchWithComparisonsInFunctionSuccess(t, transpileBatch)
}

func TestSwitchOnlyDefaultInFunctionSuccess(t *testing.T) {
	testSwitchOnlyDefaultInFunctionSuccess(t, transpileBatch)
}

func TestSwitchStringsInFunctionSuccess(t *testing.T) {
	testSwitchStringsInFunctionSuccess(t, transpileBatch)
}
