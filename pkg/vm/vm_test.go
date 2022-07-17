package vm

import (
	"bytes"
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/bytecode/assembler"
	"github.com/shoriwe/gplasma/pkg/lexer"
	"github.com/shoriwe/gplasma/pkg/parser"
	"github.com/shoriwe/gplasma/pkg/passes/checks"
	"github.com/shoriwe/gplasma/pkg/passes/simplification"
	transformations_1 "github.com/shoriwe/gplasma/pkg/passes/transformations-1"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/test-samples/fail"
	"github.com/shoriwe/gplasma/pkg/test-samples/success"
	"testing"
)

func TestSuccessSampleScripts(t *testing.T) {
	for i := 1; i <= len(success.Samples); i++ {
		sampleScript := fmt.Sprintf("sample-%d.pm", i)
		script := success.Samples[sampleScript]
		t.Logf("Testing - %s", sampleScript)
		l := lexer.NewLexer(reader.NewStringReader(script.Code))
		p := parser.NewParser(l)
		program, parseError := p.Parse()
		if parseError != nil {
			t.Fatal(parseError)
		}
		checkPass := checks.NewCheckPass()
		ast.Walk(checkPass, program)
		switch {
		case checkPass.CountInvalidLoopNodes() > 0:
			t.Fatal("invalid loop nodes")
		case checkPass.CountInvalidGeneratorNodes() > 0:
			t.Fatal("invalid generator nodes")
		case checkPass.CountInvalidFunctionNodes() > 0:
			t.Fatal("invalid function nodes")
		}
		simplified := simplification.Simplify(program)
		transformed := transformations_1.Transform(simplified)
		bytecode := assembler.Assemble(transformed)
		out := &bytes.Buffer{}
		v := NewVM(nil, out, out)
		_, err, _ := v.Execute(bytecode)
		if e := <-err; e != nil {
			t.Fatal(e)
		}
		s := out.String()
		if s != script.Result {
			t.Log("Expecting:")
			fmt.Println(script.Result)
			t.Log("But obtained:")
			fmt.Println(s)
			t.Fatal("Invalid result")
		}
	}
}

func TestFailSampleScripts(t *testing.T) {
	for i := 1; i <= len(fail.Samples); i++ {
		sampleScript := fmt.Sprintf("sample-%d.pm", i)
		script := fail.Samples[sampleScript]
		t.Logf("Testing - %s", sampleScript)
		l := lexer.NewLexer(reader.NewStringReader(script))
		p := parser.NewParser(l)
		program, parseError := p.Parse()
		if parseError != nil {
			t.Fatal(parseError)
		}
		checkPass := checks.NewCheckPass()
		ast.Walk(checkPass, program)
		switch {
		case checkPass.CountInvalidLoopNodes() > 0:
			t.Fatal("invalid loop nodes")
		case checkPass.CountInvalidGeneratorNodes() > 0:
			t.Fatal("invalid generator nodes")
		case checkPass.CountInvalidFunctionNodes() > 0:
			t.Fatal("invalid function nodes")
		}
		simplified := simplification.Simplify(program)
		transformed := transformations_1.Transform(simplified)
		bytecode := assembler.Assemble(transformed)
		out := &bytes.Buffer{}
		v := NewVM(nil, out, out)
		_, err, _ := v.Execute(bytecode)
		if e := <-err; e == nil {
			t.Fatal("should fail")
		}
	}
}
