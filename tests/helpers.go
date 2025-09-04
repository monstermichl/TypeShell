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
	exe, err := os.Executable()
	require.Nil(t, err)

	exePath := filepath.Dir(exe)

	// Copy std to executable path.
	err = copyStd(exePath)
	require.Nil(t, err)

	dir := filepath.Join(exePath, t.Name())

	err = os.MkdirAll(dir, 0700)
	require.Nil(t, err)
	defer os.RemoveAll(dir) // Make sure test dir is removed after test case execution.

	trans := transpiler.New()
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

func copyStd(dstDir string) error {
	stdFolder := "std"
	srcDir := filepath.Join("..", stdFolder)
	file, err := os.Open(srcDir)

	if err != nil {
		return err
	}
	dirEntries, err := file.ReadDir(0)

	if err != nil {
		return err
	}
	dstDir = filepath.Join(dstDir, stdFolder)
	err = os.MkdirAll(dstDir, 0700)

	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		name := dirEntry.Name()

		if strings.HasSuffix(name, ".tsh") {
			err := copyFile(name, srcDir, dstDir)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func testStdFunc(t *testing.T, transpilerCalloutFunc transpilerCalloutFunc, stdLib string, f string, args []string, quoteArgs bool, compare compareCallout) {
	transpilerCalloutFunc(t, func(dir string) (string, error) {
		if quoteArgs {
			/* Wrap parameters in quotes. */
			for i := range args {
				args[i] = wrapInQuotes(args[i])
			}
		}
		return `
			import "` + stdLib + `"

			print(` + stdLib + `.` + f + `(` + strings.Join(args, ", ") + `))
		`, nil
	}, compare)
}

func wrapInQuotes(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}
