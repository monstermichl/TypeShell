package main

import (
	"os"

	"github.com/monstermichl/typeshell/converters/bash"
	"github.com/monstermichl/typeshell/converters/batch"
	"github.com/monstermichl/typeshell/transpiler"
)

func main() {
	testFile := "test.mss"
	t := transpiler.New()

	// Dump batch file.
	batchConv := batch.New()
	dump, err := t.Transpile(testFile, batchConv)

	if err != nil {
		panic(err)
	}
	os.WriteFile(testFile+".bat", []byte(dump), 0777)

	// Dump bash file.
	bashConv := bash.New()
	dump, err = t.Transpile(testFile, bashConv)

	if err != nil {
		panic(err)
	}
	os.WriteFile(testFile+".sh", []byte(dump), 0777)
}
