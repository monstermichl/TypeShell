package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testSingleImportSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		return `
			import "strings"
			print(strings.Contains("Hello World", "Wor"))
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testMultiImportSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		return `import (
			strings1 "strings"
			strings2 "strings"
			)
			print(strings1.Contains("Hello World", "Wor"))
			print(strings2.HasPrefix("Hello World", "Hel"))
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1\n1", output)
	})
}
