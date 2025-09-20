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

func TestImportsFromExternalSourceSuccess(t *testing.T) {
	testImportsFromExternalSourceSuccess(t, transpileBashFunc)
}

func TestImportVariableSuccess(t *testing.T) {
	testImportVariableSuccess(t, transpileBashFunc)
}

func TestImportedSliceAssignmentSuccess(t *testing.T) {
	testImportedSliceAssignmentSuccess(t, transpileBashFunc)
}

func TestImportConstAssignmentFail(t *testing.T) {
	testImportConstAssignmentFail(t, transpileBashFunc)
}

func TestImportPrivateVariableFail(t *testing.T) {
	testImportPrivateVariableFail(t, transpileBashFunc)
}
