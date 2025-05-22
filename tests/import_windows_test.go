package tests

import (
	"testing"
)

func TestSingleImportSuccess(t *testing.T) {
	testSingleImportSuccess(t, transpileBatchFunc)
}

func TestMultiImportSuccess(t *testing.T) {
	testMultiImportSuccess(t, transpileBatchFunc)
}
