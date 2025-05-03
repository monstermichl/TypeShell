package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testIntEqualSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(2 == 2)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testIntNotEqualSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(2 != 1)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testIntLessSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(1 < 2)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testIntLessOrEqualSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(2 <= 2)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testIntGreaterSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(2 > 1)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testIntGreaterOrEqualSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(2 >= 2)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testComplexIntComparisonSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = (10 == 10) && (5 != 3) && (8 > 2) && (3 < 4) && (6 >= 6) && (7 <= 9) || !(2 > 5)

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testStringEqualSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print("test" == "test")
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testStringNotEqualSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print("test" != "no test")
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testBooleanEqualSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(true == true)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testBooleanNotEqualSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(true != false)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}
