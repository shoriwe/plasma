package lexer

type DirectValue uint8

const (
	SingleQuoteString DirectValue = iota
	DoubleQuoteString
	ByteString
	Integer
	HexadecimalInteger
	BinaryInteger
	OctalInteger
	Float
	ScientificFloat
	CommandOutput

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
	Continue
	Break
	Redo
	Module
	Def
	Lambda
	Interface
	Class
	Try
	Except
	Finally
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

	InvalidDirectValue
	Blank
)

/*
	Regular Expressions
*/

type Token struct {
	Contents    []rune
	DirectValue DirectValue
	Kind        Kind
	Line        int
	Column      int
	Index       int
}

func (token *Token) String() string {
	return string(token.Contents)
}

func (token *Token) append(r ...rune) {
	token.Contents = append(token.Contents, r...)
}

/*
	Keyboards
*/

var (
	CommaChar              = ','
	ColonChar              = ':'
	SemiColonChar          = ';'
	NewLineChar            = '\n'
	OpenParenthesesChar    = '('
	CloseParenthesesChar   = ')'
	OpenSquareBracketChar  = '['
	CloseSquareBracketChar = ']'
	OpenBraceChar          = '{'
	CloseBraceChar         = '}'
	DollarSignChar         = '$'
	DotChar                = '.'
	BitwiseOrChar          = '|'
	BitwiseXorChar         = '^'
	BitWiseAndChar         = '&'
	AddChar                = '+'
	SubChar                = '-'
	StarChar               = '*'
	DivChar                = '/'
	ModulusChar            = '%'
	LessThanChar           = '<'
	GreatThanChar          = '>'
	NegateBitsChar         = '~'
	SignNotChar            = '!'
	EqualsChar             = '='
	WhiteSpaceChar         = ' '
	TabChar                = '\t'
	CommentChar            = '#'
	BackSlashChar          = '\\'

	PassString      = "pass"
	SuperString     = "super"
	EndString       = "end"
	IfString        = "if"
	UnlessString    = "unless"
	ElseString      = "else"
	ElifString      = "elif"
	WhileString     = "while"
	DoString        = "do"
	ForString       = "for"
	UntilString     = "until"
	SwitchString    = "switch"
	CaseString      = "case"
	DefaultString   = "default"
	YieldString     = "yield"
	ReturnString    = "return"
	ContinueString  = "continue"
	BreakString     = "break"
	RedoString      = "redo"
	ModuleString    = "module"
	DefString       = "def"
	LambdaString    = "lambda"
	InterfaceString = "interface"
	ClassString     = "class"
	TryString       = "try"
	ExceptString    = "except"
	FinallyString   = "finally"
	AndString       = "and"
	OrString        = "or"
	XorString       = "xor"
	InString        = "in"
	BEGINString     = "BEGIN"
	ENDString       = "END"
	NotString       = "not"
	TrueString      = "True"
	FalseString     = "False"
	NoneString      = "None"
	ContextString   = "context"
	RaiseString     = "raise"
	AsString        = "as"
)
