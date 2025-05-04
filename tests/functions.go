package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testVoidFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			print("Hello World")
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello World", output)
	})
}

func testSingleReturnValueFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() int {
			return 1
		}
		print(test())
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testMultiReturnValueFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() (int, int) {
			return 1, 2
		}
		a, b := test()
		print(a, b)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2", output)
	})
}

func testSingleParamFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test(retVal int) int {
			return retVal
		}
		print(test(1))
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testMultiParamFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test(retVal1 int, retVal2 int) (int, int) {
			return retVal1, retVal2
		}
		a, b := test(1, 2)
		print(a, b)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2", output)
	})
}

func testSliceParamFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test(s []int) {
			s[0] = 1
			s[1] = 2
		}
		s := []int{}
		print(len(s))

		test(s)
		print(len(s))
		print(s[0], s[1])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0\n2\n1 2", output)
	})
}

func testCallFunctionFromFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test1(retVal1 int, retVal2 int) (int, int) {
			return retVal1, retVal2
		}
		func test2() (int, int) {
			a, b := test1(1, 2)
			return a + 1, b + 1
		}
		a, b := test2()
		print(a, b)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "2 3", output)
	})
}
