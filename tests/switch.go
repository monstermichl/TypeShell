package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testSwitchWithBoolSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		a := true

		switch a {
		case true:
			print("ok")
		case false:
			print("nok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testSwitchWithBoolDefaultSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		a := true

		switch a {
		case false:
			print("nok")
		default:
			print("ok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testSwitchWithImplicitBoolSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		a := true

		switch {
		case true:
			print("ok")
		default:
			print("nok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testSwitchWithComparisonsSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		a := 1
		b := 2

		switch {
		case a == 1 && b == 1:
			print("nok")
		case a == 1 && b > 1:
			print("ok")
		default:
			print("nok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testSwitchOnlyDefaultSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		switch {
		default:
			print("ok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}

func testSwitchStringsSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		a := "c"

		switch a {
		case "a":
			print("nok")
		case "b":
			print("nok")
		case "c":
			print("ok")
		default:
			print("nok")
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "ok", output)
	})
}
