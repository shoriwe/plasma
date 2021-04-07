package lexer

import (
	"fmt"
	"testing"
)

func test(t *testing.T, samples []string) {
	for _, sample := range samples {
		lexer := NewLexer(sample)
		for ; lexer.HasNext(); {
			token, tokenizationError := lexer.Next()
			if tokenizationError != nil {
				t.Error(tokenizationError)
				continue
			}
			fmt.Println(token)
		}
	}
}

var stringSamples = []string{
	"\"this is a string expression\n\"",
	"\"concat#{foobar}\"",
	"'concat#{foobar}'",
	"b\"Hello World\"",
}

func TestString(t *testing.T) {
	test(t, stringSamples)
}

var commandOutputSamples = []string{
	"`date`",
	"`whoami\ndate`",
}

func TestCommandOutput(t *testing.T) {
	test(t, commandOutputSamples)
}

var numericSamples = []string{
	"1\n",                   // Integer
	"0.1234\n",              // Float
	"1e-234\n",              // Scientific Number
	"0x34587_345798923",     // Hexadecimal Number
	"0b0010100100101_1001",  //  Binary Number
	"0x0004357_435345 << 1", // Octal Number
}

func TestNumeric(t *testing.T) {
	test(t, numericSamples)
}

var complexSamples = []string{
	"for a in range(1, 2)\n1-\n2*\n4\npass\nend\n",
}

func TestComplex(t *testing.T) {
	test(t, complexSamples)
}
