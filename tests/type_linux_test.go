package tests

import (
	"testing"
)

func TestTypeDeclarationAndDefinitionSuccess(t *testing.T) {
	testTypeDeclarationAndDefinitionSuccess(t, transpileBash)
}

func TestTypeDeclarationAndAssignmentFail(t *testing.T) {
	testTypeDeclarationAndAssignmentFail(t, transpileBash)
}

func TestTypeAliasAndAssignmentSuccess(t *testing.T) {
	testTypeAliasAndAssignmentSuccess(t, transpileBash)
}

func TestTypeDeclarationAndDefinitionInFunctionSuccess(t *testing.T) {
	testTypeDeclarationAndDefinitionInFunctionSuccess(t, transpileBash)
}

func TestTypeDeclarationAndAssignmentInFunctionFail(t *testing.T) {
	testTypeDeclarationAndAssignmentInFunctionFail(t, transpileBash)
}

func TestTypeAliasAndAssignmentInFunctionSuccess(t *testing.T) {
	testTypeAliasAndAssignmentInFunctionSuccess(t, transpileBash)
}

func TestTypeDeclaredInFunctionUsedOutsideFail(t *testing.T) {
	testTypeDeclaredInFunctionUsedOutsideFail(t, transpileBash)
}

func TestTypeDeclaredInIfUsedOutsideFail(t *testing.T) {
	testTypeDeclaredInIfUsedOutsideFail(t, transpileBash)
}

func TestDeclareTypeTwiceFail(t *testing.T) {
	testDeclareTypeTwiceFail(t, transpileBash)
}

func TestPassDeclaredTypeToFunctionSuccess(t *testing.T) {
	testPassDeclaredTypeToFunctionSuccess(t, transpileBash)
}

func TestPassBaseTypeToFunctionFail(t *testing.T) {
	testPassBaseTypeToFunctionFail(t, transpileBash)
}

func TestPassValueWithSameBaseTypeToFunctionSuccess(t *testing.T) {
	testPassValueWithSameBaseTypeToFunctionSuccess(t, transpileBash)
}

func TtestAssignDifferentDefinedTypeFail(t *testing.T) {
	testAssignDifferentDefinedTypeFail(t, transpileBash)
}
