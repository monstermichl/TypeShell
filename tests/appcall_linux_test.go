package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLsCallSuccess(t *testing.T) {
	transpileBashFunc(t, func(dir string) (string, error) {
		return `
			` + fmt.Sprintf(`var a, code = @ls("%s")`, dir) + `

			print(a, code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.sh\ntest.tsh 0", output)
	})
}

func TestLsCallPipeToGrepCallSuccess(t *testing.T) {
	transpileBashFunc(t, func(dir string) (string, error) {
		return `
			` + fmt.Sprintf(`var a, code = @ls("%s") | @grep(".tsh")`, dir) + `

			print(a, code)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.tsh 0", output)
	})
}
