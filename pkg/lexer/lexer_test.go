package lexer

import (
	reader2 "github.com/shoriwe/gplasma/pkg/reader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func test(t *testing.T, samples map[string][]*Token) {
	for sample, result := range samples {
		lexer := NewLexer(reader2.NewStringReader(sample))
		var computedTokens []*Token
		for lexer.HasNext() {
			token, tokenizationError := lexer.Next()
			assert.Nil(t, tokenizationError)
			computedTokens = append(computedTokens, token)
		}
		assert.Equal(t, len(result), len(computedTokens))
		for index, computedToken := range computedTokens {
			assert.Equal(t, result[index].String(), computedToken.String())
			assert.Equal(t, result[index].Kind, computedToken.Kind)
			assert.Equal(t, result[index].DirectValue, computedToken.DirectValue)
		}
	}
}

func TestString(t *testing.T) {
	test(t, stringSamples)
}

func TestCommandOutput(t *testing.T) {
	test(t, commandOutputSamples)
}

func TestNumeric(t *testing.T) {
	test(t, numericSamples)
}

func TestComplex(t *testing.T) {
	test(t, complexSamples)
}
