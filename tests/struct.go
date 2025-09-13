package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testDeclareAndDefineStructSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		var s myStruct

		s.a = "Hello"
		s.b = "World"

		print(s.a, s.b, s.c, s.d)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello World 0 0", output)
	})
}

func testDeclareAndDefineStructWithValuesSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		s := myStruct{
			a: "Hello",
			b: "World",
		}

		print(s.a, s.b, s.c, s.d)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello World 0 0", output)
	})
}

func testDeclareAndDefineStructWithValuesOneLineSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		s := myStruct{a: "Hello", b: "World"}

		print(s.a, s.b, s.c, s.d)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello World 0 0", output)
	})
}

func testDeclareAndDefineStructSliceSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		var s myStruct
		var sl []myStruct

		s.a = "Hello"
		s.b = "World"
		s.d = 2

		sl[1] = s

		for i, val := range sl {
			print(val.a, val.b, val.c, val.d)
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0 0\nHello World 0 2", output)
	})
}

func testPassStructToFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
		}

		func test(x myStruct) {
			print(x.a, x.b)

			x.a = "Bye"
			x.b = "Mars"

			print(x.a, x.b)
		}
		var s myStruct

		s.a = "Hello"
		s.b = "World"

		print(s.a, s.b)
		test(s)
		print(s.a, s.b)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello World\nHello World\nBye Mars\nHello World", output)
	})
}

func testReturnDifferentStructsSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		count := 0

		func test() myStruct {
			s := myStruct{}
			s.a = itoa(count)

			return s
		}

		s1 := test()
		print(s1.a)
		count++
		s2 := test()
		print(s1.a)
		print(s2.a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0\n0\n1", output)
	})
}

func testNestedStructSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a string
		}

		type myStruct2 struct {
			a string
			b myStruct
		}

		s2 := myStruct2{a: "Hello", b: myStruct{a: "World"}}
		s := s2.b

		print(s2.a, s.a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello World", output)
	})
}

func testStructEvaluationChainingSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a string
		}

		type myStruct2 struct {
			a string
			b myStruct
		}

		s := myStruct2{a: "Hello", b: myStruct{a: "World"}}

		print(s.a, s.b.a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello World", output)
	})
}

func testStructAssignmentChainingSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a string
		}

		type myStruct2 struct {
			a string
			b myStruct
		}

		s := myStruct2{a: "Hello", b: myStruct{a: "World"}}
		s.b.a = "Mars"

		print(s.a, s.b.a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello Mars", output)
	})
}

func testStructFieldAssignedTwiceInInitializationFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		s := myStruct{a: "Hello", a: "World"}
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "a value has already been assigned to a")
	})
}

func testStructUnknownFieldInInitializationFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		s := myStruct{a: "Hello", x: "World"}
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "struct field x doesn't exist")
	})
}

func testStructFieldWrongTypeAssignmentInInitializationFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		s := myStruct{a: "Hello", b: 1}
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected string value but got int")
	})
}

func testStructUnknownFieldInAssignmentFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		var s myStruct

		s.x = "Test"
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "struct field x doesn't exist")
	})
}

func testStructFieldWrongTypeAssignmentInAssignmentFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myStruct struct {
			a, b string
			c    bool
			d    int
		}
		var s myStruct

		s.b = 1
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected string value but got int")
	})
}
