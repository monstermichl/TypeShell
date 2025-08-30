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

func testStdStringsSplitSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	s := "ba-na-na ban-da-na ba-na-na"
	sep := "-"

	transpilerFunc(t, `
		import "strings"

		vals := strings.Split(`+strings.Join([]string{wrapInQuotes(s), wrapInQuotes(sep)}, ", ")+`)
		print(len(vals))

		for i, v := range(vals) {
			print(v)
		}
	`, func(output string, err error) {
		require.Nil(t, err)

		vals := strings.Split(s, sep)
		joined := strings.Join([]string{strconv.Itoa(len(vals)), strings.Join(vals, "\n")}, "\n")

		require.Equal(t, joined, output)
	})
}

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

func testStdStringsCutPrefixSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "Hello World"
	prefix := "Hel"

	testStringsFunc(t, transpilerCalloutFunc, "CutPrefix", []string{s, prefix}, true, func(output string, err error) {
		require.Nil(t, err)
		after, found := strings.CutPrefix(s, prefix)
		require.Equal(t, true, found)
		require.EqualValues(t, fmt.Sprintf("%s 1", after), output)
	})
}

func testStdStringsCutSuffixSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "Hello World"
	suffix := "rld"

	testStringsFunc(t, transpilerCalloutFunc, "CutSuffix", []string{s, suffix}, true, func(output string, err error) {
		require.Nil(t, err)
		after, found := strings.CutSuffix(s, suffix)
		require.Equal(t, true, found)
		require.EqualValues(t, fmt.Sprintf("%s 1", after), output)
	})
}

func testStdStringsCutSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "Hello World"
	sep := "o W"

	testStringsFunc(t, transpilerCalloutFunc, "Cut", []string{s, sep}, true, func(output string, err error) {
		require.Nil(t, err)
		before, after, found := strings.Cut(s, sep)
		require.Equal(t, true, found)
		require.EqualValues(t, fmt.Sprintf("%s %s 1", before, after), output)
	})
	sep = "Hel"

	testStringsFunc(t, transpilerCalloutFunc, "Cut", []string{s, sep}, true, func(output string, err error) {
		require.Nil(t, err)
		_, after, found := strings.Cut(s, sep)
		require.Equal(t, true, found)
		require.EqualValues(t, fmt.Sprintf("%s 1", after), output)
	})
	sep = "rld"

	testStringsFunc(t, transpilerCalloutFunc, "Cut", []string{s, sep}, true, func(output string, err error) {
		require.Nil(t, err)
		before, _, found := strings.Cut(s, sep)
		require.Equal(t, true, found)
		require.EqualValues(t, fmt.Sprintf("%s  1", before), output)
	})
	sep = "not included"

	testStringsFunc(t, transpilerCalloutFunc, "Cut", []string{s, sep}, true, func(output string, err error) {
		require.Nil(t, err)
		before, _, found := strings.Cut(s, sep)
		require.Equal(t, false, found)
		require.EqualValues(t, fmt.Sprintf("%s  0", before), output)
	})
}

func testStdStringsTrimPrefixSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "Hello World"
	prefix := "Hel"

	testStringsFunc(t, transpilerCalloutFunc, "TrimPrefix", []string{s, prefix}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.TrimPrefix(s, prefix), output)
	})
}

func testStdStringsTrimSuffixSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "Hello World"
	suffix := "rld"

	testStringsFunc(t, transpilerCalloutFunc, "TrimSuffix", []string{s, suffix}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.TrimSuffix(s, suffix), output)
	})
}

func testStdStringsTrimLeftSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "--000abc123xyz"
	cutset := "-0a1"

	testStringsFunc(t, transpilerCalloutFunc, "TrimLeft", []string{s, cutset}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.TrimLeft(s, cutset), output)
	})
}

func testStdStringsTrimRightSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "user123000---"
	cutset := "-03"

	testStringsFunc(t, transpilerCalloutFunc, "TrimRight", []string{s, cutset}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.TrimRight(s, cutset), output)
	})
}

func testStdStringsTrimSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "--000abcXYZ000--"
	cutset := "-0X"

	testStringsFunc(t, transpilerCalloutFunc, "Trim", []string{s, cutset}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.Trim(s, cutset), output)
	})
}

func testStdStringsTrimSpaceSuccess(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc) {
	s := "\t\n   Hello, GoLang!   \r\f "

	testStringsFunc(t, transpilerCalloutFunc, "TrimSpace", []string{s}, true, func(output string, err error) {
		require.Nil(t, err)
		require.EqualValues(t, strings.TrimSpace(s), output)
	})
}

func testStringsFunc(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc, f string, args []string, quoteArgs bool, compare compareCallout) {
	testStdFunc(t, transpilerCalloutFunc, "strings", f, args, quoteArgs, compare)
}
