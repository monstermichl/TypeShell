package tests

import (
	"fmt"
	"strconv"
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

func testStdStringsHasPrefixSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "Hello World"
	prefix := "Hel"

	testStringsFunc(t, transpilerCalloutFunc, "HasPrefix", []string{s, prefix}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, true, strings.HasPrefix(s, prefix))
		require.EqualValues(t, "1", output)
	})
}

func testStdStringsHasSuffixSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "Hello World"
	suffix := "rld"

	testStringsFunc(t, transpilerCalloutFunc, "HasSuffix", []string{s, suffix}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, true, strings.HasSuffix(s, suffix))
		require.EqualValues(t, "1", output)
	})
}

func testStdStringsCountSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "banana bandana banana"
	substr := "na"

	testStringsFunc(t, transpilerCalloutFunc, "Count", []string{s, substr}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strconv.Itoa(strings.Count(s, substr)), output)
	})
}

// TODO: Add test for Split-function.

func testStdStringsRepeatSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "na"
	count := 8

	testStringsFunc(t, transpilerCalloutFunc, "Repeat", []string{wrapInQuotes(s), strconv.Itoa(count)}, false, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.Repeat(s, count), output)
	})
}

func testStdStringsReplaceSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "banana bandana banana"
	old := "na"
	new := "no"
	n := 3

	testStringsFunc(t, transpilerCalloutFunc, "Replace", []string{wrapInQuotes(s), wrapInQuotes(old), wrapInQuotes(new), strconv.Itoa(n)}, false, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.Replace(s, old, new, n), output)
	})
}

func testStdStringsReplaceAllSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "banana bandana banana"
	old := "na"
	new := "no"

	testStringsFunc(t, transpilerCalloutFunc, "ReplaceAll", []string{s, old, new}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.ReplaceAll(s, old, new), output)
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
