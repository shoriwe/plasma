package lexer

import "regexp"

/*
	Token Kinds
*/

const (
	Unknown = iota
	Separator
	String1
	String2
	CommandOutput1
	CommandOutput2
	RegularExpression1
	RegularExpression2
	Keyboard
	EOF
)

/*
	Regular Expressions
*/

var (
	singleQuoteStringPattern = regexp.MustCompile("'.*'")
	doubleQuoteStringPattern = regexp.MustCompile("\".*\"")
	commandOutputPattern     = regexp.MustCompile("`.*`")
	integerPattern           = regexp.MustCompile("0|([1-9]+[_0-9]*)")
	hexadecimalPattern       = regexp.MustCompile("0[xX][0-9a-fA-F]+[_0-9a-fA-F]_")
	binaryPattern            = regexp.MustCompile("0[bB][01]+[_01]*")
	octalPattern             = regexp.MustCompile("0[0-7]+[_0-7]*")
	isPunctuationPattern     = regexp.MustCompile("[!\"#$%&'()*+,-./:;<=>?@\\[\\]\\\\^_`{|}~ ]")
)

type Token struct {
	String string
	Kind   rune
	Line   int
	Column int
	Index  [2]int
}

/*
	Keyboards
*/

var (
	BEGIN   = "BEGIN"
	Class   = "class"
	Ensure  = "ensure"
	Nil     = "nil"
	Self    = "self"
	When    = "when"
	END     = "END"
	Def     = "def"
	False   = "false"
	Not     = "not"
	Super   = "super"
	While   = "while"
	Alias   = "alias"
	Defined = "defined"
	For     = "for"
	Or      = "or"
	Then    = "then"
	Yield   = "yield"
	And     = "and"
	Do      = "do"
	If      = "if"
	Redo    = "redo"
	True    = "true"
	Begin   = "begin"
	Else    = "else"
	In      = "in"
	Rescue  = "rescue"
	Undex   = "undef"
	Break   = "break"
	Elsif   = "elsif"
	Module  = "module"
	Retry   = "retry"
	Unless  = "unless"
	Case    = "case"
	End     = "end"
	Next    = "next"
	Return  = "return"
	until   = "until"
)
