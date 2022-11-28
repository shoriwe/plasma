package simplification

import (
	"github.com/shoriwe/plasma/pkg/lexer"
	"github.com/shoriwe/plasma/pkg/parser"
	"github.com/shoriwe/plasma/pkg/reader"
	"github.com/shoriwe/plasma/pkg/test-samples/basic"
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
