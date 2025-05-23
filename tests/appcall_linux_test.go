package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLsCallSuccess(t *testing.T) {
	transpileBashFunc(t, func(dir string) (string, error) {
		return `
			` + fmt.Sprintf(`var stdout, stderr, code = @ls("%s")`, dir) + `

			print(stdout, code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.sh\ntest.tsh 0", output)
	})
}

func TestLsCallPipeToGrepCallSuccess(t *testing.T) {
	transpileBashFunc(t, func(dir string) (string, error) {
		return `
			` + fmt.Sprintf(`var stdout, stderr, code = @ls("%s") | @grep(".tsh")`, dir) + `

			print(stdout, code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.tsh 0", output)
	})
}

func TestLsCallFail(t *testing.T) {
	transpileBashFunc(t, func(dir string) (string, error) {
		return `
			var stdout, stderr, code = @ls("not-present-dir")

			print(code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.NotEqual(t, "0", output)
	})
}
