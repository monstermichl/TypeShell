package tests

import (
	"testing"
)

func TestSwitchWithBoolSuccess(t *testing.T) {
	testSwitchWithBoolSuccess(t, transpileBash)
}

func TestSwitchWithBoolDefaultSuccess(t *testing.T) {
	testSwitchWithBoolDefaultSuccess(t, transpileBash)
}

func TestSwitchWithImplicitBoolSuccess(t *testing.T) {
	testSwitchWithImplicitBoolSuccess(t, transpileBash)
}

func TestSwitchWithComparisonsSuccess(t *testing.T) {
	testSwitchWithComparisonsSuccess(t, transpileBash)
}

func TestSwitchOnlyDefaultSuccess(t *testing.T) {
	testSwitchOnlyDefaultSuccess(t, transpileBash)
}

func TestSwitchStringsSuccess(t *testing.T) {
	testSwitchStringsSuccess(t, transpileBash)
}

func TestSwitchWithBoolInFunctionSuccess(t *testing.T) {
	testSwitchWithBoolInFunctionSuccess(t, transpileBash)
}

func TestSwitchWithBoolDefaultInFunctionSuccess(t *testing.T) {
	testSwitchWithBoolDefaultInFunctionSuccess(t, transpileBash)
}

func TestSwitchWithImplicitBoolInFunctionSuccess(t *testing.T) {
	testSwitchWithImplicitBoolInFunctionSuccess(t, transpileBash)
}

func TestSwitchWithComparisonsInFunctionSuccess(t *testing.T) {
	testSwitchWithComparisonsInFunctionSuccess(t, transpileBash)
}

func TestSwitchOnlyDefaultInFunctionSuccess(t *testing.T) {
	testSwitchOnlyDefaultInFunctionSuccess(t, transpileBash)
}

func TestSwitchStringsInFunctionSuccess(t *testing.T) {
	testSwitchStringsInFunctionSuccess(t, transpileBash)
}
