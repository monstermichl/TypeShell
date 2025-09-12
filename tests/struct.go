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
