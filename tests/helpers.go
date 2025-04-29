package tests

import (
	"errors"
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

type compareCallout func(output string, err error)
type transpilerFunc func(t *testing.T, source string, compare compareCallout)

func transpile(t *testing.T, source string, targetFileName string, converter transpiler.Converter, compare compareCallout) {
	trans := transpiler.New()
	dir, err := os.MkdirTemp("", "typeshell_tests")

	require.Nil(t, err)
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "test.tsh")
	err = os.WriteFile(file, []byte(source), 0x777)

	require.Nil(t, err)
	code, err := trans.Transpile(file, converter)
	output := []byte{}

	// If transpilation was successful, run the code.
	if err == nil {
		targetFile := filepath.Join(dir, targetFileName)
		err = os.WriteFile(targetFile, []byte(code), 0x777)

		require.Nil(t, err)
		cmd := exec.Command(targetFile)
		output, err = cmd.Output()
	}
	compare(strings.TrimSpace(string(output)), err)
}

func transpileBash(t *testing.T, source string, compare compareCallout) {
	transpile(t, source, "test.sh", bash.New(), compare)
}

func transpileBatch(t *testing.T, source string, compare compareCallout) {
	transpile(t, source, "test.bat", batch.New(), compare)
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
