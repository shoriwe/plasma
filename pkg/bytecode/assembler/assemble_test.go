package assembler

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
	"github.com/shoriwe/gplasma/pkg/parser"
	"github.com/shoriwe/gplasma/pkg/passes/simplification"
	transformations_1 "github.com/shoriwe/gplasma/pkg/passes/transformations-1"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/test-samples/basic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func test(t *testing.T, samples map[string]string) {
	for _, sample := range samples {
		l := lexer.NewLexer(reader.NewStringReader(sample))
		p := parser.NewParser(l)
		program, parseError := p.Parse()
		assert.Nil(t, parseError)
		simplified, simplificationError := simplification.Simplify(program)
		assert.Nil(t, simplificationError)
		transformed, transformError := transformations_1.Transform(simplified)
		assert.Nil(t, transformError)
		_, assembleError := Assemble(transformed)
		assert.Nil(t, assembleError)
	}
}

func TestSampleScript(t *testing.T) {
	test(t, basic.Samples)
}
