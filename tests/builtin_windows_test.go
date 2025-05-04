package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLenSliceSuccess(t *testing.T) {
	testLenSliceSuccess(t, transpileBatch)
}

func TestLenStringSuccess(t *testing.T) {
	testLenStringSuccess(t, transpileBatch)
}

func TestCopySuccess(t *testing.T) {
	testCopySuccess(t, transpileBatch)
}

func TestExistsSuccess(t *testing.T) {
	transpileBatchFunc(t, func(dir string) string {
		return fmt.Sprintf(`print(exists("%s"))`, strings.ReplaceAll(dir, `\`, `\\`))
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func TestReadSuccess(t *testing.T) {
	testReadSuccess(t, transpileBatch)
}

func TestWriteSuccess(t *testing.T) {
	testWriteSuccess(t, transpileBatch)
}

func TestWriteAppendSuccess(t *testing.T) {
	testWriteAppendSuccess(t, transpileBatch)
}

func TestPanicSuccess(t *testing.T) {
	testPanicSuccess(t, transpileBatch)
}

func TestLenSliceInFunctionSuccess(t *testing.T) {
	testLenSliceInFunctionSuccess(t, transpileBatch)
}

func TestLenStringInFunctionSuccess(t *testing.T) {
	testLenStringInFunctionSuccess(t, transpileBatch)
}

func TestCopyInFunctionSuccess(t *testing.T) {
	testCopyInFunctionSuccess(t, transpileBatch)
}

func TestExistsInFunctionSuccess(t *testing.T) {
	transpileBatchFunc(t, func(dir string) string {
		return `
			func test() {
			` + fmt.Sprintf(`print(exists("%s"))`, strings.ReplaceAll(dir, `\`, `\\`)) + `
			}
			test()
		`
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func TestReadInFunctionSuccess(t *testing.T) {
	testReadInFunctionSuccess(t, transpileBatch)
}

func TestWriteInFunctionSuccess(t *testing.T) {
	testWriteInFunctionSuccess(t, transpileBatch)
}

func TestPanicInFunctionSuccess(t *testing.T) {
	testPanicInFunctionSuccess(t, transpileBatch)
}
