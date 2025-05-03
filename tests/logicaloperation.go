package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testLogicalAndSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = true
		var b = false
		var c = true && false

		print(c)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0", output)
	})
}

func testLogicalOrSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = true
		var b = false
		var c = true || false

		print(c)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testComplexLogicalOperationSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = ((true && false) || (!false && true)) && (!true || false || (true && !false))

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}
