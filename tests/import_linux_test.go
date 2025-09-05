package tests

import (
	"testing"
)

func TestSingleImportSuccess(t *testing.T) {
	testSingleImportSuccess(t, transpileBashFunc)
}

func TestSeveralSingleImportsSuccess(t *testing.T) {
	testSeveralSingleImportsSuccess(t, transpileBashFunc)
}

func TestMultiImportSuccess(t *testing.T) {
	testMultiImportSuccess(t, transpileBashFunc)
}

func TestWildlyMixedImportsSuccess(t *testing.T) {
	testWildlyMixedImportsSuccess(t, transpileBashFunc)
}
