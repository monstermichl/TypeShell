package tests

import (
	"testing"
)

func TestLogicalAndSuccess(t *testing.T) {
	testLogicalAndSuccess(t, transpileBash)
}

func TestLogicalOrSuccess(t *testing.T) {
	testLogicalOrSuccess(t, transpileBash)
}

func TestComplexLogicalOperationSuccess(t *testing.T) {
	testComplexLogicalOperationSuccess(t, transpileBash)
}
