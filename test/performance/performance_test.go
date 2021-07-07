package performance

import (
	"bytes"
	"fmt"
	"github.com/shoriwe/gplasma/pkg/compiler/plasma"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
	"path/filepath"
	"testing"
)

const (
	performance = "performance"
	fibonacci   = "fibo.pm"
	whileTest   = "while.pm"
	forTest     = "for.pm"
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
	compiler := plasma.NewCompiler(reader.NewStringReaderFromFile(fileHandler),
		plasma.Options{
			Debug: false,
		},
	)
	code, compilingError := compiler.Compile()
	if compilingError != nil {
		t.Fatal(compilingError)
		return
	}
	output := bytes.NewBuffer(make([]byte, 0))
	plasmaVm := vm.NewPlasmaVM(nil, output, output)
	_, executionError := plasmaVm.Execute(nil, code)
	if executionError != nil {
		t.Errorf("[+] %s: FAIL", script)
		t.Fatal(fmt.Sprintf("%s: %s", executionError.TypeName(), executionError.GetString()))
		return
	}

}

// At this moment, all this tests are extremely slow to consider them successful

/*

func TestFibonacci(t *testing.T) {
	test(t, filepath.Join("test", performance, fibonacci))
}
*/

func TestWhile(t *testing.T) {
	test(t, filepath.Join("test", performance, whileTest))
}

func TestFor(t *testing.T) {
	test(t, filepath.Join("test", performance, forTest))
}
