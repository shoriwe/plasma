package opcodes

const (
	Push byte = iota
	Pop
	IdentifierAssign
	SelectorAssign
	IndexAssign
	Jump
	IfJump
	Return
	Require
	DeleteIdentifier
	DeleteSelector
	Defer
	EnterBlock
	ExitBlock
	NewFunction
	NewClass
	Call
	IfOneLiner
	NewArray
	NewTuple
	NewHash
	Identifier
	Integer
	Float
	String
	Bytes
	True
	False
	None
	Selector
	Index
	Super
)
