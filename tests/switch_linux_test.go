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
