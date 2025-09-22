package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func testDifferentNumberBasesSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		print(700)
		print(0o700)
		print(0O700)
		print(0700)
		print(0x700)
		print(0X700)
		print(0b1011)
		print(0B1011)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, fmt.Sprintf("%d\n%d\n%d\n%d\n%d\n%d\n%d\n%d", 700, 0o700, 0O700, 0700, 0x700, 0X700, 0b1011, 0B1011), output)
	})
}
