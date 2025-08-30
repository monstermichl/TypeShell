package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testStdOsShellSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc, expectation string) {
	testOsFunc(t, transpilerCalloutFunc, "Shell", []string{}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, expectation, output)
	})
}

func testOsFunc(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc, f string, args []string, quoteArgs bool, compare compareCallout) {
	testStdFunc(t, transpilerCalloutFunc, "os", f, args, quoteArgs, compare)
}
