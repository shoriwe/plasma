package performance

import (
	"bytes"
	"fmt"
	"github.com/shoriwe/gplasma/pkg/compiler"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
	"path/filepath"
	"testing"
)

const (
	fibonacci = "fibo.pm"
	whileTest = "while.pm"
	forTest   = "for.pm"
)

func test(t *testing.T, script string) {
	currentDir, currentDirGetError := os.Getwd()
	if currentDirGetError != nil {
		t.Fatal(currentDirGetError)
		return
	}
	script = filepath.Join(currentDir, script)
	fileHandler, openError := os.Open(script)
	if openError != nil {
		t.Fatal(openError)
		return
	}
	bytecode, compilationError := compiler.Compile(reader.NewStringReaderFromFile(fileHandler))
	if compilationError != nil {
		t.Fatal(compilationError)
		return
	}
	output := bytes.NewBuffer(make([]byte, 0))
	plasmaVm := vm.NewPlasmaVM(nil, output, output)
	executionError, success := plasmaVm.Execute(nil, bytecode)
	if !success {
		t.Errorf("[+] %s: FAIL", script)
		t.Fatal(fmt.Sprintf("%s: %s", executionError.TypeName(), executionError.String))
		return
	}
	fmt.Println(output.String())
}

func TestWhile(t *testing.T) {
	test(t, filepath.Join(whileTest))
}

func TestFor(t *testing.T) {
	test(t, filepath.Join(forTest))
}

func TestFibonacci(t *testing.T) {
	test(t, filepath.Join(fibonacci))
}
