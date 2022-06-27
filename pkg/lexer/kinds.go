package lexer

/*
	Token Kinds
*/

type Kind uint8

const (
	Unknown Kind = iota
	PendingEscape
	Comment
	Whitespace
	Literal
	Tab
	IdentifierKind
	JunkKind
	Separator
	Punctuation
	Assignment
	Comparator
	Operator
	Keyboard
	Boolean
	NoneType
	EOF
)
