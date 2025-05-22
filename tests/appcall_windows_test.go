package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirCallSuccess(t *testing.T) {
	transpileBatchFunc(t, func(dir string) (string, error) {
		return `
			` + fmt.Sprintf(`var a, code = @dir("/B", "%s")`, strings.ReplaceAll(dir, `\`, `\\`)) + `

			print(a, code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.bat\ntest.tsh 0", output)
	})
}

func TestDirCallPipeToFindstrCallSuccess(t *testing.T) {
	transpileBatchFunc(t, func(dir string) (string, error) {
		return `
			` + fmt.Sprintf(`var a, code = @dir("/B", "%s") | @findstr(".tsh")`, strings.ReplaceAll(dir, `\`, `\\`)) + `

			print(a, code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.tsh 0", output)
	})
}
