package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testStdStringsIndexSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "The quick brown fox jumps over the lazy dog"
	substr := "fox"

	testStringsFunc(t, transpilerCalloutFunc, "Index", []string{s, substr}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, fmt.Sprintf("%d", strings.Index(s, substr)), output)
	})
}

func testStdStringsContainsSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "The quick brown fox jumps over the lazy dog"
	substr := "fox"

	testStringsFunc(t, transpilerCalloutFunc, "Contains", []string{s, substr}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, true, strings.Contains(s, substr))
		require.EqualValues(t, "1", output)
	})
}

func testStdStringsJoinSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	sep := ","

	testStringsFunc(t, transpilerCalloutFunc, "Join", []string{`[]string{"1", "2", "3"}`, wrapInQuotes(sep)}, false, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.Join([]string{"1", "2", "3"}, sep), output)
	})
}

func testStringsFunc(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc, f string, args []string, quoteArgs bool, compare compareCallout) {
	transpilerCalloutFunc(t, func(dir string) (string, error) {
		file, err := copyStd("strings", dir)

		if err != nil {
			return "", err
		}

		if quoteArgs {
			/* Wrap parameters in quotes. */
			for i := range args {
				args[i] = wrapInQuotes(args[i])
			}
		}
		return `
			import strings "` + file + `"

			print(strings.` + f + `(` + strings.Join(args, ", ") + `))
		`, nil
	}, compare)
}

func wrapInQuotes(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}
