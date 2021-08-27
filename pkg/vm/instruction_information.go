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
