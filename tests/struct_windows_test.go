package tests

import (
	"testing"
)

func TestDeclareAndDefineStructSuccess(t *testing.T) {
	testDeclareAndDefineStructSuccess(t, transpileBatch)
}

func TestDeclareAndDefineStructWithValuesSuccess(t *testing.T) {
	testDeclareAndDefineStructWithValuesSuccess(t, transpileBatch)
}

func TestDeclareAndDefineStructSliceSuccess(t *testing.T) {
	testDeclareAndDefineStructSliceSuccess(t, transpileBatch)
}
