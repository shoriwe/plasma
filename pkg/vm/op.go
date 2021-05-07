package vm

const (
	NewStringOP uint16 = iota
	PushOP
	PushN_OP
	CopyOP
	CallOP
	GetOP
	GetFromOP
	ReturnOP
)
