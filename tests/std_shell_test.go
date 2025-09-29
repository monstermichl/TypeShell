package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testStdShellGetSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc, expectation string) {
	testShellFunc(t, transpilerCalloutFunc, "Get", []string{}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, expectation, output)
	})
}

func testShellFunc(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc, f string, args []string, quoteArgs bool, compare compareCallout) {
	testStdFunc(t, transpilerCalloutFunc, "shell", f, args, quoteArgs, compare)
}
