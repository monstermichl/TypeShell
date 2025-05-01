package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testPanicSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 1

		if a == 1 {
			panic("panic")
		}
		print("got through")
	`, func(output string, err error) {
		require.Error(t, err)
		require.Equal(t, "panic: panic", output)
	})
}
