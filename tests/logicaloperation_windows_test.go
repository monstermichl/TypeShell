package tests

import (
	"testing"
)

func TestLogicalAndSuccess(t *testing.T) {
	testLogicalAndSuccess(t, transpileBatch)
}

func TestLogicalOrSuccess(t *testing.T) {
	testLogicalOrSuccess(t, transpileBatch)
}

func TestComplexLogicalOperationSuccess(t *testing.T) {
	testComplexLogicalOperationSuccess(t, transpileBatch)
}
