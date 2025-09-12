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
