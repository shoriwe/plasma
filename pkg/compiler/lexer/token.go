package lexer

/*
	Token Kinds
*/

const (
	Unknown = iota
	Whitespace
	ConstantKind
	IdentifierKind
	Separator
	Punctuation
	Assignment
	Comparator
	Operator
	SingleQuoteString
	DoubleQuoteString
	Integer
	HexadecimalInteger
	BinaryInteger
	OctalInteger
	Float
	ScientificFloat
	CommandOutput
	ByteString
	Keyboard
	EOF
)

/*
	Regular Expressions
*/

type Token struct {
	String string
	Kind   rune
	Line   int
	Column int
	Index  int
}

/*
	Keyboards
*/

var (
	Pass         = "pass"
	Super        = "super"
	End          = "end"
	If           = "if"
	Else         = "else"
	Elif         = "elif"
	While        = "while"
	For          = "for"
	Until        = "until"
	Switch       = "switch"
	Case         = "case"
	Yield        = "yield"
	Return       = "return"
	Retry        = "retry"
	Break        = "break"
	Redo         = "redo"
	Module       = "module"
	Def          = "def"
	Lambda       = "lambda"
	Struct       = "struct"
	Interface    = "interface"
	Go           = "go"
	Class        = "class"
	Try          = "try"
	Except       = "except"
	Finally      = "finally"
	And          = "and"
	Or           = "or"
	Xor          = "xor"
	In           = "in"
	IsInstanceOf = "isinstanceof"
	Async        = "async"
	Await        = "await"
	BEGIN        = "BEGIN"
	END          = "END"
	Enum         = "enum"
	Not          = "not"
)
