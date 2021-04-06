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

var string1Samples = []string{
	"\"this is a string expression\n\"",
	"\"concat#{foobar}\"",
	"'concat#{foobar}'",
}

func TestString1(t *testing.T) {
	test(t, string1Samples)
}

var string2Samples = []string{
	"%q!I said, \"You said, 'She said it.'\"!",
	"%!I said, \"You said, 'She said it.'\"!",
	"%Q('This is it.'\\n)",
}

func TestString2(t *testing.T) {
	test(t, string2Samples)
}

var commandOutput1Samples = []string{
	"`date`",
	"`whoami\ndate`",
}

func TestCommandOutput1(t *testing.T) {
	test(t, commandOutput1Samples)
}

var commandOutput2Samples = []string{
	"%x{Hello}",
	"%x\nhello\n",
}

func TestCommandOutput2(t *testing.T) {
	test(t, commandOutput2Samples)
}