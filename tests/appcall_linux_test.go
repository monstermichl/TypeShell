package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLsCallSuccess(t *testing.T) {
	transpileBashFunc(t, func(dir string) string {
		return `
			` + fmt.Sprintf(`var a = @ls("%s")`, dir) + `

			print(a)
		`
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.sh\ntest.tsh", output)
	})
}

func TestLsCallPipeToGrepCallSuccess(t *testing.T) {
	transpileBashFunc(t, func(dir string) string {
		return `
			` + fmt.Sprintf(`var a = @ls("%s") | @grep(".tsh")`, dir) + `

			print(a)
		`
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test.tsh", output)
	})
}
