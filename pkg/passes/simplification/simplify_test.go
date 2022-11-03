package simplification

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
	"github.com/shoriwe/gplasma/pkg/parser"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/test-samples/basic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func test(t *testing.T, samples map[string]string) {
	for script, sample := range samples {
		l := lexer.NewLexer(reader.NewStringReader(sample))
		p := parser.NewParser(l)
		program, parseError := p.Parse()
		assert.Nil(t, parseError, script)
		_, err := Simplify(program)
		assert.Nil(t, err, script)
	}
}

func TestSampleScript(t *testing.T) {
	test(t, basic.Samples)
}
