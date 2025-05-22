package tests

import (
	"testing"
)

func TestSingleImportSuccess(t *testing.T) {
	testSingleImportSuccess(t, transpileBashFunc)
}

func TestMultiImportSuccess(t *testing.T) {
	testMultiImportSuccess(t, transpileBashFunc)
}
