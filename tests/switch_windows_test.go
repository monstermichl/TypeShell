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
