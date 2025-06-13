package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/monstermichl/typeshell/converters/bash"
	"github.com/monstermichl/typeshell/converters/batch"
	"github.com/monstermichl/typeshell/transpiler"
)

const (
	typeBatch string = "batch"
	typeBash  string = "bash"
)

var convMapping = map[string]transpiler.Converter{
	typeBatch: batch.New(),
	typeBash:  bash.New(),
}

type options struct {
	in         string
	out        string
	converters []transpiler.Converter
}

func parseOptions() options {
	args := os.Args
	options := options{}
	types := []string{}

	for k := range convMapping {
		types = append(types, k)
	}

	for i := 1; i < (len(args) - 1); i += 2 {
		cSwitch := args[i]
		cValue := args[i+1]

		switch cSwitch {
		case "-i", "--in":
			// Make sure input file exists.
			if stat, err := os.Stat(cValue); err != nil {
				panic(fmt.Errorf("input file %s doesn't exist", cValue))
			} else if stat.IsDir() {
				panic(fmt.Errorf("input %s is not a file", cValue))
			}
			options.in = cValue
		case "-o", "--out":
			// Make sure output path exists.
			if stat, err := os.Stat(cValue); err != nil {
				panic(fmt.Errorf("output path %s doesn't exist", cValue))
			} else if !stat.IsDir() {
				panic(fmt.Errorf("output path %s is not a directory", cValue))
			}
			options.out = cValue
		case "-t", "--type":
			conv, ok := convMapping[cValue]

			if !ok {
				panic(fmt.Errorf("unknown converter type %s. Allowed types are %s", cValue, strings.Join(types, ", ")))
			}
			options.converters = append(options.converters, conv)
		default:
			panic(fmt.Errorf("unknown option %s", cSwitch))
		}
	}

	if len(options.in) == 0 {
		panic("no input file provided (-i/--in)")
	} else if len(options.out) == 0 {
		panic("no output directory provided (-o/--out)")
	} else if len(options.converters) == 0 {
		panic("no output type provided (-t/--type)")
	}
	return options
}

func main() {
	options := parseOptions()
	t := transpiler.New()

	for _, conv := range options.converters {
		in := options.in
		dump, err := t.Transpile(in, conv)

		if err != nil {
			panic(err)
		}
		file := filepath.Base(in)
		file = file[0 : len(file)-len(filepath.Ext(in))] // Remove extension.

		os.WriteFile(filepath.Join(options.out, fmt.Sprintf("%s.%s", file, conv.Extension())), []byte(dump), 0777)
	}
}
