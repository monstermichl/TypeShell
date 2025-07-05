package tests

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirCallSuccess(t *testing.T) {
	transpileBatchFunc(t, func(dir string) (string, error) {
		return `
			` + fmt.Sprintf(`var stdout, stderr, code = @dir("/B", "%s")`, strings.ReplaceAll(dir, `\`, `\\`)) + `

			print(stdout, code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.bat\ntest.tsh 0", output)
	})
}

func TestDirCallPipeToFindstrCallSuccess(t *testing.T) {
	transpileBatchFunc(t, func(dir string) (string, error) {
		return `
			` + fmt.Sprintf(`var stdout, stderr, code = @dir("/B", "%s") | @findstr(".tsh")`, strings.ReplaceAll(dir, `\`, `\\`)) + `

			print(stdout, code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.tsh 0", output)
	})
}

func TestBatFileFromSubDirCallPipeToFindstrCallSuccess(t *testing.T) {
	transpileBatchFunc(t, func(dir string) (string, error) {
		subdir := path.Join(dir, "subdir")
		err := os.Mkdir(subdir, 0700)

		if err != nil {
			return "", err
		}
		batFile := path.Join(subdir, "bat.bat")
		err = os.WriteFile(batFile, []byte("dir /b .."), 0700)

		if err != nil {
			return "", err
		}
		return `
			` + fmt.Sprintf(`var stdout, stderr, code = @"%s"() | @findstr(".tsh")`, strings.ReplaceAll(strings.ReplaceAll(batFile, `/`, `\`), `\`, `\\`)) + `

			print(stdout, code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.tsh 0", output)
	})
}

func TestDirCallFail(t *testing.T) {
	transpileBatchFunc(t, func(dir string) (string, error) {
		return `
			var stdout, stderr, code = @dir("/B", "not-present-dir")

			print(code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.NotEqual(t, "0", output)
	})
}
