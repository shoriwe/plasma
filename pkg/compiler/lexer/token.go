package lexer

/*
	Token Kinds
*/

const (
	Unknown = iota
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
	ScientificFloat
	CommandOutput
	ByteString
	Keyboard
	Boolean
	NoneType
	EOF

	// Punctuation
	Comma
	Colon
	SemiColon
	NewLine
	// Reserved Keyboards
	Pass
	Super
	End
	If
	Unless
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
	Defer
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
	IsInstanceOf // This maybe can be a regular identifier
	Async
	Await
	BEGIN
	END
	Enum
	GoTo
	Context

	// Assigns
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
	// Unary Operators
	Not
	SignNot
	NegateBits
	// Binary Operators
	//// Logical operators
	And
	Or
	Xor
	In
	////
	Equals
	NotEqual
	GreaterThan
	GreaterOrEqualThan
	LessThan
	LessOrEqualThan
	//// Bitwise Operations
	BitwiseOr
	BitwiseXor
	BitWiseAnd
	BitwiseLeft
	BitwiseRight
	//// Basic Binary expressions
	Add // This is also an unary operator
	Sub // This is also an unary operator
	Star
	Div
	FloorDiv
	Modulus
	PowerOf
	// Reserved Literals
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
	UnlessString             = "unless"
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
	DeferString              = "defer"
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
	TrueString               = "True"
	FalseString              = "False"
	NoneString               = "None"
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
	GoToString               = "goto"
	ContextString            = "context"
)
