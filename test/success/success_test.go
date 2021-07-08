package success

import (
	"bytes"
	"fmt"
	"github.com/shoriwe/gplasma/pkg/compiler/plasma"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
	Bunch of samples that should only compile and execute (no result checking)
*/
const (
	expressionSamples   = "expressions"
	literals            = "literals"
	composites          = "composites"
	unaryExpressions    = "unary-expressions"
	binaryExpressions   = "binary-expressions"
	indexExpressions    = "index-expressions"
	selectorExpressions = "selector-expressions"
	lambdaExpressions   = "lambdas"
	ifUnlessOneLiners   = "if_unless-one-liner"
	generatorExpression = "generators"

	statementSamples   = "statements"
	assignStatement    = "assignment"
	functionDefinition = "function-definition"
	ifStatement        = "if-statement"
	doWhileStatement   = "do_while-statement"
	beginEnd           = "begin-end"
	whileStatement     = "while-statement"
	untilStatement     = "until-statement"
	forStatement       = "for-statement"
	tryStatement       = "try-blocks"
	moduleStatement    = "module-statement"
	classStatement     = "class-statement"
	interfaceStatement = "interface-statement"
	switchStatement    = "switch-statement"
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
			t.Errorf("[+] %s: FAIL", file.Name())
			t.Logf("Output:\n%s", output.String())
			t.Fatal(fmt.Sprintf("%s: %s", executionError.TypeName(), executionError.GetString()))
			return
		}
		if strings.Contains(output.String(), "False\n") {
			t.Errorf("[+] %s: FAIL", file.Name())
			t.Logf("Output:\n%s", output.String())
			return
		}
		fmt.Println(fmt.Sprintf("[+] %s: SUCCESS", file.Name()))
	}
}

func TestLiterals(t *testing.T) {
	test(t, filepath.Join("success", expressionSamples, literals))
}

func TestComposites(t *testing.T) {
	test(t, filepath.Join("success", expressionSamples, composites))
}

func TestUnaryExpressions(t *testing.T) {
	test(t, filepath.Join("success", expressionSamples, unaryExpressions))
}

func TestBinaryExpressions(t *testing.T) {
	test(t, filepath.Join("success", expressionSamples, binaryExpressions))
}

func TestIndexExpressions(t *testing.T) {
	test(t, filepath.Join("success", expressionSamples, indexExpressions))
}

func TestSelectorExpressions(t *testing.T) {
	test(t, filepath.Join("success", expressionSamples, selectorExpressions))
}

func TestLambdaExpressions(t *testing.T) {
	test(t, filepath.Join("success", expressionSamples, lambdaExpressions))
}

func TestIfAndUnlessOneLinersExpressions(t *testing.T) {
	test(t, filepath.Join("success", expressionSamples, ifUnlessOneLiners))
}

func TestGeneratorExpressions(t *testing.T) {
	test(t, filepath.Join("success", expressionSamples, generatorExpression))
}

// Statement tests

func TestAssignStatement(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, assignStatement))
}

func TestIfStatement(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, ifStatement))
}

func TestForStatements(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, forStatement))
}

func TestWhileStatements(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, whileStatement))
}

func TestUntilStatements(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, untilStatement))
}

func TestDoWhileStatement(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, doWhileStatement))
}

func TestFunctionDefinitionStatement(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, functionDefinition))
}

func TestBeginEndStatements(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, beginEnd))
}

func TestTryStatements(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, tryStatement))
}

func TestModuleStatements(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, moduleStatement))
}

func TestClassStatements(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, classStatement))
}

func TestInterfaceStatements(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, interfaceStatement))
}

func TestSwitchStatements(t *testing.T) {
	test(t, filepath.Join("success", statementSamples, switchStatement))
}
