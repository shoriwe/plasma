package parser

import (
	"fmt"
	lexer2 "github.com/shoriwe/gruby/pkg/compiler/lexer"
	"testing"
)

func test(t *testing.T, samples []string) {
	for _, sample := range samples {
		lexer := lexer2.NewLexer(sample)
		parser := NewParser(lexer)
		program, parsingError := parser.Parse()
		if parsingError != nil {
			t.Error(parsingError)
			return
		}
		fmt.Println(program)
	}
}

var basicSamples = []string{
	"hello = 1 + 2",
}

func TestParseBasic(t *testing.T) {
	test(t, basicSamples)
}
