package tests

import "testing"

func TestStdShellGetSuccess(t *testing.T) {
	testStdShellGetSuccess(t, transpileBatchFunc, "batch")
}
