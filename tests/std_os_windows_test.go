package tests

import "testing"

func TestStdOsShellSuccess(t *testing.T) {
	testStdOsShellSuccess(t, transpileBatchFunc, "batch")
}
