package gruby

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/compiler/plasma"
	"github.com/shoriwe/gruby/pkg/reader"
	"github.com/shoriwe/gruby/pkg/vm"
	"os"
	"path/filepath"
	"testing"
)

const (
	package_          = "pkg"
	samplesDirectory  = "tests_samples"
	literals          = "literals"
	composites        = "composites"
	unaryExpressions  = "unary-expressions"
	binaryExpressions = "binary-expressions"
)

func test(t *testing.T, directory string) {
	currentDir, currentDirGetError := os.Getwd()
	if currentDirGetError != nil {
		t.Fatal(currentDirGetError)
		return
	}
	directory = filepath.Join(currentDir, directory)
	directoryContent, directoryReadingError := os.ReadDir(directory)
	if directoryReadingError != nil {
		t.Fatal(directoryReadingError)
		return
	}
	for _, file := range directoryContent {
		if file.IsDir() {
			continue
		}
		fileHandler, openError := os.Open(filepath.Join(directory, file.Name()))
		if openError != nil {
			t.Fatal(openError)
			return
		}
		compiler := plasma.NewCompiler(reader.NewStringReaderFromFile(fileHandler))
		// content, _ := io.ReadAll(fileHandler)
		// compiler := plasma.NewCompiler(reader.NewStringReader(string(content)))
		code, compilingError := compiler.Compile()
		if compilingError != nil {
			t.Fatal(compilingError)
			return
		}

		plasmaVm := vm.NewPlasmaVM()
		plasmaVm.InitializeByteCode(code)
		result, executionError := plasmaVm.Execute()
		// _, executionError := plasmaVm.Execute()
		if executionError != nil {
			t.Fatal(executionError)
			return
		}
		resultToString, getError := result.Get(vm.ToString)
		if getError != nil {
			t.Fatal(getError)
			return
		}
		if _, ok := resultToString.(*vm.Function); !ok {
			t.Fatal("Expecting ToString function")
			return
		}
		stringResult, callError := vm.CallFunction(resultToString.(*vm.Function), plasmaVm, result.SymbolTable())
		if callError != nil {
			t.Fatal(callError)
			return
		}
		fmt.Println(stringResult.GetString())
		fmt.Println(fmt.Sprintf("[+] %s: SUCCESS", file.Name()))
	}
}

func TestLiterals(t *testing.T) {
	test(t, filepath.Join(package_, samplesDirectory, literals))
}

func TestComposites(t *testing.T) {
	test(t, filepath.Join(package_, samplesDirectory, composites))
}

func TestUnaryExpressions(t *testing.T) {
	test(t, filepath.Join(package_, samplesDirectory, unaryExpressions))
}

func TestBinaryExpressions(t *testing.T) {
	test(t, filepath.Join(package_, samplesDirectory, binaryExpressions))
}
