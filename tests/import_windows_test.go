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

func TestImportVariableSuccess(t *testing.T) {
	testImportVariableSuccess(t, transpileBatchFunc)
}

func TestImportedSliceAssignmentSuccess(t *testing.T) {
	testImportedSliceAssignmentSuccess(t, transpileBatchFunc)
}

func TestImportConstAssignmentFail(t *testing.T) {
	testImportConstAssignmentFail(t, transpileBatchFunc)
}

func TestImportPrivateVariableFail(t *testing.T) {
	testImportPrivateVariableFail(t, transpileBatchFunc)
}
