package lexer

import (
	"testing"
)

func test(t *testing.T, samples map[string][]*Token) {
	for sample, result := range samples {
		lexer := NewLexer(sample)
		var computedTokens []*Token
		for ; lexer.HasNext(); {
			token, tokenizationError := lexer.Next()
			if tokenizationError != nil {
				t.Error(tokenizationError)
				continue
			}
			computedTokens = append(computedTokens, token)
		}
		if len(computedTokens) != len(result) {
			t.Errorf("Not equal number of token in sample: %s; expecting %d recevied %d", sample, len(result), len(computedTokens))
			return
		}
		for index, computedToken := range computedTokens {
			if computedToken.String != result[index].String {
				t.Errorf("Expecting Token: '%s' but Received Token: %s in sample %s", result[index].String, computedToken.String, sample)
				return
			}
			if computedToken.Kind != result[index].Kind {
				t.Errorf("Expecting Kind: %d, but Received Kind: %d in sample %s", result[index].Kind, computedToken.Kind, sample)
				return
			}
			if computedToken.DirectValue != result[index].DirectValue {
				t.Errorf("Expecting DirectValue: %d, but Received DirectValue: %d in sample %s", result[index].DirectValue, computedToken.DirectValue, sample)
				return
			}
		}
	}
}

var stringSamples = map[string][]*Token{
	"\"this is a string expression\n\"": {
		{
			String:      "\"this is a string expression\n\"",
			DirectValue: DoubleQuoteString,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"\"concat#{foobar}\"": {
		{
			String:      "\"concat#{foobar}\"",
			DirectValue: DoubleQuoteString,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"'concat#{foobar}'": {
		{
			String:      "'concat#{foobar}'",
			DirectValue: SingleQuoteString,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"b\"Hello World\"": {
		{
			String:      "b\"Hello World\"",
			DirectValue: ByteString,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
}

func TestString(t *testing.T) {
	test(t, stringSamples)
}

var commandOutputSamples = map[string][]*Token{
	"`date`": {
		{
			String:      "`date`",
			DirectValue: CommandOutput,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"`whoami\ndate`": {
		{
			String:      "`whoami\ndate`",
			DirectValue: CommandOutput,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
}

func TestCommandOutput(t *testing.T) {
	test(t, commandOutputSamples)
}

var numericSamples = map[string][]*Token{
	"1\n": {
		{
			String:      "1",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "\n",
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, // Integer
	"0.1234\n": {
		{
			String:      "0.1234",
			DirectValue: Float,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "\n",
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, // Float
	"1e-234\n": {
		{
			String:      "1e-234",
			DirectValue: ScientificFloat,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "\n",
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, // Scientific Number
	"0x34587_345798923": {
		{
			String:      "0x34587_345798923",
			DirectValue: HexadecimalInteger,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, // Hexadecimal Number
	"0b0010100100101_1001": {
		{
			String:      "0b0010100100101_1001",
			DirectValue: BinaryInteger,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, //  Binary Number
	"0o0004357_435345 << 1": {
		{
			String:      "0o0004357_435345",
			DirectValue: OctalInteger,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "<<",
			DirectValue: BitwiseLeft,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "1",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"1 // 2": {
		{
			String:      "1",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "//",
			DirectValue: FloorDiv,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "2",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"1**2": {
		{
			String:      "1",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "**",
			DirectValue: PowerOf,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "2",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
}

func TestNumeric(t *testing.T) {
	test(t, numericSamples)
}

var complexSamples = map[string][]*Token{
	"for a in range(1, 2)\n1-\n2*\n4\npass\nend\n": {
		{
			String:      "for",
			DirectValue: For,
			Kind:        Keyboard,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "a",
			DirectValue: -1,
			Kind:        IdentifierKind,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "in",
			DirectValue: In,
			Kind:        Keyboard,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "range",
			DirectValue: -1,
			Kind:        IdentifierKind,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "(",
			DirectValue: OpenParentheses,
			Kind:        Punctuation,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "1",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      ",",
			DirectValue: Comma,
			Kind:        Punctuation,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "2",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      ")",
			DirectValue: CloseParentheses,
			Kind:        Punctuation,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "\n",
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "1",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "-",
			DirectValue: Sub,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "2",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "*",
			DirectValue: Star,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "4",
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "\n",
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "pass",
			DirectValue: Pass,
			Kind:        Keyboard,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "\n",
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "end",
			DirectValue: End,
			Kind:        Keyboard,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "\n",
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
}

func TestComplex(t *testing.T) {
	test(t, complexSamples)
}
