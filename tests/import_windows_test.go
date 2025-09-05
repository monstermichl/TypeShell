package tests

import (
	"testing"
)

func TestSingleImportSuccess(t *testing.T) {
	testSingleImportSuccess(t, transpileBatchFunc)
}

func TestSeveralSingleImportsSuccess(t *testing.T) {
	testSeveralSingleImportsSuccess(t, transpileBatchFunc)
}

func TestMultiImportSuccess(t *testing.T) {
	testMultiImportSuccess(t, transpileBatchFunc)
}

func TestWildlyMixedImportsSuccess(t *testing.T) {
	testWildlyMixedImportsSuccess(t, transpileBatchFunc)
}

func TestImportsFromExternalSourceSuccess(t *testing.T) {
	testImportsFromExternalSourceSuccess(t, transpileBatchFunc)
}
