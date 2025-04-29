package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testIfComparisonSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 1

		if a == 1 {
			print("ok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testNonBoolIfConditionFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		if "1" {
			print("ok")
		}
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected boolean expression")
	})
}

func testIfWithAndComparisonSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 1
		var b = 2

		if a == 1 && b == 2 {
			print("ok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testIfWithOrComparisonSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 1
		var b = 2

		if a == 1 || b == 2 {
			print("ok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testElseIfSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 1
		var b = 2

		if a == 1 && b == 1 {
			print("nok")
		} else if b == 2 {
		 	print("ok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testElseSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 1
		var b = 3

		if a == 1 && b == 1 {
			print("nok")
		} else if b == 2 {
		 	print("nok")
		} else {
		 	print("ok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testIfComparisonInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = 1

			if a == 1 {
				print("ok")
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testIfWithAndComparisonInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = 1
			var b = 2

			if a == 1 && b == 2 {
				print("ok")
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testIfWithOrComparisonInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = 1
			var b = 2

			if a == 1 || b == 2 {
				print("ok")
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testElseIfInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = 1
			var b = 2

			if a == 1 && b == 1 {
				print("nok")
			} else if b == 2 {
				print("ok")
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testElseInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = 1
			var b = 3

			if a == 1 && b == 1 {
				print("nok")
			} else if b == 2 {
				print("nok")
			} else {
				print("ok")
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}
