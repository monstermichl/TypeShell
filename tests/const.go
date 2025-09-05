package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testDefineConstantsSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		const a = 0
		const b int = 1
		const c, d = 2, 4 - 1
		const e, f int = 4, c + d

		print(a, b, c, d, e, f)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0 1 2 3 4 5", output)
	})
}

func testDefineConstantsInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func main() {
			const a = 0
			const b int = 1
			const c, d = 2, 4 - 1
			const e, f int = 4, c + d

			print(a, b, c, d, e, f)
		}
		main()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0 1 2 3 4 5", output)
	})
}

func testDefineConstantsGroupedSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		const (
			a = 0
			b, c = 1, 2
			d = iota
			e = 4
			f
			g, h = iota, iota
		)

		print(a, b, c, d, e, f, g, h)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0 1 2 3 4 5 6 7", output)
	})
}

func testDefineConstantsMissingValueFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		const (
			a = 0
			b, c = 1, 2
			d
			e = 4
			f, g = iota, iota
		)
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected data type or value assignment")
	})
}

func testDefineSameConstantFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		const a = 1
		const a = 2
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "constant a has already been defined")
	})
}

func testAssignFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		const a = 1
		a = 2
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "cannot assign a value to constant a")
	})
}

func testAssignFromFunctionFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func one() int {
			return 1
		}
		const a = one()
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected constant value")
	})
}
