package vm

type ClassInformation struct {
	Name       string
	BodyLength int
}

type FunctionInformation struct {
	Name              string
	BodyLength        int
	NumberOfArguments int
}

type ConditionInformation struct {
	BodyLength     int
	ElseBodyLength int
}

type LoopInformation struct {
	BodyLength      int
	ConditionLength int
	Receivers       []string
}
