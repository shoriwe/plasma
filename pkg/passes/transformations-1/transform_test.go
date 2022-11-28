package transformations_1

import (
	"github.com/shoriwe/plasma/pkg/lexer"
	"github.com/shoriwe/plasma/pkg/parser"
	"github.com/shoriwe/plasma/pkg/passes/simplification"
	"github.com/shoriwe/plasma/pkg/reader"
	"github.com/shoriwe/plasma/pkg/test-samples/basic"
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
		_, transformError := Transform(simplified)
		assert.Nil(t, transformError)
	}
}

func TestSampleScript(t *testing.T) {
	test(t, basic.Samples)
}
