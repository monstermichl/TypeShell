package tests

import (
	"testing"
)

func TestTypeDeclarationAndDefinitionSuccess(t *testing.T) {
	testTypeDeclarationAndDefinitionSuccess(t, transpileBatch)
}

func TestTypeDeclarationAndAssignmentFail(t *testing.T) {
	testTypeDeclarationAndAssignmentFail(t, transpileBatch)
}

func TestTypeAliasAndAssignmentSuccess(t *testing.T) {
	testTypeAliasAndAssignmentSuccess(t, transpileBatch)
}

func TestTypeDeclarationAndDefinitionInFunctionSuccess(t *testing.T) {
	testTypeDeclarationAndDefinitionInFunctionSuccess(t, transpileBatch)
}

func TestTypeDeclarationAndAssignmentInFunctionFail(t *testing.T) {
	testTypeDeclarationAndAssignmentInFunctionFail(t, transpileBatch)
}

func TestTypeAliasAndAssignmentInFunctionSuccess(t *testing.T) {
	testTypeAliasAndAssignmentInFunctionSuccess(t, transpileBatch)
}

func TestTypeDeclaredInFunctionUsedOutsideFail(t *testing.T) {
	testTypeDeclaredInFunctionUsedOutsideFail(t, transpileBatch)
}

func TestTypeDeclaredInIfUsedOutsideFail(t *testing.T) {
	testTypeDeclaredInIfUsedOutsideFail(t, transpileBatch)
}

func TestDeclareTypeTwiceFail(t *testing.T) {
	testDeclareTypeTwiceFail(t, transpileBatch)
}

func TestPassDeclaredTypeToFunctionSuccess(t *testing.T) {
	testPassDeclaredTypeToFunctionSuccess(t, transpileBatch)
}

func TestPassBaseTypeToFunctionFail(t *testing.T) {
	testPassBaseTypeToFunctionFail(t, transpileBatch)
}

func TestPassValueWithSameBaseTypeToFunctionSuccess(t *testing.T) {
	testPassValueWithSameBaseTypeToFunctionSuccess(t, transpileBatch)
}

func TtestAssignDifferentDefinedTypeFail(t *testing.T) {
	testAssignDifferentDefinedTypeFail(t, transpileBatch)
}
