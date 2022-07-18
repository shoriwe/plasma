package main

import (
	"github.com/shoriwe/gplasma/pkg/compiler"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
)

func executeFiles() {
	for _, arg := range os.Args {
		switch arg {
		case "-h", "--help":
			help()
		}
	}
	files := make([][]byte, 0, len(os.Args[1:]))
	for _, file := range os.Args[1:] {
		contents, readError := os.ReadFile(file)
		if readError != nil {
			onError(file, readError)
		}
		files = append(files, contents)
	}
	plasma := vm.NewVM(os.Stdin, os.Stdout, os.Stderr)
	for index, file := range files {
		bytecode, compileError := compiler.Compile(string(file))
		if compileError != nil {
			onError(os.Args[1:][index], compileError)
		}
		_, errorChan, _ := plasma.Execute(bytecode)
		executeError := <-errorChan
		if executeError != nil {
			onError(os.Args[1:][index], executeError)
		}
	}
}
