package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testVarDefinition(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func one() int {
			return 7
		}
		func two() (int, int) {
			return 8, 9
		}
		var a int
		var b int = 1
		c := 2
		var d, e int = 3, 4
		f, g := 5, 6
		var h = one()
		var i, j = two()
		k := one()
		l, m := two()

		print(a, b, c, d, e, f, g, h, i, j, k, l, m)
	`, func(output string) {
		require.Equal(t, "0 1 2 3 4 5 6 7 8 9 7 8 9", output)
	})
}
