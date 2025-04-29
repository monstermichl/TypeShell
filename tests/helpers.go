package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/monstermichl/typeshell/converters/bash"
	"github.com/monstermichl/typeshell/converters/batch"
	"github.com/monstermichl/typeshell/transpiler"
	"github.com/stretchr/testify/require"
)

type compareCallout func(output string)
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

	require.Nil(t, err)
	targetFile := filepath.Join(dir, targetFileName)
	err = os.WriteFile(targetFile, []byte(code), 0x777)

	require.Nil(t, err)
	cmd := exec.Command(targetFile)
	output, err := cmd.Output()

	require.Nil(t, err)
	compare(strings.TrimSpace(string(output)))
}

func transpileBash(t *testing.T, source string, compare compareCallout) {
	transpile(t, source, "test.sh", bash.New(), compare)
}

func transpileBatch(t *testing.T, source string, compare compareCallout) {
	transpile(t, source, "test.bat", batch.New(), compare)
}
