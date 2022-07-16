package vm

import (
	"bufio"
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
	"github.com/shoriwe/gplasma/pkg/test-samples/success"
	"testing"
)

func TestSampleScripts(t *testing.T) {
	for i := 1; i <= len(success.Samples); i++ {
		sampleScript := fmt.Sprintf("sample-%d.pm", i)
		sampleCode := success.Samples[sampleScript]
		t.Logf("Testing - %s", sampleScript)
		l := lexer.NewLexer(reader.NewStringReader(sampleCode))
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
		writer := bufio.NewWriter(out)
		v := NewVM(nil, writer, writer)
		_, err, _ := v.Execute(bytecode)
		if e := <-err; e != nil {
			t.Fatal(e)
		}
	}
}
