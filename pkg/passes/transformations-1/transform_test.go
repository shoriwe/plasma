package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
	"github.com/shoriwe/gplasma/pkg/parser"
	"github.com/shoriwe/gplasma/pkg/passes/simplification"
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
		simplified := simplification.Simplify(program)
		Transform(simplified)
	}
}

func TestSampleScript(t *testing.T) {
	test(t, basic.Samples)
}
