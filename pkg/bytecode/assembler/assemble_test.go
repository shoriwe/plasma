package assembler

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
	"github.com/shoriwe/gplasma/pkg/parser"
	"github.com/shoriwe/gplasma/pkg/passes/simplification"
	transformations_1 "github.com/shoriwe/gplasma/pkg/passes/transformations-1"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/test-samples/basic"
	"testing"
)

func test(t *testing.T, samples map[string]string) {
	for script, sample := range samples {
		t.Logf("Testing: %s", script)
		l := lexer.NewLexer(reader.NewStringReader(sample))
		p := parser.NewParser(l)
		program, parseError := p.Parse()
		if parseError != nil {
			t.Fatalf("Failed in %s with error %s", script, parseError.Error())
		}
		simplified, simplificationError := simplification.Simplify(program)
		if simplificationError != nil {
			t.Fatal(simplificationError)
		}
		transformed, transformError := transformations_1.Transform(simplified)
		if transformError != nil {
			t.Fatal(transformError)
		}
		bytecode, assembleError := Assemble(transformed)
		if assembleError != nil {
			t.Fatal(assembleError)
		}
		bytecodeSize := float64(len(bytecode)) / 1024
		t.Logf("Size of %s: %db => %fkb", script, len(bytecode), bytecodeSize)
	}
}

func TestSampleScript(t *testing.T) {
	test(t, basic.Samples)
}
