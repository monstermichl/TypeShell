package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func testLenSliceSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		a := []int{1, 2, 3}

		print(len(a))
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "3", output)
	})
}

func testLenStringSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(len("abc"))
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "3", output)
	})
}

func testCopySuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		s1 := []int{}
		s2 := []int{1, 2, 3}

		a := copy(s1, s2)

		print(a)

		for i := 0; i < len(s1); i++ {
			print(s1[i])
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "3\n1\n2\n3", output)
	})
}

func testReadSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	file := "read-test.txt"
	content := "Hello World"
	os.WriteFile(file, []byte(content), 0x777)
	defer os.Remove(file)

	transpilerFunc(t, `
		`+fmt.Sprintf(`a := read("%s")`, file)+`
		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, content, output)
	})
}

func testWriteSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	file := "read-test.txt"
	content := "Hello Moon"
	defer os.Remove(file)

	transpilerFunc(t, `
		`+fmt.Sprintf(`write("%s", "%s")`, file, content)+`
		`+fmt.Sprintf(`a := read("%s")`, file)+`
		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, content, output)
	})
}

func testWriteAppendSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	file := "read-test.txt"
	content := "Hello Moon"
	defer os.Remove(file)

	transpilerFunc(t, `
		`+fmt.Sprintf(`write("%s", "%s")`, file, content)+`
		`+fmt.Sprintf(`write("%s", "%s", true)`, file, content)+`
		`+fmt.Sprintf(`a := read("%s")`, file)+`
		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, fmt.Sprintf("%s\n%s", content, content), output)
	})
}

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

func testLenSliceInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			a := []int{1, 2, 3}

			print(len(a))
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "3", output)
	})
}

func testLenStringInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			print(len("abc"))
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "3", output)
	})
}

func testCopyInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			s1 := []int{}
			s2 := []int{1, 2, 3}

			a := copy(s1, s2)

			print(a)

			for i := 0; i < len(s1); i++ {
				print(s1[i])
			}
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "3\n1\n2\n3", output)
	})
}

func testReadInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	file := "read-test.txt"
	content := "Hello World"
	os.WriteFile(file, []byte(content), 0x777)
	defer os.Remove(file)

	transpilerFunc(t, `
		func test() {
			`+fmt.Sprintf(`a := read("%s")`, file)+`
			print(a)
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, content, output)
	})
}

func testWriteInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	file := "read-test.txt"
	content := "Hello Moon"
	defer os.Remove(file)

	transpilerFunc(t, `
		func test() {
			`+fmt.Sprintf(`write("%s", "%s")`, file, content)+`
			`+fmt.Sprintf(`a := read("%s")`, file)+`
			print(a)
		}
		test()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, content, output)
	})
}

func testPanicInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func test() {
			var a = 1

			if a == 1 {
				panic("panic")
			}
			print("got through")
		}
		test()
	`, func(output string, err error) {
		require.Error(t, err)
		require.Equal(t, "panic: panic", output)
	})
}
