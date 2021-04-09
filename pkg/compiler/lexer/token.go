package lexer

/*
	Token Kinds
*/

const (
	Unknown = iota
	PendingEscape
	Comment
	Whitespace
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
	ScientificFloat
	CommandOutput
	ByteString
	Keyboard
	EOF

	Comma
	Colon
	SemiColon
	NewLine
	Pass
	Super
	End
	If
	Else
	Elif
	While
	For
	Until
	Switch
	Case
	Yield
	Return
	Retry
	Break
	Redo
	Module
	Def
	Lambda
	Struct
	Interface
	Go
	Class
	Try
	Except
	Finally
	And
	Or
	Xor
	In
	IsInstanceOf
	Async
	Await
	BEGIN
	END
	Enum
	Not
	Assign
	Equals
	SignNot
	NotEqual
	NegateBits
	NegateBitsAssign
	GreaterThan
	GreaterOrEqualThan
	LessThan
	LessOrEqualThan
	BitwiseOr
	BitwiseXor
	BitWiseAnd
	BitwiseLeft
	BitwiseRight
	Add
	Sub
	Star
	Div
	FloorDiv
	Modulus
	PowerOf

	OpenParentheses
	CloseParentheses
	OpenSquareBracket
	CloseSquareBracket
	OpenBrace
	CloseBrace
	DollarSign
	Dot

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
)

/*
	Regular Expressions
*/

type Token struct {
	String      string
	DirectValue int
	Kind        int
	Line        int
	Column      int
	Index       int
}

/*
	Keyboards
*/

var (
	CommaString              = ","
	ColonString              = ":"
	SemiColonString          = ";"
	NewLineString            = "\n"
	PassString               = "pass"
	SuperString              = "super"
	EndString                = "end"
	IfString                 = "if"
	ElseString               = "else"
	ElifString               = "elif"
	WhileString              = "while"
	ForString                = "for"
	UntilString              = "until"
	SwitchString             = "switch"
	CaseString               = "case"
	YieldString              = "yield"
	ReturnString             = "return"
	RetryString              = "retry"
	BreakString              = "break"
	RedoString               = "redo"
	ModuleString             = "module"
	DefString                = "def"
	LambdaString             = "lambda"
	StructString             = "struct"
	InterfaceString          = "interface"
	GoString                 = "go"
	ClassString              = "class"
	TryString                = "try"
	ExceptString             = "except"
	FinallyString            = "finally"
	AndString                = "and"
	OrString                 = "or"
	XorString                = "xor"
	InString                 = "in"
	IsInstanceOfString       = "isinstanceof"
	AsyncString              = "async"
	AwaitString              = "await"
	BEGINString              = "BEGIN"
	ENDString                = "END"
	EnumString               = "enum"
	NotString                = "not"
	OpenParenthesesString    = "("
	CloseParenthesesString   = ")"
	OpenSquareBracketString  = "["
	CloseSquareBracketString = "]"
	OpenBraceString          = "{"
	CloseBraceString         = "}"
	DollarSignString         = "$"
	DotString                = "."
	BitwiseOrString          = "|"
	BitwiseXorString         = "^"
	BitWiseAndString         = "&"
	BitwiseLeftString        = "<<"
	BitwiseRightString       = ">>"
	AddString                = "+"
	SubString                = "-"
	StarString               = "*"
	DivString                = "/"
	FloorDivString           = "//"
	ModulusString            = "%"
	PowerOfString            = "**"
	LessThanString           = "<"
	GreatThanString          = ">"
	NegateBitsString         = "~"
	SignNotString            = "!"
	EqualsString             = "="
	WhiteSpaceString         = " "
	TabString                = "\t"
	CommentString            = "#"
)
