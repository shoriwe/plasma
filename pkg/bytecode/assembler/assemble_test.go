package assembler

import (
	"fmt"
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
		l := lexer.NewLexer(reader.NewStringReader(sample))
		p := parser.NewParser(l)
		program, parseError := p.Parse()
		if parseError != nil {
			t.Fatalf("Failed in %s with error %s", script, parseError.Error())
		}
		simplifiedProgram := simplification.Simplify(program)
		transformedProgram := transformations_1.Transform(simplifiedProgram)
		bytecode := Assemble(transformedProgram)
		fmt.Println(len(bytecode))
	}
}

func TestSampleScript(t *testing.T) {
	test(t, basic.Samples)
}
