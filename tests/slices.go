package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testDefineSliceSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = []int{1, 2}

		print(a[0], a[1])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2", output)
	})
}

func testSliceAssignValuesSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = []int{1, 2}

		a[0] = 3
		a[1] = 4
		a[2] = 5

		print(a[0], a[1], a[2])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "3 4 5", output)
	})
}

func testSliceAssignUndefinedSubscriptSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = []int{1, 2}

		a[3] = 4

		print(a[0], a[1], a[2], a[3])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2 0 4", output)
	})
}

func testSliceLengthSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = []int{1, 2}
		
		print(len(a))
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "2", output)
	})
}

func testIterateSliceSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = []int{1, 2}

		for i := 0; i < len(a); i++ {
			print(a[i])
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1\n2", output)
	})
}

func testCopySliceSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = []int{1, 2}
		var b = []int{}

		amount := copy(b, a)
		print(amount)

		for i := 0; i < len(b); i++ {
			print(b[i])
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "2\n1\n2", output)
	})
}

func testDefineSliceInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = []int{1, 2}

			print(a[0], a[1])
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2", output)
	})
}

func testSliceAssignValuesInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = []int{1, 2}

			a[0] = 3
			a[1] = 4
			a[2] = 5

			print(a[0], a[1], a[2])
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "3 4 5", output)
	})
}

func testSliceAssignUndefinedSubscriptInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = []int{1, 2}

			a[3] = 4

			print(a[0], a[1], a[2], a[3])
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1 2 0 4", output)
	})
}

func testSliceLengthInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = []int{1, 2}
			
			print(len(a))
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "2", output)
	})
}

func testIterateSliceInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = []int{1, 2}

			for i := 0; i < len(a); i++ {
				print(a[i])
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1\n2", output)
	})
}

func testCopySliceInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = []int{1, 2}
			var b = []int{}

			amount := copy(b, a)
			print(amount)

			for i := 0; i < len(b); i++ {
				print(b[i])
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "2\n1\n2", output)
	})
}
