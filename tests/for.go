package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testForComparisonSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 0

		for a < 2 {
			print("ok")
			a++
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testNonBoolForConditionFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		for "1" {
			print("ok")
		}
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected boolean expression")
	})
}

func testForWithAndComparisonSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 0
		var b = 1

		for a < 2 && b == 1 {
			print("ok")
			a++
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithOrComparisonSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 0
		var b = 1

		for a < 2 || b == 1 {
			print("ok")
			a++

			if a == 2 {
				b = 2
			}
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithCountingVariableSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		for i := 0; i < 2; i++ {
			print("ok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithSeparateCountingVariableSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		i := 0
		for ; i < 2; i++ {
			print("ok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithSeparateCountingVariableAndSeparateIncrementSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		i := 0
		for ; i < 2; {
			print("ok")
			i++
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		i := 0
		for ; ; {
			if i >= 2 {
				break
			}
			print("ok")
			i++
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithNoConditionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			i := 0

			for {
				if i >= 2 {
					break
				}
				print("ok")
				i++
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}




func testForComparisonInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = 0

			for a < 2 {
				print("ok")
				a++
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testNonBoolForConditionInFunctionFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			for "1" {
				print("ok")
			}
		}
		test()
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected boolean expression")
	})
}

func testForWithAndComparisonInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = 0
			var b = 1

			for a < 2 && b == 1 {
				print("ok")
				a++
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithOrComparisonInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = 0
			var b = 1

			for a < 2 || b == 1 {
				print("ok")
				a++

				if a == 2 {
					b = 2
				}
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithCountingVariableInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			for i := 0; i < 2; i++ {
				print("ok")
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithSeparateCountingVariableInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			i := 0
			for ; i < 2; i++ {
				print("ok")
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithSeparateCountingVariableAndSeparateIncrementInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			i := 0
			for ; i < 2; {
				print("ok")
				i++
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithSeparateCountingVariableAndSepareteConditionAndSeparateIncrementInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			i := 0
			for ; ; {
				if i >= 2 {
					break
				}
				print("ok")
				i++
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}

func testForWithNoConditionInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			i := 0

			for {
				if i >= 2 {
					break
				}
				print("ok")
				i++
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok\nok", output)
	})
}
