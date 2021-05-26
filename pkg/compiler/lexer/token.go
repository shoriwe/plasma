package lexer

/*
	Token Kinds
*/

const (
	Unknown uint8 = iota
	PendingEscape
	Comment
	Whitespace
	Literal
	Tab
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
	AwaitKeyboard
	ScientificFloat
	CommandOutput
	ByteString
	Keyboard
	Boolean
	NoneType
	EOF

	Comma
	Colon
	SemiColon
	NewLine

	Pass
	Super
	End
	If
	Unless
	As
	Raise
	Else
	Elif
	While
	Do
	For
	Until
	Switch
	Case
	Default
	Yield
	Return
	Retry
	Break
	Redo
	Defer
	Module
	Def
	Lambda
	Interface
	Class
	Try
	Except
	Finally
	IsInstanceOf // This maybe can be a regular identifier
	Async
	Await
	BEGIN
	END
	Context

	Assign
	NegateBitsAssign
	BitwiseOrAssign
	BitwiseXorAssign
	BitWiseAndAssign
	BitwiseLeftAssign
	BitwiseRightAssign
	AddAssign
	SubAssign
	StarAssign
	DivAssign
	FloorDivAssign
	ModulusAssign
	PowerOfAssign

	Not
	SignNot
	NegateBits

	And
	Or
	Xor
	In

	Equals
	NotEqual
	GreaterThan
	GreaterOrEqualThan
	LessThan
	LessOrEqualThan

	BitwiseOr
	BitwiseXor
	BitWiseAnd
	BitwiseLeft
	BitwiseRight

	Add // This is also an unary operator
	Sub // This is also an unary operator
	Star
	Div
	FloorDiv
	Modulus
	PowerOf

	True
	False
	None

	OpenParentheses
	CloseParentheses
	OpenSquareBracket
	CloseSquareBracket
	OpenBrace
	CloseBrace
	DollarSign
	Dot
)

/*
	Regular Expressions
*/

type Token struct {
	String      string
	DirectValue uint8
	Kind        uint8
	Line        int
	Column      int
	Index       int
}

/*
	Keyboards
*/

var (
	CommaChar              uint8 = ','
	ColonChar              uint8 = ':'
	SemiColonChar          uint8 = ';'
	NewLineChar            uint8 = '\n'
	OpenParenthesesChar    uint8 = '('
	CloseParenthesesChar   uint8 = ')'
	OpenSquareBracketChar  uint8 = '['
	CloseSquareBracketChar uint8 = ']'
	OpenBraceChar          uint8 = '{'
	CloseBraceChar         uint8 = '}'
	DollarSignChar         uint8 = '$'
	DotChar                uint8 = '.'
	BitwiseOrChar          uint8 = '|'
	BitwiseXorChar         uint8 = '^'
	BitWiseAndChar         uint8 = '&'
	AddChar                uint8 = '+'
	SubChar                uint8 = '-'
	StarChar               uint8 = '*'
	DivChar                uint8 = '/'
	ModulusChar            uint8 = '%'
	LessThanChar           uint8 = '<'
	GreatThanChar          uint8 = '>'
	NegateBitsChar         uint8 = '~'
	SignNotChar            uint8 = '!'
	EqualsChar             uint8 = '='
	WhiteSpaceChar         uint8 = ' '
	TabChar                uint8 = '\t'
	CommentChar            uint8 = '#'
	BackSlashChar          uint8 = '\\'

	PassString         = "pass"
	SuperString        = "super"
	EndString          = "end"
	IfString           = "if"
	UnlessString       = "unless"
	ElseString         = "else"
	ElifString         = "elif"
	WhileString        = "while"
	DoString           = "do"
	ForString          = "for"
	UntilString        = "until"
	SwitchString       = "switch"
	CaseString         = "case"
	DefaultString      = "default"
	YieldString        = "yield"
	ReturnString       = "return"
	RetryString        = "retry"
	BreakString        = "break"
	RedoString         = "redo"
	DeferString        = "defer"
	ModuleString       = "module"
	DefString          = "def"
	LambdaString       = "lambda"
	InterfaceString    = "interface"
	ClassString        = "class"
	TryString          = "try"
	ExceptString       = "except"
	FinallyString      = "finally"
	AndString          = "and"
	OrString           = "or"
	XorString          = "xor"
	InString           = "in"
	IsInstanceOfString = "isinstanceof"
	AsyncString        = "async"
	AwaitString        = "await"
	BEGINString        = "BEGIN"
	ENDString          = "END"
	NotString          = "not"
	TrueString         = "True"
	FalseString        = "False"
	NoneString         = "None"
	ContextString      = "context"
	RaiseString        = "raise"
	AsString           = "as"
)
