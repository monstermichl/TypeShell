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

func TestCallImportedFunctionAsStatementSuccess(t *testing.T) {
	testCallImportedFunctionAsStatementSuccess(t, transpileBatchFunc)
}

func TestImportVariableSuccess(t *testing.T) {
	testImportVariableSuccess(t, transpileBatchFunc)
}

func TestImportedSliceAssignmentSuccess(t *testing.T) {
	testImportedSliceAssignmentSuccess(t, transpileBatchFunc)
}

func TestImportSimpleTypeSuccess(t *testing.T) {
	testImportSimpleTypeSuccess(t, transpileBatchFunc)
}

func TestImportStructTypeSuccess(t *testing.T) {
	testImportStructTypeSuccess(t, transpileBatchFunc)
}

func TestImportStructTypePrivateFieldInitFail(t *testing.T) {
	testImportStructTypePrivateFieldInitFail(t, transpileBatchFunc)
}

func TestImportStructTypePrivateFieldEvaluationFail(t *testing.T) {
	testImportStructTypePrivateFieldEvaluationFail(t, transpileBatchFunc)
}

func TestImportConstAssignmentFail(t *testing.T) {
	testImportConstAssignmentFail(t, transpileBatchFunc)
}

func TestImportPrivateVariableFail(t *testing.T) {
	testImportPrivateVariableFail(t, transpileBatchFunc)
}
