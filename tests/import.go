package tests

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func testSingleImportSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		file := "strings.tsh"
		err := copyFile(file, filepath.Join("..", "std"), dir)

		if err != nil {
			return "", err
		}
		return fmt.Sprintf(`import strings "%s"`, file) + `
			print(strings.Contains("Hello World", "Wor"))
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testMultiImportSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		file := "strings.tsh"
		err := copyFile(file, filepath.Join("..", "std"), dir)

		if err != nil {
			return "", err
		}
		return `import (
			` + fmt.Sprintf(`strings1 "%s"`, file) + `
			` + fmt.Sprintf(`strings2 "%s"`, file) + `
			)
			print(strings1.Contains("Hello World", "Wor"))
			print(strings2.HasPrefix("Hello World", "Hel"))
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1\n1", output)
	})
}

func copyFile(fileName string, srcDir string, dstDir string) error {
	src, err := os.Open(filepath.Join(srcDir, fileName))

	if err != nil {
		return err
	}
	dst, err := os.Create(filepath.Join(dstDir, fileName))

	if err != nil {
		return err
	}
	_, err = io.Copy(dst, src)
	return err
}
