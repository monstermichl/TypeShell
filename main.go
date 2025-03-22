package main

import (
	"os"

	"github.com/monstermichl/typeshell/converters/bash"
	"github.com/monstermichl/typeshell/converters/batch"
	"github.com/monstermichl/typeshell/lexer"
	"github.com/monstermichl/typeshell/parser"
	"github.com/monstermichl/typeshell/transpiler"
)

func main() {
	testFile := "test.mss"
	data, _ := os.ReadFile(testFile)
	tokens, err := lexer.Tokenize(string(data))

	if err != nil {
		panic(err)
	}
	p := parser.New(tokens)
	prog, err := p.Parse()
	//fmt.Println(prog)
	if err != nil {
		panic(err)
	}
	i := transpiler.New(prog)

	// Dump batch file.
	batchConv := batch.New()
	dump, err := i.Transpile(batchConv)

	if err != nil {
		panic(err)
	}
	os.WriteFile(testFile+".bat", []byte(dump), 0777)

	// Dump bash file.
	bashConv := bash.New()
	dump, err = i.Transpile(bashConv)

	if err != nil {
		panic(err)
	}
	os.WriteFile(testFile+".sh", []byte(dump), 0777)
}
