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

func testDefineSliceRowValuesSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = []int{
			1,
			2,
		}

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

func testReassignSliceSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = []int{1, 2}
		a[3] = 4

		print(len(a))
		print(a[0], a[1], a[2], a[3])

		a = []int{5, 6}
		print(len(a))
		print(a[0], a[1])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "4\n1 2 0 4\n2\n5 6", output)
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

func testReassignSliceInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = []int{1, 2}
			a[3] = 4

			print(len(a))
			print(a[0], a[1], a[2], a[3])

			a = []int{5, 6}
			print(len(a))
			print(a[0], a[1])
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "4\n1 2 0 4\n2\n5 6", output)
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

func testSliceReturnedFromFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		i := 0
		func test() []string {
			return []string{"test" + itoa(i)}
		}
		s1 := test()
		i++
		s2 := test()
		print(s1[0])
		print(s2[0])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test0\ntest1", output)
	})
}

func testComplexSliceOperationsSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test1() string {
			return "hello world"
		}

		func test2() string {
			return "hello mars"
		}

		func test3() []string {
			return []string{test1(), test2()}
		}

		func test4(nextString string) ([]string, int) {
			s := test3()
			s[3] = nextString

			return s, len(s)
		}
		slice1, amount := test4("hello mum")
		print("amount 1: " + itoa(amount))

		var slice2 []string
		amount = copy(slice2, slice1)
		print("amount 2: " + itoa(amount))

		for i, v := range slice2 {
			print(v)
		}
		print(slice1[0])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "amount 1: 4\namount 2: 4\nhello world\nhello mars\n\nhello mum\nhello world", output)
	})
}
