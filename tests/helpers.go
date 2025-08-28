package tests

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/monstermichl/typeshell/converters/bash"
	"github.com/monstermichl/typeshell/converters/batch"
	"github.com/monstermichl/typeshell/transpiler"
	"github.com/stretchr/testify/require"
)

type sourceCallout func(dir string) (string, error)
type compareCallout func(output string, err error)
type transpilerFunc func(t *testing.T, source string, compare compareCallout)
type transpilerCalloutFunc func(t *testing.T, callout sourceCallout, compare compareCallout)

func transpileFunc(t *testing.T, source sourceCallout, targetFileName string, converter transpiler.Converter, compare compareCallout) {
	trans := transpiler.New()
	dir, err := os.MkdirTemp("", "typeshell_tests")

	require.Nil(t, err)
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "test.tsh")
	outputString := ""
	src, err := source(dir)

	if err == nil {
		var code string
		err = os.WriteFile(file, []byte(src), 0700)

		require.Nil(t, err)
		code, err = trans.Transpile(file, converter)
		output := []byte{}

		// If transpilation was successful, run the code.
		if err == nil {
			targetFile := filepath.Join(dir, targetFileName)
			err = os.WriteFile(targetFile, []byte(code), 0700)

			require.Nil(t, err)
			cmd := exec.Command(targetFile)
			output, err = cmd.Output()
		}
		outputString = string(output)
		outputString = strings.ReplaceAll(outputString, "\r\n", "\n")
		outputString = strings.TrimSpace(outputString)
	}
	compare(outputString, err)
}

func transpile(t *testing.T, source string, targetFileName string, converter transpiler.Converter, compare compareCallout) {
	transpileFunc(t, func(_ string) (string, error) {
		return source, nil
	}, targetFileName, converter, compare)
}

func transpileBash(t *testing.T, source string, compare compareCallout) {
	transpile(t, source, "test.sh", bash.New(), compare)
}

func transpileBashFunc(t *testing.T, source sourceCallout, compare compareCallout) {
	transpileFunc(t, source, "test.sh", bash.New(), compare)
}

func transpileBatch(t *testing.T, source string, compare compareCallout) {
	transpile(t, source, "test.bat", batch.New(), compare)
}

func transpileBatchFunc(t *testing.T, source sourceCallout, compare compareCallout) {
	transpileFunc(t, source, "test.bat", batch.New(), compare)
}

func shortenError(err error) error {
	if err != nil {
		s := err.Error()

		if matches := regexp.MustCompile(`(.+)\s+at row`).FindStringSubmatch(s); matches != nil {
			s = matches[1]
		}
		err = errors.New(s)
	}
	return err
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

func copyStd(fileName string, dstDir string) (string, error) {
	fileParts := strings.Split(fileName, ".")
	fileName = fmt.Sprintf("%s.tsh", fileParts[0])

	return fileName, copyFile(fileName, filepath.Join("..", "std"), dstDir)
}

func testStdFunc(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc, std string, f string, args []string, quoteArgs bool, compare compareCallout) {
	transpilerCalloutFunc(t, func(dir string) (string, error) {
		file, err := copyStd(std, dir)

		if err != nil {
			return "", err
		}

		if quoteArgs {
			/* Wrap parameters in quotes. */
			for i := range args {
				args[i] = wrapInQuotes(args[i])
			}
		}
		return `
			import strings "` + file + `"

			print(strings.` + f + `(` + strings.Join(args, ", ") + `))
		`, nil
	}, compare)
}

func wrapInQuotes(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}
