package tests

import (
	"testing"
)

func TestPanicSuccess(t *testing.T) {
	testPanicSuccess(t, transpileBatch)
}
