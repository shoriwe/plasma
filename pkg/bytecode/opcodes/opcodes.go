package opcodes

const (
	Push byte = iota
	Pop
	IdentifierAssign
	SelectorAssign
	Label
	Jump
	IfJump
	Return
	DeleteIdentifier
	DeleteSelector
	Defer
	EnterBlock
	ExitBlock
	NewFunction
	NewClass
	Call
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
	Super
)
