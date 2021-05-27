package vm

const (
	// Literal initialization
	NewStringOP uint8 = iota
	NewBytesOP
	NewIntegerOP
	NewFloatOP
	NewTrueBoolOP
	NewFalseBoolOP
	GetNoneOP

	NoOP
	PushOP
	PushN_OP
	CopyOP
	CallOP
	GetOP
	GetFromOP
	ReturnOP
)
