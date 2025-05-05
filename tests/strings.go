package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testStringConcatSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		s1 := "Hello"
		s2 := "World"
		s3 := s1 + " " + s2

		print(s3)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello World", output)
	})
}

func testStringLengthSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		s := "test"

		print(len(s))
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "4", output)
	})
}

func testStringSingleSubscriptSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		s := "test"

		print(s[1], s[3])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "e t", output)
	})
}

func testStringStartSubscriptSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		s := "test"

		print(s[2:])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test"[2:], output)
	})
}

func testStringEndSubscriptSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		s := "test"

		print(s[:1])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test"[:1], output)
	})
}

func testStringRangeSubscriptSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		s := "test"

		print(s[1:3])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test"[1:3], output)
	})
}

func testStringRangeNoIndicesSubscriptSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		s := "test"

		print(s[:])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test"[:], output)
	})
}

func testItoaSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print("Hello World " + itoa(24))
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello World 24", output)
	})
}
