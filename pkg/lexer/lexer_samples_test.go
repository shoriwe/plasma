package lexer

var stringSamples = map[string][]*Token{
	"\"this is a string expression\n\"": {
		{
			Contents:    []rune("\"this is a string expression\n\""),
			DirectValue: DoubleQuoteString,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"\"concat#{foobar}\"": {
		{
			Contents:    []rune("\"concat#{foobar}\""),
			DirectValue: DoubleQuoteString,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"'concat#{foobar}'": {
		{
			Contents:    []rune("'concat#{foobar}'"),
			DirectValue: SingleQuoteString,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"b\"Hello World\"": {
		{
			Contents:    []rune("b\"Hello World\""),
			DirectValue: ByteString,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
}

var commandOutputSamples = map[string][]*Token{
	"`date`": {
		{
			Contents:    []rune("`date`"),
			DirectValue: CommandOutput,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"`whoami\ndate`": {
		{
			Contents:    []rune("`whoami\ndate`"),
			DirectValue: CommandOutput,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
}

var numericSamples = map[string][]*Token{
	"1\n": {
		{
			Contents:    []rune("1"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("\n"),
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, // Integer
	"0.1234\n": {
		{
			Contents:    []rune("0.1234"),
			DirectValue: Float,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("\n"),
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, // Float
	"1e-234\n": {
		{
			Contents:    []rune("1e-234"),
			DirectValue: ScientificFloat,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("\n"),
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, // Scientific Number
	"0x34587_345798923": {
		{
			Contents:    []rune("0x34587_345798923"),
			DirectValue: HexadecimalInteger,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, // Hexadecimal Number
	"0b0010100100101_1001": {
		{
			Contents:    []rune("0b0010100100101_1001"),
			DirectValue: BinaryInteger,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	}, //  Binary Number
	"0o0004357_435345 << 1": {
		{
			Contents:    []rune("0o0004357_435345"),
			DirectValue: OctalInteger,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("<<"),
			DirectValue: BitwiseLeft,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("1"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"1 // 2": {
		{
			Contents:    []rune("1"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("//"),
			DirectValue: FloorDiv,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("2"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
	"1**2": {
		{
			Contents:    []rune("1"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("**"),
			DirectValue: PowerOf,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("2"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
}

var complexSamples = map[string][]*Token{
	"for a in range(1, 2)\n1-\n2*\n4\npass\nend\n": {
		{
			Contents:    []rune("for"),
			DirectValue: For,
			Kind:        Keyword,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("a"),
			DirectValue: InvalidDirectValue,
			Kind:        IdentifierKind,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("in"),
			DirectValue: In,
			Kind:        Comparator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("range"),
			DirectValue: InvalidDirectValue,
			Kind:        IdentifierKind,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("("),
			DirectValue: OpenParentheses,
			Kind:        Punctuation,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("1"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune(","),
			DirectValue: Comma,
			Kind:        Punctuation,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("2"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune(")"),
			DirectValue: CloseParentheses,
			Kind:        Punctuation,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("\n"),
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("1"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("-"),
			DirectValue: Sub,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("2"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("*"),
			DirectValue: Star,
			Kind:        Operator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("4"),
			DirectValue: Integer,
			Kind:        Literal,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("\n"),
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("pass"),
			DirectValue: Pass,
			Kind:        Keyword,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("\n"),
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("end"),
			DirectValue: End,
			Kind:        Keyword,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    []rune("\n"),
			DirectValue: NewLine,
			Kind:        Separator,
			Line:        0,
			Column:      0,
			Index:       0,
		},
		{
			Contents:    nil,
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        0,
			Column:      0,
			Index:       0,
		},
	},
}