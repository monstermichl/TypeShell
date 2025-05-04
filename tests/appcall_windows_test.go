package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirCallSuccess(t *testing.T) {
	transpileBatchFunc(t, func(dir string) string {
		return `
			` + fmt.Sprintf(`var a = @dir("/B", "%s")`, strings.ReplaceAll(dir, `\`, `\\`)) + `

			print(a)
		`
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.bat\ntest.tsh", output)
	})
}

func TestDirCallPipeToFindstrCallSuccess(t *testing.T) {
	transpileBatchFunc(t, func(dir string) string {
		return `
			` + fmt.Sprintf(`var a = @dir("/B", "%s") | @findstr(".tsh")`, strings.ReplaceAll(dir, `\`, `\\`)) + `

			print(a)
		`
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.tsh", output)
	})
}
