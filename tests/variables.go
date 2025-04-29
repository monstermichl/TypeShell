package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testDefineVariablesSuccess(t *testing.T, transpilerFunc transpilerFunc) {
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
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0 1 2 3 4 5 6 7 8 9 7 8 9", output)
	})
}

func testDefineSameVariableFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a int
		var a int
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "variable a has already been defined")
	})
}

func testNoNewVariableFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a, b int
		a, b := 1, 2
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "no new variables")
	})
}

func testAssignSuccessful(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a, b int
		a, b = 1, 2

		print(a, b)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2", output)
	})
}

func testAssignFromFunctionSuccessful(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func two() (int, int) {
			return 1, 2
		}
		var a, b int
		a, b = two()

		print(a, b)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2", output)
	})
}

func testAssignToUndefinedFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		a = 1
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "variable a has not been defined")
	})
}

func testDefineVariablesInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func one() int {
			return 7
		}
		func two() (int, int) {
			return 8, 9
		}
		func test() {
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
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0 1 2 3 4 5 6 7 8 9 7 8 9", output)
	})
}

func testDefineSameVariableInFunctionFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a int
			var a int
		}
		test()
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "variable a has already been defined")
	})
}

func testNoNewVariableInFunctionFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a, b int
			a, b := 1, 2
		}
		test()
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "no new variables")
	})
}

func testAssignInFunctionSuccessful(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a, b int
			a, b = 1, 2

			print(a, b)
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2", output)
	})
}

func testAssignFromFunctionInFunctionSuccessful(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func two() (int, int) {
			return 1, 2
		}
		func test() {
			var a, b int
			a, b = two()

			print(a, b)
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2", output)
	})
}

func testAssignToUndefinedInFunctionFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			a = 1
		}
		test()
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "variable a has not been defined")
	})
}
