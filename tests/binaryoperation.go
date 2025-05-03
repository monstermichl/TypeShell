package tests

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func testAdditionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 2
		var b = 3
		var c = a + b

		print(c)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, strconv.Itoa(2+3), output)
	})
}

func testSubtractionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 2
		var b = 3
		var c = a - b

		print(c)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, strconv.Itoa(2-3), output)
	})
}

func testMultiplicationSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 2
		var b = 3
		var c = a * b

		print(c)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, strconv.Itoa(2*3), output)
	})
}

func testDivisionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 3
		var b = 2
		var c = a / b

		print(c)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, strconv.Itoa(3/2), output)
	})
}

func testModuloSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 2
		var b = 3
		var c = a % b

		print(c)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, strconv.Itoa(2%3), output)
	})
}

func testMoreComplexCalculationSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 5
		var b = 4
		var c = 3
		var d = 2
		var e = a * b + c / d

		print(e)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, strconv.Itoa(5*4+3/2), output)
	})
}

func testMoreComplexCalculationWithBracketsSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 5
		var b = 4
		var c = 3
		var d = 2
		var e = a * (b + c) / d

		print(e)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, strconv.Itoa(5*(4+3)/2), output)
	})
}

func testComplexCalculationSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = (145 + 37 * 2) - 18 + 64 / 4 * 3 - (250 % 23 + 11) + 7 * (81 - 9) / 6 + 999

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, strconv.Itoa((145+37*2)-18+64/4*3-(250%23+11)+7*(81-9)/6+999), output)
	})
}

func testCompoundAssignmentAdditionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 2
		a += 2

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "4", output)
	})
}

func testCompoundAssignmentSubtractionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 2
		a -= 2

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0", output)
	})
}

func testCompoundAssignmentMultiplicationSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 2
		a *= 2

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "4", output)
	})
}

func testCompoundAssignmentDivisionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 4
		a /= 2

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "2", output)
	})
}

func testCompoundAssignmentModuloSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		var a = 4
		a %= 2

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "0", output)
	})
}
