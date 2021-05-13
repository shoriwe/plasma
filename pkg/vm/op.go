package vm

const (
	NewStringOP uint8 = iota
	NoOP
	PushOP
	PushN_OP
	CopyOP
	CallOP
	GetOP
	GetFromOP
	ReturnOP
)
